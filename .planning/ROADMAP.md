# ClearGate v1.1 Roadmap

**5 phases** | **17 requirements mapped** | Covers all v1 + v1.1 requirements

## Phases

### Phase 6: Rootless Execution & ToolSpec Seeding ✓
**Goal:** Implement real Podman rootless isolation and wire up the backend database to replace frontend UI mocks with real ToolSpecs.
**Requirements:** EXEC-06, EXEC-07, TOOL-05, UI-04
**Status:** Complete (2026-05-08)

### Phase 7: Core Execution and API Wiring ✓
**Goal:** Wire the complete execution lifecycle from React UI → argv[] compilation → Podman container → live SSE log streaming back to the browser.
**Requirements:** EXEC-01, EXEC-03, EXEC-04, EXEC-05, UI-01, UI-02, UI-03
**Depends on:** Phase 6
**Status:** Complete (2026-05-08)

### Phase 8: Persistent Presets and Audit Logging
**Goal:** Store user-saved presets in DuckDB so they survive server restarts. Implement audit logging that records who ran what, when, with which inputs.
**Requirements:** PRESET-01, AUDIT-01
**Depends on:** Phase 7

**Success Criteria:**
1. User-saved presets are persisted in DuckDB and survive server restarts.
2. Presets are scoped per tool (stored with tool_id reference).
3. Every job execution is logged to an `audit_log` table with: job_id, tool_id, user, timestamp, input values, toolspec version, exit code.
4. Audit log entries are queryable via `GET /api/admin/audit`.

### Phase 9: ToolSpec Validation and Timeout Enforcement
**Goal:** Validate ToolSpec YAML against the schema before approval, reject unknown flags at execution time, and enforce container execution timeouts.
**Requirements:** TOOL-03, TOOL-04
**Depends on:** Phase 7

**Success Criteria:**
1. ToolSpec YAML is validated against the Go struct schema on `POST /api/admin/drafts` and `SyncFromDirectory`.
2. Invalid specs are rejected with clear error messages.
3. Unknown flags (not declared in ToolSpec) submitted in execution requests are rejected.
4. `runtime.timeoutSeconds` is enforced — containers exceeding the limit are killed and the job is marked `failed` with `timeout` status.
5. ToolSpec versioning: approving a new version creates a new row, old version remains queryable.

### Phase 10: File Handling and Output Artifacts
**Goal:** Wire up secure file upload for tool inputs, enforce size/extension limits, and enable downloading output artifacts from completed jobs.
**Requirements:** EXEC-02, FILE-01, FILE-02, FILE-03
**Depends on:** Phase 7

**Success Criteria:**
1. `POST /api/upload?job_id=` accepts multipart files and stores them in the job workspace `input/` directory.
2. File extension allowlists from the ToolSpec `inputs[].allowedExtensions` are enforced.
3. File size limits from `inputs[].maxSizeMB` are enforced (request rejected if exceeded).
4. Output artifacts in the workspace `output/` directory are listed and downloadable via `GET /api/download?job_id=&filename=`.
5. Path traversal is prevented (already implemented, needs test coverage).
6. Exit code, stdout/stderr sizes, and output file list are captured as job metadata (EXEC-02).

## Traceability Map
- EXEC-01: Phase 7
- EXEC-02: Phase 10
- EXEC-03: Phase 7
- EXEC-04: Phase 7
- EXEC-05: Phase 7
- EXEC-06: Phase 6
- EXEC-07: Phase 6
- FILE-01: Phase 10
- FILE-02: Phase 10
- FILE-03: Phase 10
- UI-01: Phase 7
- UI-02: Phase 7
- UI-03: Phase 7
- UI-04: Phase 6
- TOOL-03: Phase 9
- TOOL-04: Phase 9
- TOOL-05: Phase 6
- AUDIT-01: Phase 8
- PRESET-01: Phase 8
