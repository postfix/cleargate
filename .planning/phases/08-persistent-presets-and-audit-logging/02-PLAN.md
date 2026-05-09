---
wave: 2
depends_on: [01]
files_modified:
  - internal/api/preset.go
  - internal/api/admin.go
  - internal/api/execute.go
  - cmd/cleargate/main.go
autonomous: true
---

# 02-PLAN: Backend API Wiring

<objective>
Refactor the API layer to use the DuckDB repositories for presets and implement the audit logging hooks.
</objective>

<requirements_addressed>
- PRESET-01
- AUDIT-01
</requirements_addressed>

## Tasks

```xml
<task>
  <description>Refactor Preset API</description>
  <read_first>
    - internal/api/preset.go
    - internal/repository/preset_repo.go
  </read_first>
  <action>
    Update `PresetHandler` in `internal/api/preset.go` to depend on `*repository.PresetRepository` instead of the in-memory array/mutex.
    - `HandleSavePreset`: parse `models.Preset` and call `repo.Save(&preset)`.
    - `HandleListPresets`: parse `tool_id` from query parameter. Call `repo.ListByTool(toolID)`.
    - `HandleDeletePreset`: parse `{id}` from mux/URL and call `repo.Delete(id)`.
  </action>
  <acceptance_criteria>
    `grep "PresetRepository" internal/api/preset.go` exits 0.
    `grep "HandleDeletePreset" internal/api/preset.go` exits 0.
  </acceptance_criteria>
</task>
<task>
  <description>Implement Audit API</description>
  <read_first>
    - internal/repository/audit_repo.go
  </read_first>
  <action>
    Create `internal/api/admin.go` (if it doesn't exist, or add to it). Define `AdminHandler` with a dependency on `*repository.AuditRepository`.
    Implement `HandleListAuditLogs` which calls `repo.List()` and returns JSON.
  </action>
  <acceptance_criteria>
    `grep "HandleListAuditLogs" internal/api/admin.go` exits 0.
  </acceptance_criteria>
</task>
<task>
  <description>Wire Audit Logging into Execution Pipeline</description>
  <read_first>
    - internal/api/execute.go
  </read_first>
  <action>
    Update `ExecuteHandler` in `internal/api/execute.go` to accept `*repository.AuditRepository`.
    In `HandleExecute`, right after `Wait` returns the `exitCode` (around the time the job completes and status is updated), create a `models.AuditLog` struct (using the `jobID`, `req.ToolID`, the captured `exitCode`, and `time.Now()`), and call `auditRepo.Log(&log)`.
  </action>
  <acceptance_criteria>
    `grep "AuditRepository" internal/api/execute.go` exits 0.
  </acceptance_criteria>
</task>
<task>
  <description>Wire Dependencies in main.go</description>
  <read_first>
    - cmd/cleargate/main.go
  </read_first>
  <action>
    In `cmd/cleargate/main.go`:
    1. Pass `db` to `repository.NewPresetRepository(db)` and `repository.NewAuditRepository(db)`.
    2. Pass `presetRepo` to `api.NewPresetHandler(...)`.
    3. Pass `auditRepo` to `api.NewExecuteHandler(...)` and `api.NewAdminHandler(...)`.
    4. Register routes:
       - `r.HandleFunc("/api/presets", presetHandler.HandleDeletePreset).Methods("DELETE")`
       - `r.HandleFunc("/api/admin/audit", adminHandler.HandleListAuditLogs).Methods("GET")`
  </action>
  <acceptance_criteria>
    `grep "NewPresetRepository" cmd/cleargate/main.go` exits 0.
    `grep "NewAuditRepository" cmd/cleargate/main.go` exits 0.
    `grep "admin/audit" cmd/cleargate/main.go` exits 0.
  </acceptance_criteria>
</task>
```

<verification>
`go build ./cmd/cleargate` exits 0.
</verification>

<must_haves>
- Delete route is registered.
- Audit repo logs execution accurately when container stops.
</must_haves>
