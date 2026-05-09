# Phase 8: Persistent Presets and Audit Logging - Context

**Gathered:** 2026-05-09
**Status:** Ready for planning

<domain>
## Phase Boundary

Move presets from ephemeral in-memory storage to DuckDB persistence, scope them by tool, and add full CRUD lifecycle. Implement an audit trail that logs every job execution to DuckDB and displays it in a new React frontend tab.
</domain>

<decisions>
## Implementation Decisions

### Preset Scoping
- **D-01:** Presets must be scoped per-tool. The database schema needs a `tool_id` column, and the API should filter by it (`GET /api/presets?tool_id=...`).

### Audit Data Shape
- **D-02:** Minimal logging. The audit record should only capture: `job_id`, `tool_id`, `timestamp`, and `exit_code`.

### Audit Queryability
- **D-03:** Full UI. Add an "Audit Log" tab to the React frontend to view execution history, backed by a new `GET /api/admin/audit` endpoint.

### Preset Lifecycle
- **D-04:** Full CRUD. Users must be able to not only save, but also delete and update existing presets from the UI.
</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Architecture & Specs
- `.planning/ROADMAP.md` — Phase 8 goal and constraints
- `.planning/REQUIREMENTS.md` — PRESET-01 and AUDIT-01

No other external specs referenced.
</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `ToolSpecRepository` in `internal/repository/toolspec_repo.go` — established pattern for DuckDB connection, table creation, and CRUD operations.
- `models.Preset` in `internal/models/toolspec.go` — struct already has JSON and YAML tags.

### Established Patterns
- In-memory registry hooks: `internal/job/registry.go` handles active jobs. Audit log hooks could be added here or in the execution pipeline.

### Integration Points
- `internal/api/preset.go` — currently uses in-memory `sync.RWMutex`, needs complete rewrite to use a DuckDB repository.
- React UI: Needs a new tab component for the Audit Log and updates to the existing execution view to support preset deletion/updating.
</code_context>

<specifics>
## Specific Ideas

No specific references provided.
</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope
</deferred>

---

*Phase: 08-persistent-presets-and-audit-logging*
*Context gathered: 2026-05-09*
