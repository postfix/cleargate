package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/postfix/cleargate/internal/models"
	_ "github.com/marcboeker/go-duckdb"
	"gopkg.in/yaml.v3"
)

// ToolSpecRepository manages ToolSpecs in a DuckDB database.
type ToolSpecRepository struct {
	db *sql.DB
}

// ToolSpecRecord represents a row in the toolspecs table.
type ToolSpecRecord struct {
	ID        string
	Name      string
	Version   string
	Status    string
	Content   string
	CreatedAt time.Time
}

// NewToolSpecRepository initializes a DuckDB connection and creates the schema.
func NewToolSpecRepository(dbPath string) (*ToolSpecRepository, error) {
	db, err := sql.Open("duckdb", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open duckdb: %w", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS toolspecs (
		id VARCHAR PRIMARY KEY,
		name VARCHAR,
		version VARCHAR,
		status VARCHAR,
		content TEXT,
		created_at TIMESTAMP
	)`
	
	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create toolspecs table: %w", err)
	}

	return &ToolSpecRepository{db: db}, nil
}

// Close closes the database connection.
func (r *ToolSpecRepository) Close() error {
	return r.db.Close()
}

// SaveDraft inserts or updates a ToolSpec as a draft.
func (r *ToolSpecRepository) SaveDraft(spec *models.ToolSpec) error {
	contentBytes, err := yaml.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to marshal toolspec: %w", err)
	}

	id := fmt.Sprintf("%s-%s", spec.Metadata.Name, spec.Metadata.Version)

	query := `
	INSERT INTO toolspecs (id, name, version, status, content, created_at)
	VALUES (?, ?, ?, 'draft', ?, ?)
	ON CONFLICT (id) DO UPDATE SET
		status = 'draft',
		content = EXCLUDED.content,
		created_at = EXCLUDED.created_at
	`

	_, err = r.db.Exec(query, id, spec.Metadata.Name, spec.Metadata.Version, string(contentBytes), time.Now())
	if err != nil {
		return fmt.Errorf("failed to save draft: %w", err)
	}

	return nil
}

// ListDrafts returns all ToolSpecs with 'draft' status.
func (r *ToolSpecRepository) ListDrafts() ([]ToolSpecRecord, error) {
	rows, err := r.db.Query("SELECT id, name, version, status, content, created_at FROM toolspecs WHERE status = 'draft'")
	if err != nil {
		return nil, fmt.Errorf("failed to query drafts: %w", err)
	}
	defer rows.Close()

	var drafts []ToolSpecRecord
	for rows.Next() {
		var rec ToolSpecRecord
		if err := rows.Scan(&rec.ID, &rec.Name, &rec.Version, &rec.Status, &rec.Content, &rec.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		drafts = append(drafts, rec)
	}

	return drafts, nil
}

// Approve updates the status of a ToolSpec to 'approved'.
func (r *ToolSpecRepository) Approve(id string) error {
	res, err := r.db.Exec("UPDATE toolspecs SET status = 'approved' WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to approve toolspec: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("toolspec not found: %s", id)
	}

	return nil
}
