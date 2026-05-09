# Phase 8: Persistent Presets and Audit Logging - Research

## Goal
Store user-saved presets in DuckDB (surviving restarts) with full CRUD, scoped per-tool. Implement an audit log for every job execution to DuckDB and display it in the frontend.

## Architecture & Existing Patterns

### Persistence (DuckDB)
The existing application uses DuckDB via `github.com/marcboeker/go-duckdb` with a repository pattern. The `ToolSpecRepository` creates its table with `CREATE TABLE IF NOT EXISTS` upon initialization. We must replicate this pattern.

**Preset Repository:**
- New table: `presets`
- Fields: `id`, `tool_id`, `name`, `values` (JSON string), `created_at`
- `models.Preset` needs a `ToolID string` field added.
- Existing `PresetHandler` in `internal/api/preset.go` will be refactored to use this new repository instead of its `sync.RWMutex` + slice.
- Endpoints:
  - `POST /api/presets` (Upsert/Save)
  - `GET /api/presets?tool_id=...` (List)
  - `DELETE /api/presets/{id}` (Delete)

**Audit Log Repository:**
- New table: `audit_logs`
- Fields: `job_id`, `tool_id`, `exit_code`, `created_at`
- New struct `models.AuditLog`.
- Endpoints:
  - `GET /api/admin/audit` (List)

### Job Execution & Hooks
Job tracking currently lives in `internal/job/registry.go`. 
- `Registry.Register(...)`
- `Registry.Complete(jobID string, exitCode int, err error)`
The audit log should be written when the job completes. The `Registry` or the execution routine in `internal/api/execute.go` (which waits for the Podman container) needs access to the `AuditRepository` to write the entry.

### Frontend
Currently, `ExecutionPage.tsx` handles presets. 
- It has a `handleSavePreset` that calls `POST /api/presets`. It will need updates to pass `tool_id` in the body.
- It needs a delete mechanism (e.g., a trash icon next to user-created preset pills).
- The `CatalogPage.tsx` or a new top-level `AuditPage.tsx` tab needs to be created to fetch and display `/api/admin/audit`.

## Validation Architecture
- **Preset Persistence**: Saving a preset, restarting the backend, and fetching it again must succeed.
- **Audit Logging**: Running a job (success or failure) must result in a new entry in the `audit_logs` table.

## Implementation Steps
1. Add `ToolID` to `models.Preset` and create `models.AuditLog`.
2. Create `internal/repository/preset_repo.go` and `internal/repository/audit_repo.go`.
3. Refactor `internal/api/preset.go` to use `PresetRepository` and add delete handler.
4. Add `GET /api/admin/audit` to `internal/api/admin.go` or a new handler.
5. Inject `AuditRepository` into `internal/api/execute.go` so it writes to the audit log upon job completion.
6. Wire dependencies in `cmd/cleargate/main.go`.
7. Update React frontend: add `tool_id` to preset payloads, implement preset deletion, and add the Audit Log page.
