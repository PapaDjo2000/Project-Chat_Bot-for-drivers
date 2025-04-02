package bot

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestDB реализует минимальный интерфейс sql.DB для тестов
type TestDB struct {
	queryRowResult *sql.Row
	execError      error
}

func (db *TestDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.queryRowResult
}

func (db *TestDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if db.execError != nil {
		return nil, db.execError
	}
	return &testResult{}, nil
}

type testResult struct{}

func (r *testResult) LastInsertId() (int64, error) {
	return 1, nil
}

func (r *testResult) RowsAffected() (int64, error) {
	return 1, nil
}
func TestReportsStorage_SaveReport(t *testing.T) {
	tests := []struct {
		name        string
		db          *TestDB
		report      *models.Reports
		expectedErr bool
	}{
		{
			name: "successful save",
			db:   &TestDB{},
			report: &models.Reports{
				ID:       uuid.New(),
				UserID:   123,
				Date:     time.Now(),
				Request:  json.RawMessage(`{"key":"value"}`),
				Response: json.RawMessage(`{"key":"value"}`),
			},
			expectedErr: false,
		},
		{
			name: "database error",
			db: &TestDB{
				execError: errors.New("database error"),
			},
			report: &models.Reports{
				ID:       uuid.New(),
				UserID:   123,
				Date:     time.Now(),
				Request:  json.RawMessage(`{"key":"value"}`),
				Response: json.RawMessage(`{"key":"value"}`),
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewReportsStorage(tt.db)
			err := storage.SaveReport(context.Background(), tt.report)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
