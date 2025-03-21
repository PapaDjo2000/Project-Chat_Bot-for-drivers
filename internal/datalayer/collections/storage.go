package collections

import (
	"context"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/models"
)

type (
	Users interface {
		GetByChatID(ctx context.Context, chatID int64) (*models.Users, error)
		Create(ctx context.Context, user models.Users) error
	}
)
