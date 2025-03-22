package businesslayer

import (
	"context"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/dto"
)

type (
	ChatBot interface {
		SendMessage()
	}

	Users interface {
		CreateIfNotExist(ctx context.Context, user dto.User) error
		LoadByChatID(ctx context.Context, chatID int64) (*dto.User, error)
		Work(userName string) string
	}
)
