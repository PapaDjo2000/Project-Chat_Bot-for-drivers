package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer"
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
	//u.Timeout = math.MaxInt
	updates := p.apiBot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			p.logger.Debug().Msg("context is dode")
			return nil
		// закрытие контекста
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
				msg.ReplyMarkup = getGeneral()
				if _, err := p.apiBot.Send(msg); err != nil {
					p.logger.Err(err).Send()
				}
				go p.handleStart(ctx, update)
			default:
				switch update.Message.Text {
				case "Посчитать":
					// Запускаем процесс ввода данных для расчета
					if !isChannelFound {
						userChannel = make(chan tgbotapi.Update)
						p.usersChannels[update.Message.Chat.ID] = userChannel
					}
					go p.handleWork(ctx, update, userChannel)

				case "Выдать":
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

func (p *Processor) isUserAuthorized(ctx context.Context, chatID int64) bool {
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

	var request dto.UserRequest
	var err error

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введи расход:")
	p.apiBot.Send(msg)

	responseTimer := time.NewTimer(2 * time.Minute)
	defer responseTimer.Stop()

	select {
	case <-ctx.Done():
		p.logger.Debug().Msg("ctx is done")
	case <-responseTimer.C:
		p.logger.Debug().Msg("no response from user")
		return
	case response := <-userChannel:
		for {
			request.Consumption, err = strconv.ParseFloat(response.Message.Text, 64)
			if err != nil {
				msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Должно быть число")
				if _, err := p.apiBot.Send(msg); err != nil {
					p.logger.Err(err).Msg("failed to send error message")
					return
				}
				p.logger.Debug().Msg("user entered incorrect Consumption value")
				response = <-userChannel
				continue
			}
			break
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
			for {
				request.Capacity, err = strconv.Atoi(response.Message.Text)
				if err != nil {
					msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Должно быть число")
					if _, err := p.apiBot.Send(msg); err != nil {
						p.logger.Err(err).Msg("failed to send error message")
						return
					}
					p.logger.Debug().Msg("user entered incorrect Consumption value")
					response = <-userChannel
					continue
				}
				break
			}
			msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Введи Остаток топлива:")
			p.apiBot.Send(msg)
			responseTimer.Reset(1 * time.Minute)

			select {
			case <-ctx.Done():
				p.logger.Debug().Msg("ctx is done")
			case <-responseTimer.C:
				p.logger.Debug().Msg("no response from user")
				return
			case response = <-userChannel:

				for {
					request.FuelResidue, err = strconv.ParseFloat(response.Message.Text, 64)
					if err != nil {
						msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Должно быть число")
						if _, err := p.apiBot.Send(msg); err != nil {
							p.logger.Err(err).Msg("failed to send error message")
							return
						}
						p.logger.Debug().Msg("user entered incorrect Consumption value")
						response = <-userChannel
						continue
					}
					break
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

					for {
						request.SpeedometerResidue, err = strconv.Atoi(response.Message.Text)
						if err != nil {
							msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Должно быть число")
							if _, err := p.apiBot.Send(msg); err != nil {
								p.logger.Err(err).Msg("failed to send error message")
								return
							}
							p.logger.Debug().Msg("user entered incorrect Consumption value")
							response = <-userChannel
							continue
						}
						break
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

						for {
							request.Refuel, err = strconv.Atoi(response.Message.Text)
							if err != nil {
								msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Должно быть число")
								if _, err := p.apiBot.Send(msg); err != nil {
									p.logger.Err(err).Msg("failed to send error message")
									return
								}
								p.logger.Debug().Msg("user entered incorrect Consumption value")
								response = <-userChannel
								continue
							}
							break
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

							for {
								request.Distance, err = strconv.Atoi(response.Message.Text)
								if err != nil {
									msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Должно быть число")
									if _, err := p.apiBot.Send(msg); err != nil {
										p.logger.Err(err).Msg("failed to send error message")
										return
									}
									p.logger.Debug().Msg("user entered incorrect Consumption value")
									response = <-userChannel
									continue
								}
								break
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
								for {
									request.QuantityTrips, err = strconv.Atoi(response.Message.Text)
									if err != nil {
										msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Должно быть число")
										if _, err := p.apiBot.Send(msg); err != nil {
											p.logger.Err(err).Msg("failed to send error message")
											return
										}
										p.logger.Debug().Msg("user entered incorrect Consumption value")
										response = <-userChannel
										continue
									}
									break
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

									for {
										request.Tons, err = strconv.Atoi(response.Message.Text)
										if err != nil {
											msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Должно быть число")
											if _, err := p.apiBot.Send(msg); err != nil {
												p.logger.Err(err).Msg("failed to send error message")
												return
											}
											p.logger.Debug().Msg("user entered incorrect Consumption value")
											response = <-userChannel
											continue
										}
										break
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
										for {
											request.Backload, err = strconv.Atoi(response.Message.Text)
											if err != nil {
												msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Должно быть число")
												if _, err := p.apiBot.Send(msg); err != nil {
													p.logger.Err(err).Msg("failed to send error message")
													return
												}
												p.logger.Debug().Msg("user entered incorrect Consumption value")
												response = <-userChannel
												continue
											}
											break
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
											for {
												request.Lifting, err = strconv.ParseFloat(response.Message.Text, 64)
												if err != nil {
													msg := tgbotapi.NewMessage(response.Message.Chat.ID, "Должно быть число")
													if _, err := p.apiBot.Send(msg); err != nil {
														p.logger.Err(err).Msg("failed to send error message")
														return
													}
													p.logger.Debug().Msg("user entered incorrect Consumption value")
													response = <-userChannel
													continue
												}
												break
											}

											vitaldata := p.executorProcessor.Calculate(request)
											str := vitaldata.ToString(request)

											msg := tgbotapi.NewMessage(response.Message.Chat.ID, str)
											p.apiBot.Send(msg)

											err = p.handleUserSaveReport(ctx, update, request, vitaldata)
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

func (p *Processor) handleUserSaveReport(ctx context.Context, update tgbotapi.Update, request dto.UserRequest, vitaldata dto.VitalData) error {
	requestData, err := json.Marshal(request)
	if err != nil {
		p.logger.Err(err).Msg("Failed to marshal request data")
		return fmt.Errorf("failed to marshal request data: %w", err)
	}
	vataldata, err := json.Marshal(vitaldata)
	if err != nil {
		p.logger.Err(err).Msg("Failed to marshal request data")
		return fmt.Errorf("failed to marshal request data: %w", err)
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
