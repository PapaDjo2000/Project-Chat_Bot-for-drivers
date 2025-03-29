package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/dto"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/collections"
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
	//u.Timeout = math.MaxInt
	updates := p.apiBot.GetUpdatesChan(u)
	for {
		select {
		case <-ctx.Done():
			p.logger.Debug().Msg("context is dode")
			return nil
		// закрытие контекста
		case update := <-updates:
			if update.CallbackQuery != nil {
				go p.handleCallbackQuery(update.CallbackQuery)
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
				go p.handleStart(ctx, update)
			case "work":
				if !p.isUserAuthorized(ctx, update.Message.Chat.ID, update.Message.Chat.UserName) {
					p.suggestToRunStartCommand(update.Message.Chat.ID, update.Message.Chat.UserName)

					continue
				}

				if !isChannelFound {
					userChannel = make(chan tgbotapi.Update)
					p.usersChannels[update.Message.Chat.ID] = userChannel
				}
				go p.handleWork(ctx, update, userChannel)
			case "lalal":
			default:
				// послать сообщение, что не понимаем че он хочет
				if err := p.SendMessage(update.Message.Chat.ID, fmt.Sprintf("Дорогой, %s! Я не понимаю.", update.Message.Chat.UserName)); err != nil {
					p.logger.Err(err).Send()
				}
			}
		}
	}
}

func (p *Processor) handleStart(ctx context.Context, update tgbotapi.Update) {
	// сохранить в базе
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
	// приветствие
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

func (p *Processor) isUserAuthorized(ctx context.Context, chatID int64, userName string) bool {
	// проверить, что в базе есть такой пользователь
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

	// выполняем работу
	if err := p.SendMessage(update.Message.Chat.ID, p.usersProcessor.Work(update.Message.Chat.UserName)); err != nil {
		p.logger.Err(err).Send()
		return
	}

	// задал вопрос.
	if err := p.sendQuestion(update.Message.Chat.ID, 1); err != nil {
		p.logger.Err(err).Send()
		return
	}

	responseTimer := time.NewTimer(1 * time.Minute)
	defer responseTimer.Stop()

	var request dto.UserRequest
	var err error

	// получили ответ.
	select {
	case <-ctx.Done():
		p.logger.Debug().Msg("ctx is done")
	case <-responseTimer.C:
		p.logger.Debug().Msg("no response from user")
		return
	case response := <-userChannel:
		// в зависимости от ответа попросили ввести 2 значения.
		request.Consumption, err = strconv.ParseFloat(response.Message.Text, 64)
		if err != nil {
			msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Расход должен быть дробным")
			p.apiBot.Send(msg)

			p.logger.Debug().Msg("user entered incorrect Consumption value")
			return
		}

		msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Введи Грузоподъемность:")
		p.apiBot.Send(msg)

		responseTimer.Reset(1 * time.Minute)

		select {
		case <-ctx.Done():
			p.logger.Debug().Msg("ctx is done")
		case <-responseTimer.C:
			p.logger.Debug().Msg("no response from user")
			return
		case response = <-userChannel:
			request.Capacity, err = strconv.Atoi(response.Message.Text)
			if err != nil {
				msg := tgbotapi.NewMessage(response.Message.Chat.ID, "должно быть числом")
				p.apiBot.Send(msg)

				p.logger.Debug().Msg("user entered incorrect Capacity value")
				return
			}

			msg = tgbotapi.NewMessage(response.Message.Chat.ID, "Введи Остаток топлива:")
			p.apiBot.Send(msg)

			responseTimer.Reset(1 * time.Minute)

			select {
			case <-ctx.Done():
				p.logger.Debug().Msg("ctx is done")
			case <-responseTimer.C:
				p.logger.Debug().Msg("no response from user")
				return
			case response = <-userChannel:
				request.FuelResidue, err = strconv.ParseFloat(response.Message.Text, 64)
				if err != nil {
					msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Остаток должен быть дробным")
					p.apiBot.Send(msg)

					p.logger.Debug().Msg("user entered incorrect FuelResidue value")
					return
				}

				msg = tgbotapi.NewMessage(response.Message.Chat.ID, "Введи Остаток по спидометру:")
				p.apiBot.Send(msg)

				responseTimer.Reset(1 * time.Minute)

				select {
				case <-ctx.Done():
					p.logger.Debug().Msg("ctx is done")
				case <-responseTimer.C:
					p.logger.Debug().Msg("no response from user")
					return
				case response = <-userChannel:
					request.SpeedometerResidue, err = strconv.Atoi(response.Message.Text)
					if err != nil {
						msg := tgbotapi.NewMessage(response.Message.Chat.ID, "должно быть числом")
						p.apiBot.Send(msg)

						p.logger.Debug().Msg("user entered incorrect SpeedometerResidue value")
						return
					}

					msg = tgbotapi.NewMessage(response.Message.Chat.ID, "Введи Заправку:")
					p.apiBot.Send(msg)

					responseTimer.Reset(1 * time.Minute)

					select {
					case <-ctx.Done():
						p.logger.Debug().Msg("ctx is done")
					case <-responseTimer.C:
						p.logger.Debug().Msg("no response from user")
						return
					case response = <-userChannel:
						request.Refuel, err = strconv.Atoi(response.Message.Text)
						if err != nil {
							msg := tgbotapi.NewMessage(response.Message.Chat.ID, "должно быть числом")
							p.apiBot.Send(msg)

							p.logger.Debug().Msg("user entered incorrect Refuel value")
							return
						}

						msg = tgbotapi.NewMessage(response.Message.Chat.ID, "Введи Расстояние в одну сторону:")
						p.apiBot.Send(msg)

						responseTimer.Reset(1 * time.Minute)

						select {
						case <-ctx.Done():
							p.logger.Debug().Msg("ctx is done")
						case <-responseTimer.C:
							p.logger.Debug().Msg("no response from user")
							return
						case response = <-userChannel:
							request.Distance, err = strconv.Atoi(response.Message.Text)
							if err != nil {
								msg := tgbotapi.NewMessage(response.Message.Chat.ID, "должно быть числом")
								p.apiBot.Send(msg)

								p.logger.Debug().Msg("user entered incorrect Distance value")
								return
							}

							msg = tgbotapi.NewMessage(response.Message.Chat.ID, "Введи Желаемое кол-во рейсов:")
							p.apiBot.Send(msg)

							responseTimer.Reset(1 * time.Minute)

							select {
							case <-ctx.Done():
								p.logger.Debug().Msg("ctx is done")
							case <-responseTimer.C:
								p.logger.Debug().Msg("no response from user")
								return
							case response = <-userChannel:
								request.QuantityTrips, err = strconv.Atoi(response.Message.Text)
								if err != nil {
									msg := tgbotapi.NewMessage(response.Message.Chat.ID, "должно быть числом")
									p.apiBot.Send(msg)

									p.logger.Debug().Msg("user entered incorrect QuantityTrips value")
									return
								}

								msg = tgbotapi.NewMessage(response.Message.Chat.ID, "Введи количество тонн:")
								p.apiBot.Send(msg)

								responseTimer.Reset(1 * time.Minute)

								select {
								case <-ctx.Done():
									p.logger.Debug().Msg("ctx is done")
								case <-responseTimer.C:
									p.logger.Debug().Msg("no response from user")
									return
								case response = <-userChannel:
									request.Tons, err = strconv.Atoi(response.Message.Text)
									if err != nil {
										msg := tgbotapi.NewMessage(response.Message.Chat.ID, "должно быть числом")
										p.apiBot.Send(msg)

										p.logger.Debug().Msg("user entered incorrect Tons value")
										return
									}

									msg = tgbotapi.NewMessage(response.Message.Chat.ID, "Введи Обратные тонны (если нет то 0)")
									p.apiBot.Send(msg)

									responseTimer.Reset(1 * time.Minute)

									select {
									case <-ctx.Done():
										p.logger.Debug().Msg("ctx is done")
									case <-responseTimer.C:
										p.logger.Debug().Msg("no response from user")
										return
									case response = <-userChannel:
										request.Backload, err = strconv.Atoi(response.Message.Text)
										if err != nil {
											msg := tgbotapi.NewMessage(response.Message.Chat.ID, "должно быть числом")
											p.apiBot.Send(msg)

											p.logger.Debug().Msg("user entered incorrect Backload value")
											return
										}

										msg = tgbotapi.NewMessage(response.Message.Chat.ID, "Введи расход на подъемы:")
										p.apiBot.Send(msg)

										responseTimer.Reset(1 * time.Minute)

										select {
										case <-ctx.Done():
											p.logger.Debug().Msg("ctx is done")
										case <-responseTimer.C:
											p.logger.Debug().Msg("no response from user")
											return
										case response = <-userChannel:
											request.Lifting, err = strconv.ParseFloat(response.Message.Text, 64)
											if err != nil {
												msg := tgbotapi.NewMessage(response.Message.Chat.ID, "должно быть числом")
												p.apiBot.Send(msg)

												p.logger.Debug().Msg("user entered incorrect Lifting value")
												return
											}

											vitaldata := p.executorProcessor.Calculate(request)
											str := vitaldata.ToString(request)

											// some logic
											msg := tgbotapi.NewMessage(response.Message.Chat.ID, str)
											p.apiBot.Send(msg)
											err = p.handleUserInteraction(ctx, update, request)
											if err != nil {
												msg := tgbotapi.NewMessage(response.Message.Chat.ID, "No save")
												p.apiBot.Send(msg)

												p.logger.Err(err).Send()
												return
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func (p *Processor) sendQuestion(chatID int64, questionID int) error {
	question, exists := questions[questionID]
	if !exists {
		return errors.New("question is not found")
	}

	var keyboard [][]tgbotapi.InlineKeyboardButton
	for _, option := range question.Options {
		callbackData := strconv.Itoa(option.NextQuestionID)
		row := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(option.Text, callbackData),
		}
		keyboard = append(keyboard, row)
	}

	msg := tgbotapi.NewMessage(chatID, question.Text)
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}

	_, err := p.apiBot.Send(msg)
	if err != nil {
		p.logger.Err(err).Send()
		return err
	}

	return nil
}

func (p *Processor) handleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) {
	chatID := callbackQuery.Message.Chat.ID
	nextQuestionID, err := strconv.Atoi(callbackQuery.Data)
	if err != nil {
		log.Println("Ошибка преобразования callback data:", err)
		return
	}

	switch nextQuestionID {
	case finishOfFirstQuestion:
		msg := tgbotapi.NewMessage(chatID, "Введите а:")
		p.apiBot.Send(msg)
	case defaultFinish:
		msg := tgbotapi.NewMessage(chatID, "Спасибо за ответы! Диалог завершен.")
		p.apiBot.Send(msg)
	default:
		p.sendQuestion(chatID, nextQuestionID)
	}
}
func (p *Processor) handleUserInteraction(ctx context.Context, update tgbotapi.Update, request dto.UserRequest) error {

	report := &models.Reports{
		ID:       uuid.New(),
		UserID:   update.Message.Chat.ID,
		Date:     time.Now(),
		Request:  request,   // Текст сообщения пользователя
		Response: "Готово!", // Ответ бота (можно изменить)
	}

	// Сохраняем отчет в базу данных
	if err := p.reportsCollection.SaveReport(ctx, report); err != nil {
		p.logger.Err(err).Msg("Failed to save report")
		return fmt.Errorf("failed to save user interaction: %w", err)
	}

	// Логируем успешное сохранение
	p.logger.Info().Msgf("Report saved for user %d", update.Message.Chat.ID)

	return nil
}
