package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReportsStorage(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	storage := NewReportsStorage(db)
	ctx := context.Background()

	userID := int64(123)
	reportID := uuid.New()
	date := time.Now().UTC().Truncate(time.Second)
	request := json.RawMessage(`{"Tons": 5, "Refuel": 10}`)
	response := json.RawMessage(`{"TotalFuel": 20}`)

	_, err := db.Exec(`INSERT INTO users (id, name, chat_id) VALUES ($1, $2, $3)`,
		userID, "Test User", userID)
	require.NoError(t, err)

	t.Run("SaveReport success", func(t *testing.T) {
		report := &models.Reports{
			ID:       reportID,
			UserID:   userID,
			Date:     date,
			Request:  request,
			Response: response,
		}

		err := storage.SaveReport(ctx, report)
		assert.NoError(t, err)
	})

	t.Run("GetReportsByChatID success", func(t *testing.T) {
		report, err := storage.GetReportsByChatID(ctx, userID)
		require.NoError(t, err)

		assert.Equal(t, reportID, report.ID)
		assert.Equal(t, userID, report.UserID)
		assert.Equal(t, date.UTC(), report.Date.UTC())
		assert.JSONEq(t, string(request), string(report.Request))
		assert.JSONEq(t, string(response), string(report.Response))
	})

	t.Run("GetUserReports success", func(t *testing.T) {
		reports, err := storage.GetUserReports(ctx, userID)
		require.NoError(t, err)
		require.Len(t, reports, 1)

		assert.Equal(t, reportID, reports[0].ID)
	})

	t.Run("DeleteUserReports success", func(t *testing.T) {
		err := storage.DeleteUserReports(ctx, userID)
		assert.NoError(t, err)

		reports, err := storage.GetUserReports(ctx, userID)
		assert.NoError(t, err)
		assert.Empty(t, reports)
	})

	t.Run("GetReportsByChatID not found", func(t *testing.T) {
		_, err := storage.GetReportsByChatID(ctx, 999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()
	connStr := "postgres://user:pass@localhost/test_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY,
			name TEXT,
			chat_id BIGINT
		);
		
		CREATE TABLE IF NOT EXISTS pr.reports (
			id UUID PRIMARY KEY,
			user_id BIGINT REFERENCES users(id),
			date TIMESTAMP,
			request JSONB,
			response JSONB
		);
	`)
	require.NoError(t, err)

	return db, func() {
		db.Exec("DROP TABLE IF EXISTS pr.reports")
		db.Exec("DROP TABLE IF EXISTS users")
		db.Close()
	}
}
