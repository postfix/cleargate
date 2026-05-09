package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/postfix/cleargate/internal/models"
)

// PresetRepository manages Presets in a DuckDB database.
type PresetRepository struct {
	db *sql.DB
}

// NewPresetRepository initializes the schema and returns the repo.
func NewPresetRepository(db *sql.DB) (*PresetRepository, error) {
	query := `
	CREATE TABLE IF NOT EXISTS presets (
		id VARCHAR PRIMARY KEY,
		tool_id VARCHAR,
		name VARCHAR,
		description VARCHAR,
		visibility VARCHAR,
		locked BOOLEAN,
		values TEXT
	)`
	
	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create presets table: %w", err)
	}

	return &PresetRepository{db: db}, nil
}

// Save inserts or updates a preset.
func (r *PresetRepository) Save(preset *models.Preset) error {
	valuesBytes, err := json.Marshal(preset.Values)
	if err != nil {
		return fmt.Errorf("failed to marshal preset values: %w", err)
	}

	query := `
	INSERT INTO presets (id, tool_id, name, description, visibility, locked, values)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT (id) DO UPDATE SET
		name = EXCLUDED.name,
		description = EXCLUDED.description,
		visibility = EXCLUDED.visibility,
		locked = EXCLUDED.locked,
		values = EXCLUDED.values
	`

	_, err = r.db.Exec(query, preset.ID, preset.ToolID, preset.Name, preset.Description, preset.Visibility, preset.Locked, string(valuesBytes))
	if err != nil {
		return fmt.Errorf("failed to save preset: %w", err)
	}

	return nil
}

// ListByTool returns all presets for a specific tool.
func (r *PresetRepository) ListByTool(toolID string) ([]models.Preset, error) {
	rows, err := r.db.Query("SELECT id, tool_id, name, description, visibility, locked, values FROM presets WHERE tool_id = ?", toolID)
	if err != nil {
		return nil, fmt.Errorf("failed to query presets: %w", err)
	}
	defer rows.Close()

	var presets []models.Preset
	for rows.Next() {
		var preset models.Preset
		var valuesStr string
		if err := rows.Scan(&preset.ID, &preset.ToolID, &preset.Name, &preset.Description, &preset.Visibility, &preset.Locked, &valuesStr); err != nil {
			return nil, fmt.Errorf("failed to scan preset row: %w", err)
		}
		
		var values map[string]interface{}
		if err := json.Unmarshal([]byte(valuesStr), &values); err != nil {
			return nil, fmt.Errorf("failed to unmarshal preset values: %w", err)
		}
		preset.Values = values

		presets = append(presets, preset)
	}
	
	return presets, nil
}

// Delete removes a preset by ID.
func (r *PresetRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM presets WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete preset: %w", err)
	}
	return nil
}
