package users

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/businesslayer/dto"
	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Реализация Users с использованием SQLite
type SQLiteUsers struct {
	db *sql.DB
}

func NewSQLiteUsers(db *sql.DB) *SQLiteUsers {
	return &SQLiteUsers{db: db}
}

func (s *SQLiteUsers) GetUserByChatID(ctx context.Context, ChatID int64) (*models.Users, error) {
	var user models.Users
	query := `SELECT id, name, chat_id FROM users WHERE chat_id = ?`
	err := s.db.QueryRowContext(ctx, query, ChatID).Scan(&user.ID, &user.Name, &user.ChatID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *SQLiteUsers) CreateUser(ctx context.Context, user *models.Users) error {
	query := `INSERT INTO users (id, name, chat_id) VALUES (?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, user.ID, user.Name, user.ChatID)
	return err
}

func (s *SQLiteUsers) UpdateUser(ctx context.Context, user *models.Users) error {
	panic("not implemented")
}

func (s *SQLiteUsers) DeleteUser(ctx context.Context, id int64) error {
	panic("not implemented")
}

// Настройка тестовой базы данных
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	// Открываем SQLite в памяти
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Создаем таблицу users
	_, err = db.Exec(`
        CREATE TABLE users (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            chat_id INTEGER NOT NULL UNIQUE
        )
    `)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	// Функция очистки
	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

// Тесты
func TestProcessor_CreateIfNotExist_UserExists(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	usersRepo := NewSQLiteUsers(db)
	logger := zerolog.Nop()
	processor := NewProcessor(logger, usersRepo)

	ctx := context.Background()
	chatID := int64(12345)
	existingUser := models.Users{
		ID:     uuid.New(),
		Name:   "John Doe",
		ChatID: chatID,
	}

	// Вставка существующего пользователя
	_, err := db.Exec(`INSERT INTO users (id, name, chat_id) VALUES (?, ?, ?)`, existingUser.ID, existingUser.Name, existingUser.ChatID)
	assert.NoError(t, err)

	userRequest := dto.User{
		Name:   "John Doe",
		ChatID: chatID,
	}

	err = processor.CreateIfNotExist(ctx, userRequest)
	assert.NoError(t, err)

	// Проверка, что пользователь не был создан повторно
	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM users WHERE chat_id = ?`, chatID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestProcessor_CreateIfNotExist_UserDoesNotExist(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	usersRepo := NewSQLiteUsers(db)
	logger := zerolog.Nop()
	processor := NewProcessor(logger, usersRepo)

	ctx := context.Background()
	chatID := int64(12345)
	userRequest := dto.User{
		Name:   "John Doe",
		ChatID: chatID,
	}

	err := processor.CreateIfNotExist(ctx, userRequest)
	assert.NoError(t, err)

	// Проверка, что пользователь был создан
	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM users WHERE chat_id = ?`, chatID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestProcessor_LoadByChatID_UserFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	usersRepo := NewSQLiteUsers(db)
	logger := zerolog.Nop()
	processor := NewProcessor(logger, usersRepo)

	ctx := context.Background()
	chatID := int64(12345)
	existingUser := models.Users{
		ID:     uuid.New(),
		Name:   "John Doe",
		ChatID: chatID,
	}

	// Вставка существующего пользователя
	_, err := db.Exec(`INSERT INTO users (id, name, chat_id) VALUES (?, ?, ?)`, existingUser.ID, existingUser.Name, existingUser.ChatID)
	assert.NoError(t, err)

	user, err := processor.LoadByChatID(ctx, chatID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, existingUser.ID, user.ID)
	assert.Equal(t, existingUser.Name, user.Name)
	assert.Equal(t, existingUser.ChatID, user.ChatID)
}

func TestProcessor_LoadByChatID_UserNotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	usersRepo := NewSQLiteUsers(db)
	logger := zerolog.Nop()
	processor := NewProcessor(logger, usersRepo)

	ctx := context.Background()
	chatID := int64(12345)

	user, err := processor.LoadByChatID(ctx, chatID)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "not found", err.Error())
}
