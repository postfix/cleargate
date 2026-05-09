package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/postfix/cleargate/internal/models"
)

// AuditRepository manages AuditLogs in a DuckDB database.
type AuditRepository struct {
	db *sql.DB
}

// NewAuditRepository initializes the schema and returns the repo.
func NewAuditRepository(db *sql.DB) (*AuditRepository, error) {
	query := `
	CREATE TABLE IF NOT EXISTS audit_logs (
		job_id VARCHAR PRIMARY KEY,
		tool_id VARCHAR,
		exit_code INTEGER,
		created_at TIMESTAMP
	)`
	
	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create audit_logs table: %w", err)
	}

	return &AuditRepository{db: db}, nil
}

// Log records a new execution.
func (r *AuditRepository) Log(entry *models.AuditLog) error {
	query := `
	INSERT INTO audit_logs (job_id, tool_id, exit_code, created_at)
	VALUES (?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, entry.JobID, entry.ToolID, entry.ExitCode, entry.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to log audit entry: %w", err)
	}

	return nil
}

// List returns all audit logs ordered by creation time descending.
func (r *AuditRepository) List() ([]models.AuditLog, error) {
	rows, err := r.db.Query("SELECT job_id, tool_id, exit_code, created_at FROM audit_logs ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var entry models.AuditLog
		if err := rows.Scan(&entry.JobID, &entry.ToolID, &entry.ExitCode, &entry.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan audit log row: %w", err)
		}
		logs = append(logs, entry)
	}
	
	// If the table is empty or error scanning next rows we return an empty array instead of nil
	if logs == nil {
		logs = []models.AuditLog{}
	}

	return logs, nil
}
