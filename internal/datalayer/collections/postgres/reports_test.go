package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Настройка тестовой базы данных
func setupTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Создание таблицы users
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY,
            name TEXT NOT NULL,
            chat_id BIGINT NOT NULL UNIQUE
        );
    `)
	if err != nil {
		t.Fatalf("failed to create table 'users': %v", err)
	}

	// Создание таблицы reports
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS reports (
            id INTEGER PRIMARY KEY,
            user_id BIGINT NOT NULL,
            date TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%S', 'now')),
            request TEXT,
            response TEXT,
            FOREIGN KEY (user_id) REFERENCES users(chat_id) ON DELETE CASCADE
        );
    `)
	if err != nil {
		t.Fatalf("failed to create table 'reports': %v", err)
	}

	// Функция очистки
	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

// Тесты
func TestReportsStorage_GetReportsByChatID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	storage := NewReportsStorage(db)

	ctx := context.Background()
	reportID := int64(1)
	userID := int64(123)
	date := "2023-10-01 12:00:00"
	request := json.RawMessage(`{"Tons": 5, "Refuel": 10}`)
	response := json.RawMessage(`{"TotalFuel": 20}`)

	// Вставка тестового пользователя
	_, err := db.Exec(`INSERT INTO users (id, name, chat_id) VALUES (?, ?, ?)`, userID, "John Doe", userID)
	assert.NoError(t, err)

	// Вставка тестового отчета
	_, err = db.Exec(
		`INSERT INTO reports (id, user_id, date, request, response) VALUES (?, ?, ?, ?, ?)`,
		reportID, userID, date, request, response,
	)
	assert.NoError(t, err)

	// Успешный запрос
	report, err := storage.GetReportsByChatID(ctx, reportID)
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, reportID, report.ID)
	assert.Equal(t, userID, report.UserID)
	assert.Equal(t, date, report.Date)
	assert.Equal(t, request, report.Request)
	assert.Equal(t, response, report.Response)

	// Запрос несуществующего отчета
	_, err = storage.GetReportsByChatID(ctx, 999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
