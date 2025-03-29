package collections

import (
	"context"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/models"
)

type (
	Users interface {
		GetUserByChatID(ctx context.Context, ChatID int64) (*models.Users, error)
		CreateUser(ctx context.Context, user *models.Users) error
		UpdateUser(ctx context.Context, user *models.Users) error
		DeleteUser(ctx context.Context, id int64) error
	}
	Reports interface {
		GetReportsByChatID(ctx context.Context, ChatID int64) (*models.Reports, error)
		SaveReport(ctx context.Context, report *models.Reports) error
		GetUserReports(ctx context.Context, userID int64) ([]*models.Reports, error)
	}
)
