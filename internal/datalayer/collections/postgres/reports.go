package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/PapaDjo2000/Project-Chat_Bot-for-drivers/internal/datalayer/models"
)

// ReportStorage предоставляет методы для работы с таблицей reports.
type ReportsStorage struct {
	db *sql.DB
}

// NewReportStorage создает новый экземпляр ReportStorage.
func NewReportsStorage(db *sql.DB) *ReportsStorage {
	return &ReportsStorage{db: db}
}

// GetReportByID получает отчет по ID.
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

// SaveReport сохраняет новый отчет.
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

// GetUserReports получает все отчеты пользователя.
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
