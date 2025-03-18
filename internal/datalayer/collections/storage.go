package collections

import (
	"context"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/models"
)

type (
	Users interface {
		LoadByChatID(ctx context.Context, chatID int)
		CreateUser(ctx context.Context, user models.Users)
	}
)
