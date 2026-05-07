---
wave: 1
depends_on: []
files_modified:
  - internal/api/execute.go
  - internal/api/catalog.go
  - internal/api/preset.go
autonomous: true
requirements_addressed: [EXEC-01, EXEC-02, EXEC-03, EXEC-04, EXEC-05, UI-02, UI-03, PRESET-01, AUDIT-01]
---

# Wave 1: Backend API Handlers

## Objective
Create the remaining API handlers required to support the ClearGate frontend, specifically execution, catalog, and presets.

## Tasks

```xml
<task id="05-01-01">
  <read_first>
    - internal/api/upload.go
    - internal/runtime/runtime.go
    - internal/job/logger.go
  </read_first>
  <action>
    Create `internal/api/execute.go`.
    Implement `ExecutionHandler` struct with `runtime.Runtime`, `workspace.Manager`, and `job.Logger` dependencies.
    Implement `HandleExecute` to accept a POST request with job parameters, start the job in the sandbox, and stream output via SSE (Server-Sent Events).
  </action>
  <acceptance_criteria>
    - `internal/api/execute.go` contains `type ExecutionHandler struct`
    - `internal/api/execute.go` contains `HandleExecute` method
  </acceptance_criteria>
</task>

<task id="05-01-02">
  <read_first>
    - internal/repository/toolspec_repo.go
  </read_first>
  <action>
    Create `internal/api/catalog.go`.
    Implement `CatalogHandler` struct with `*repository.ToolSpecRepository` dependency.
    Implement `HandleListCatalog` to return a JSON array of all approved ToolSpecs from the repository.
  </action>
  <acceptance_criteria>
    - `internal/api/catalog.go` contains `type CatalogHandler struct`
    - `internal/api/catalog.go` contains `HandleListCatalog` method
  </acceptance_criteria>
</task>

<task id="05-01-03">
  <read_first>
    - internal/api/catalog.go
  </read_first>
  <action>
    Create `internal/api/preset.go`.
    Implement `PresetHandler` struct.
    Implement `HandleSavePreset` (POST) to save job form state as a preset.
    Implement `HandleListPresets` (GET) to return saved presets.
    (For MVP, you may use an in-memory map or write to a simple JSON file in the workspace).
  </action>
  <acceptance_criteria>
    - `internal/api/preset.go` contains `HandleSavePreset` and `HandleListPresets`
  </acceptance_criteria>
</task>
```
