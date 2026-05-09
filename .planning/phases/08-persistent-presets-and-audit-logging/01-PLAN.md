---
wave: 1
depends_on: []
files_modified:
  - internal/models/toolspec.go
  - internal/repository/preset_repo.go
  - internal/repository/audit_repo.go
autonomous: true
---

# 01-PLAN: Database Repositories

<objective>
Define the data models and DuckDB repository logic for presets and audit logs.
</objective>

<requirements_addressed>
- PRESET-01
- AUDIT-01
</requirements_addressed>

## Tasks

```xml
<task>
  <description>Update models for Presets and AuditLogs</description>
  <read_first>
    - internal/models/toolspec.go
  </read_first>
  <action>
    Add `ToolID string \`json:"tool_id" yaml:"tool_id"\`` to `Preset` struct in `internal/models/toolspec.go`.
    Add a new `AuditLog` struct to `internal/models/toolspec.go`:
    ```go
    type AuditLog struct {
        JobID     string    `json:"job_id"`
        ToolID    string    `json:"tool_id"`
        ExitCode  int       `json:"exit_code"`
        CreatedAt time.Time `json:"created_at"`
    }
    ```
  </action>
  <acceptance_criteria>
    `grep "ToolID" internal/models/toolspec.go` exits 0.
    `grep "AuditLog struct" internal/models/toolspec.go` exits 0.
  </acceptance_criteria>
</task>
<task>
  <description>Implement PresetRepository</description>
  <read_first>
    - internal/repository/toolspec_repo.go
  </read_first>
  <action>
    Create `internal/repository/preset_repo.go`. Define `PresetRepository` struct wrapping `*sql.DB`. 
    In `NewPresetRepository(db *sql.DB)`, run:
    `CREATE TABLE IF NOT EXISTS presets (id VARCHAR PRIMARY KEY, tool_id VARCHAR, name VARCHAR, values TEXT)` (DuckDB).
    Implement methods:
    - `Save(preset *models.Preset) error` (uses `INSERT ... ON CONFLICT (id) DO UPDATE SET ...`)
    - `ListByTool(toolID string) ([]models.Preset, error)` (select filtering by tool_id)
    - `Delete(id string) error`
  </action>
  <acceptance_criteria>
    `grep "PresetRepository struct" internal/repository/preset_repo.go` exits 0.
    `grep "CREATE TABLE IF NOT EXISTS presets" internal/repository/preset_repo.go` exits 0.
  </acceptance_criteria>
</task>
<task>
  <description>Implement AuditRepository</description>
  <read_first>
    - internal/repository/toolspec_repo.go
  </read_first>
  <action>
    Create `internal/repository/audit_repo.go`. Define `AuditRepository` struct wrapping `*sql.DB`.
    In `NewAuditRepository(db *sql.DB)`, run:
    `CREATE TABLE IF NOT EXISTS audit_logs (job_id VARCHAR PRIMARY KEY, tool_id VARCHAR, exit_code INTEGER, created_at TIMESTAMP)`
    Implement methods:
    - `Log(entry *models.AuditLog) error`
    - `List() ([]models.AuditLog, error)` (ORDER BY created_at DESC)
  </action>
  <acceptance_criteria>
    `grep "AuditRepository struct" internal/repository/audit_repo.go` exits 0.
    `grep "CREATE TABLE IF NOT EXISTS audit_logs" internal/repository/audit_repo.go` exits 0.
  </acceptance_criteria>
</task>
```

<verification>
`go build ./internal/repository` exits 0.
</verification>

<must_haves>
- Tables created using `CREATE TABLE IF NOT EXISTS`.
- `models.Preset` properly aligns with DB representation (`values` stored as JSON string in DB, unmarshaled back to `map[string]interface{}`).
</must_haves>
