package businesslayer

import (
	"context"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/dto"
)

type (
	ChatBot interface {
		SendMessage(chatID int64, message string) error
		Listen(ctx context.Context) error
	}

	Users interface {
		CreateIfNotExist(ctx context.Context, userRequest dto.User) error
		LoadByChatID(ctx context.Context, chatID int64) (*dto.User, error)
		Work(userName string) string
	}
	Executor interface {
		Executor() dto.VitalData
	}
)
