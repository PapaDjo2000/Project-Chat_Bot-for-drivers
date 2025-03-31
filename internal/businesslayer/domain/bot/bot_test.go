package bot

import (
	"context"
	"errors"
	"testing"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/bisinesslayer/domain/bot"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/dto"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type MockBotAPI struct {
	SendFunc func(msg tgbotapi.MessageConfig) (tgbotapi.Message, error)
}
type MockReportsCollection struct {
	GetUserReportsFunc    func(ctx context.Context, userID int64) ([]*models.Reports, error)
	DeleteUserReportsFunc func(ctx context.Context, userID int64) error
	SaveReportFunc        func(ctx context.Context, report *models.Reports) error
}

func (m *MockReportsCollection) GetUserReports(ctx context.Context, userID int64) ([]*models.Reports, error) {
	return m.GetUserReportsFunc(ctx, userID)
}

func (m *MockReportsCollection) DeleteUserReports(ctx context.Context, userID int64) error {
	return m.DeleteUserReportsFunc(ctx, userID)
}

func (m *MockReportsCollection) SaveReport(ctx context.Context, report *models.Reports) error {
	return m.SaveReportFunc(ctx, report)
}
func (m *MockBotAPI) Send(msg tgbotapi.MessageConfig) (tgbotapi.Message, error) {
	return m.SendFunc(msg)
}

func TestSendMessage(t *testing.T) {
	mockBot := &MockBotAPI{
		SendFunc: func(msg tgbotapi.MessageConfig) (tgbotapi.Message, error) {
			if msg.Text == "Test Message" {
				return tgbotapi.Message{}, nil
			}
			return tgbotapi.Message{}, errors.New("failed to send message")
		},
	}

	logger := zerolog.New(nil)
	processor := &bot.Processor{
		apiBot: mockBot,
		logger: logger,
	}

	err := processor.SendMessage(123456789, "Test Message")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = processor.SendMessage(123456789, "Invalid Message")
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
}
func TestHandleWork(t *testing.T) {
	mockBot := &MockBotAPI{
		SendFunc: func(msg tgbotapi.MessageConfig) (tgbotapi.Message, error) {
			return tgbotapi.Message{}, nil
		},
	}

	mockUsersProcessor := &MockUsersProcessor{}
	mockExecutorProcessor := &MockExecutorProcessor{}
	mockReportsCollection := &MockReportsCollection{}

	logger := zerolog.New(nil)
	processor := &bot.Processor{
		apiBot:            mockBot,
		logger:            logger,
		usersProcessor:    mockUsersProcessor,
		executorProcessor: mockExecutorProcessor,
		reportsCollection: mockReportsCollection,
		usersChannels:     make(map[int64]chan tgbotapi.Update),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userChannel := make(chan tgbotapi.Update, 10)
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: 123456789},
		},
	}

	go func() {
		userChannel <- tgbotapi.Update{
			Message: &tgbotapi.Message{Text: "10"},
		}
		userChannel <- tgbotapi.Update{
			Message: &tgbotapi.Message{Text: "5000"},
		}
		close(userChannel)
	}()

	processor.handleWork(ctx, update, userChannel)

	// Проверяем, что данные были обработаны корректно
	if len(processor.usersChannels) != 0 {
		t.Fatalf("Expected usersChannels to be empty, got %d", len(processor.usersChannels))
	}
}

type MockReportsCollection struct {
	SaveReportFunc func(ctx context.Context, report *models.Reports) error
}

func (m *MockReportsCollection) SaveReport(ctx context.Context, report *models.Reports) error {
	return m.SaveReportFunc(ctx, report)
}

func TestHandleUserSaveReport(t *testing.T) {
	mockReportsCollection := &MockReportsCollection{
		SaveReportFunc: func(ctx context.Context, report *models.Reports) error {
			if report.UserID != 123456789 {
				return errors.New("invalid user ID")
			}
			return nil
		},
	}

	logger := zerolog.New(nil)
	processor := &bot.Processor{
		reportsCollection: mockReportsCollection,
		logger:            logger,
	}

	ctx := context.Background()
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: 123456789},
		},
	}

	request := dto.UserRequest{
		Consumption: 10.5,
		Capacity:    5000,
	}
	vitaldata := dto.VitalData{
		Result: "Success",
	}

	err := processor.handleUserSaveReport(ctx, update, request, vitaldata)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
