package businesslayer

import (
	"Project-Chat_Bot-for-drivers/interanl/businesslayer/dto"
	"context"
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
