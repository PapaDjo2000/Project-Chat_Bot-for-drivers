package collections

import (
	"context"
	"database/sql"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/collections/postgres"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/models"
)

type (
	Users interface {
		NewUserStorage(db *sql.DB) *postgres.UserStorage
		GetUserByChatID(ctx context.Context, ChatID int64) (*models.Users, error)
		CreateUser(ctx context.Context, user *models.Users) error
		UpdateUser(ctx context.Context, user *models.Users) error
		DeleteUser(ctx context.Context, id string) error
	}
	Reports interface {
		NewReportsStorage(db *sql.DB) *postgres.ReportsStorage
		GetReportsByChatID(ctx context.Context, ChatID int64) (*models.Reports, error)
		SaveReport(ctx context.Context, report *models.Reports) error
		GetUserReports(ctx context.Context, userID string)
	}
)
