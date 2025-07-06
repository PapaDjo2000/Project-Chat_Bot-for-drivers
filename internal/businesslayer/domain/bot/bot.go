package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/domain/bot/keyboard"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/dto"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/collections"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/collections/postgres"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Processor struct {
	apiBot *tgbotapi.BotAPI
	logger zerolog.Logger

	usersProcessor    businesslayer.Users
	executorProcessor businesslayer.Executor
	reportsCollection collections.Reports

	usersChannels map[int64]chan tgbotapi.Update
}

func New(
	token string,
	logger zerolog.Logger,
	usersProcessor businesslayer.Users,
	executorProcessor businesslayer.Executor,
	reportsCollection collections.Reports,
) (*Processor, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Processor{
		apiBot:            bot,
		logger:            logger,
		usersProcessor:    usersProcessor,
		executorProcessor: executorProcessor,
		reportsCollection: reportsCollection,
		usersChannels:     make(map[int64]chan tgbotapi.Update),
	}, nil
}

func (p *Processor) SendMessage(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	if _, err := p.apiBot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (p *Processor) Listen(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)

	updates := p.apiBot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			p.logger.Debug().Msg("context is dode")
			return nil

		case update := <-updates:
			if update.Message == nil || update.Message.Text == "" {
				continue
			}
			userChannel, isChannelFound := p.usersChannels[update.Message.Chat.ID]
			if isChannelFound {
				go func() {
					userChannel <- update
				}()
				continue
			}

			switch update.Message.Command() {
			case "start":
				if !p.isUserAuthorized(ctx, update.Message.Chat.ID) {
					go p.handleStart(ctx, update)
					continue
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать! Выберите действие:")
				msg.ReplyMarkup = keyboard.GetGeneral()
				if _, err := p.apiBot.Send(msg); err != nil {
					p.logger.Err(err).Send()
				}
				go p.handleStart(ctx, update)
			default:
				switch update.Message.Text {
				case "✏️Посчитать📝":

					if !isChannelFound {
						userChannel = make(chan tgbotapi.Update)
						p.usersChannels[update.Message.Chat.ID] = userChannel
					}
					go p.handleWork(ctx, update, userChannel)

				case "🫡Выдать данные📂":
					if !p.isUserAuthorized(ctx, update.Message.Chat.ID) {
						p.suggestToRunStartCommand(update.Message.Chat.ID, update.Message.Chat.UserName)
						continue
					}
					var rep []*models.Reports
					rep, err := p.reportsCollection.GetUserReports(ctx, update.Message.Chat.ID)
					if err != nil {
						p.logger.Err(err).Msg("Failed to get user reports")
						return fmt.Errorf("failed to get user reports: %w", err)
					}
					p.logger.Info().
						Int64("user_id", update.Message.Chat.ID).
						Int("reports_count", len(rep)).
						Msg("Fetched user reports")

					for _, report := range rep {
						RenamedRequest, err := postgres.RenameKeys(report.Request, postgres.RequestKeyMapping)
						if err != nil {
							p.logger.Err(err).Msg("Failed to rename request keys")
							continue
						}
						RenamedResponse, err := postgres.RenameKeys(report.Response, postgres.ResponseKeyMapping)
						if err != nil {
							p.logger.Err(err).Msg("Failed to rename response keys")
							continue
						}
						requestJSON, _ := json.MarshalIndent(RenamedRequest, "", "  ")
						responseJSON, _ := json.MarshalIndent(RenamedResponse, "", "  ")
						msg := fmt.Sprintf("Дата: %s\nВведенные данные : %s\nПолученные Данные: %s",
							report.Date.Format(time.RFC3339),
							requestJSON,
							responseJSON,
						)

						if err := p.SendMessage(update.Message.Chat.ID, msg); err != nil {
							p.logger.Err(err).Send()
						}
					}
				case "🗑Удалить мои данные!":
					err := p.reportsCollection.DeleteUserReports(ctx, update.Message.Chat.ID)
					if err != nil {
						p.logger.Err(err).Send()
					}
					if err := p.SendMessage(update.Message.Chat.ID, "Твои данные удалены!"); err != nil {
						p.logger.Err(err).Send()
					}
				}
			}
		}
	}
}

func (p *Processor) handleStart(ctx context.Context, update tgbotapi.Update) {
	if err := p.usersProcessor.CreateIfNotExist(
		ctx,
		dto.User{
			ID:     uuid.New(),
			Name:   update.Message.Chat.UserName,
			ChatID: update.Message.Chat.ID,
		},
	); err != nil {
		p.logger.Err(err).Send()
		return
	}
	if err := p.SendMessage(update.Message.Chat.ID, fmt.Sprintf("Привет, %s!", update.Message.Chat.UserName)); err != nil {
		p.logger.Err(err).Send()
		return
	}
}

func (p *Processor) suggestToRunStartCommand(chatID int64, userName string) {
	if err := p.SendMessage(chatID, fmt.Sprintf("Дорогой, %s! Воспользуйся командой /start", userName)); err != nil {
		p.logger.Err(err).Send()
	}
}

func (p *Processor) isUserAuthorized(ctx context.Context, chatID int64) bool {
	if _, err := p.usersProcessor.LoadByChatID(ctx, chatID); err != nil {
		p.logger.Err(err).Send()
		return false
	}

	return true
}

func (p *Processor) handleWork(ctx context.Context, update tgbotapi.Update, userChannel chan tgbotapi.Update) {
	defer func() {
		close(userChannel)
		delete(p.usersChannels, update.Message.Chat.ID)
	}()
	var request dto.UserRequest
	type Question struct {
		Prompt       string
		Handler      func(string) error
		ResponseTime time.Duration
	}
	questions := []Question{
		{"Введи расход:", func(input string) error {
			value, err := strconv.ParseFloat(input, 64)
			if err != nil {
				return fmt.Errorf("invalid input for consumption: %w", err)
			}
			if value < 0 {
				return fmt.Errorf("value must be non-negative")
			}
			request.Consumption = value
			return nil
		}, 2 * time.Minute},

		{"Введи расход на подъемы:", func(input string) error {
			value, err := strconv.ParseFloat(input, 64)
			if err != nil {
				return fmt.Errorf("invalid input for lifting: %w", err)
			}
			if value < 0 {
				return fmt.Errorf("value must be non-negative")
			}
			request.Lifting = value
			return nil
		}, 1 * time.Minute},

		{"Введи Грузоподъемность:", func(input string) error {
			value, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("invalid input for capacity: %w", err)
			}
			if value < 0 {
				return fmt.Errorf("value must be non-negative")
			}
			request.Capacity = value
			return nil
		}, 1 * time.Minute},

		{"Введи Остаток по спидометру:", func(input string) error {
			value, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("invalid input for speedometer residue: %w", err)
			}
			if value < 0 {
				return fmt.Errorf("value must be non-negative")
			}
			request.SpeedometerResidue = value
			return nil
		}, 1 * time.Minute},

		{"Введи Остаток топлива:", func(input string) error {
			value, err := strconv.ParseFloat(input, 64)
			if err != nil {
				return fmt.Errorf("invalid input for fuel residue: %w", err)
			}
			if value < 0 {
				return fmt.Errorf("value must be non-negative")
			}
			request.FuelResidue = value
			return nil
		}, 1 * time.Minute},

		{"Введи Заправку:", func(input string) error {
			value, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("invalid input for refuel: %w", err)
			}
			if value < 0 {
				return fmt.Errorf("value must be non-negative")
			}
			request.Refuel = value
			return nil
		}, 1 * time.Minute},

		{"Введи Расстояние в одну сторону:", func(input string) error {
			value, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("invalid input for distance: %w", err)
			}
			if value < 0 {
				return fmt.Errorf("value must be non-negative")
			}
			request.Distance = value
			return nil
		}, 1 * time.Minute},

		{"Введи Желаемое кол-во рейсов:", func(input string) error {
			value, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("invalid input for quantity trips: %w", err)
			}
			if value < 0 {
				return fmt.Errorf("value must be non-negative")
			}
			request.QuantityTrips = value
			return nil
		}, 1 * time.Minute},

		{"Введи количество тонн:", func(input string) error {
			value, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("invalid input for tons: %w", err)
			}
			if value < 0 {
				return fmt.Errorf("value must be non-negative")
			}
			request.Tons = value
			return nil
		}, 1 * time.Minute},

		{"Введи Обратные тонны (если нет то 0):", func(input string) error {
			value, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("invalid input for backload: %w", err)
			}
			if value < 0 {
				return fmt.Errorf("value must be non-negative")
			}
			request.Backload = value
			return nil
		}, 1 * time.Minute},
	}

	processInput := func(question Question) bool {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, question.Prompt)
		msg.ReplyMarkup = keyboard.GetCancel()

		p.apiBot.Send(msg)

		responseTimer := time.NewTimer(question.ResponseTime)
		defer responseTimer.Stop()

		select {
		case <-ctx.Done():
			p.logger.Debug().Msg("ctx is done")
			return false
		case <-responseTimer.C:
			p.logger.Debug().Msg("no response from user")
			return false
		case response := <-userChannel:
			if response.Message.Text == "😬Отмена⚠️" {
				msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Выберите действие:")
				msg.ReplyMarkup = keyboard.GetGeneral()
				p.apiBot.Send(msg)
				return false
			}

			maxAttempts := 3
			attempts := 0

			for {
				attempts++
				err := question.Handler(response.Message.Text)
				if err == nil {
					break
				}

				if attempts >= maxAttempts {
					msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Слишком много некорректных попыток. Попробуйте позже.")
					p.apiBot.Send(msg)
					return false
				}

				msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Должно быть число. Попробуйте снова:")
				p.apiBot.Send(msg)
				response = <-userChannel

				if response.Message.Text == "😬Отмена⚠️" {
					keyboard.GetGeneral()
					return false
				}
			}
		}
		return true
	}
	for _, question := range questions {
		if !processInput(question) {
			return
		}
	}
	vd := p.executorProcessor.Calculate(request)
	str := vd.ToString(request)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)

	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	msg.ReplyMarkup = keyboard.GetGeneral()
	p.apiBot.Send(msg)

	err := p.uSaveRep(ctx, update, request, vd)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось сохранить данные.")
		p.apiBot.Send(msg)
		p.logger.Err(err).Send()
		return
	}
}

func (p *Processor) uSaveRep(ctx context.Context, update tgbotapi.Update, req dto.UserRequest, vd dto.VitalData) error {
	requestData, err := json.Marshal(req)
	if err != nil {
		p.logger.Err(err).Msg("Failed to marshal req data")
		return fmt.Errorf("failed to marshal req data: %w", err)
	}
	vataldata, err := json.Marshal(vd)
	if err != nil {
		p.logger.Err(err).Msg("Failed to marshal req data")
		return fmt.Errorf("failed to marshal req data: %w", err)
	}
	report := &models.Reports{
		ID:       uuid.New(),
		UserID:   update.Message.Chat.ID,
		Date:     time.Now(),
		Request:  json.RawMessage(requestData),
		Response: json.RawMessage(vataldata),
	}
	if err := p.reportsCollection.SaveReport(ctx, report); err != nil {
		p.logger.Err(err).Msg("Failed to save report")
		return fmt.Errorf("failed to save user interaction: %w", err)
	}
	p.logger.Info().Msgf("Report saved for user %d", update.Message.Chat.ID)

	return nil
}
