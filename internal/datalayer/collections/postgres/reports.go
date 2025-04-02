package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/models"
)

var RequestKeyMapping = map[string]string{
	"Tons":               "Тонны",
	"Refuel":             "Заправка",
	"Lifting":            "Подъемы",
	"Backload":           "Обратная загрузка",
	"Capacity":           "Вместимость Груза",
	"Distance":           "Расстояние",
	"Consumption":        "Расход топлива",
	"FuelResidue":        "Остаток топлива",
	"QuantityTrips":      "Количество рейсов",
	"SpeedometerResidue": "Остаток спидометра",
}

var ResponseKeyMapping = map[string]string{
	"UserId":            "ID пользователя",
	"Lifting":           "Подъемы",
	"Wastage":           "Расход на пробег",
	"DailyRun":          "Суточный пробег",
	"DailyRate":         "Суточная расход",
	"TotalFuel":         "Общий расход топлива",
	"Underfuel":         "Либо расход либо недостаток 3 действие",
	"Undelivery":        "Недоставлено тонн",
	"OperatingDistance": "Пройденное расстояние за день",
}

type ReportsStorage struct {
	db *sql.DB
}

func NewReportsStorage(db *sql.DB) *ReportsStorage {
	return &ReportsStorage{db: db}
}

func (s *ReportsStorage) GetReportsByChatID(ctx context.Context, ChatID int64) (*models.Reports, error) {
	var report models.Reports
	query := `SELECT id, user_id, date, requestб response FROM pr.reports WHERE id = $1`
	err := s.db.QueryRowContext(ctx, query, ChatID).Scan(&report.ID, &report.UserID, &report.Date, &report.Request, &report.Response)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("report with id %v not found", ChatID)
		}
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	return &report, nil
}

func (s *ReportsStorage) SaveReport(ctx context.Context, report *models.Reports) error {
	query := `
        INSERT INTO pr.reports (id, user_id, date, request, response)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := s.db.ExecContext(ctx, query, report.ID, report.UserID, report.Date, report.Request, report.Response)
	if err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}
	return nil
}

func (s *ReportsStorage) GetUserReports(ctx context.Context, userID int64) ([]*models.Reports, error) {
	query := `SELECT id, user_id, date, request, response FROM pr.reports WHERE user_id = $1`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user reports: %w", err)
	}
	defer rows.Close()

	var reports []*models.Reports
	for rows.Next() {
		var report models.Reports
		if err := rows.Scan(&report.ID, &report.UserID, &report.Date, &report.Request, &report.Response); err != nil {
			return nil, fmt.Errorf("failed to scan report: %w", err)
		}
		reports = append(reports, &report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration over rows: %w", err)
	}

	return reports, nil
}
func (s *ReportsStorage) DeleteUserReports(ctx context.Context, userID int64) error {
	query := `DELETE FROM pr.reports WHERE user_id = $1`

	result, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user reports: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no reports found for user_id: %d", userID)
	}
	return nil
}

func RenameKeys(data json.RawMessage, keyMapping map[string]string) (map[string]interface{}, error) {
	var parsedData map[string]interface{}
	if err := json.Unmarshal(data, &parsedData); err != nil {
		return nil, err
	}
	renamedData := make(map[string]interface{})
	for key, value := range parsedData {
		newKey, exists := keyMapping[key]
		if exists {
			renamedData[newKey] = value
		} else {
			renamedData[key] = value
		}
	}
	return renamedData, nil
}
