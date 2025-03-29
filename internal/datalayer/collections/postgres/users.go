package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/models"
)

// UserStorage предоставляет методы для работы с таблицей users.
type UserStorage struct {
	db *sql.DB
}

// NewUserStorage создает новый экземпляр UserStorage.
func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db: db}
}

// GetUserByID получает пользователя по ID.
func (s *UserStorage) GetUserByChatID(ctx context.Context, id int64) (*models.Users, error) {
	var user models.Users
	query := `SELECT id, name, chat_id FROM pr.users WHERE Chat_ID = $1`
	err := s.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.ChatID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// CreateUser создает нового пользователя.
func (s *UserStorage) CreateUser(ctx context.Context, user *models.Users) error {
	query := `
        INSERT INTO pr.users (id, name, chat_id)
        VALUES ($1, $2, $3)
    `
	_, err := s.db.ExecContext(ctx, query, user.ID, user.Name, user.ChatID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// UpdateUser обновляет данные пользователя.
func (s *UserStorage) UpdateUser(ctx context.Context, user *models.Users) error {
	query := `
        UPDATE pr.users
        SET name = $2, chat_id = $3
        WHERE id = $1
    `
	_, err := s.db.ExecContext(ctx, query, user.ID, user.Name, user.ChatID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser удаляет пользователя по ID.
func (s *UserStorage) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM pr.users WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
