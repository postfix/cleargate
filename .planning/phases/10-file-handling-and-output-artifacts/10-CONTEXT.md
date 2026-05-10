# Phase 10: File Handling and Output Artifacts - Context

**Gathered:** 2026-05-10
**Status:** Ready for planning

<domain>
## Phase Boundary
Wire up secure file upload for tool inputs, enforce size/extension limits, and enable downloading output artifacts from completed jobs. Includes gathering execution metadata.
</domain>

<decisions>
## Implementation Decisions

### Upload Validation Strategy
- **D-01:** Stream-based validation. Enforce limits dynamically while parsing `multipart/form-data` using `io.LimitReader`. This safely prevents disk/memory exhaustion from malicious uploads.

### Job Metadata Storage
- **D-02:** Save a `metadata.json` directly into the job's workspace directory when the job completes, capturing stdout/stderr sizes and output file list.

### Workspace Cleanup
- **D-03:** Leave workspaces on disk indefinitely for now.

### the agent's Discretion
- Specific JSON structure for `metadata.json`.
- Approach for retrieving the output file list (e.g., using `filepath.Walk` on the `output/` directory).
</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Architecture & Specs
- `.planning/ROADMAP.md` — Phase 10 goal and constraints
- `.planning/REQUIREMENTS.md` — EXEC-02, FILE-01, FILE-02, FILE-03
</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/api/upload.go` — Existing upload handler (may need updates for limits/extensions).
- `internal/api/download.go` — Existing download handler.
- `internal/workspace/manager.go` — Manages the workspace paths (`input/`, `output/`).
- `internal/job/docker_runner.go` / `internal/api/execute.go` — Execution pipeline where metadata gathering should occur post-execution.
</code_context>

<specifics>
## Specific Ideas
- In `internal/api/upload.go`, cross-reference the uploaded file with the `ToolSpec` to read `inputs[].maxSizeMB` and `inputs[].allowedExtensions`.
- Update `ExecutionHandler` or a post-run hook to write `metadata.json` to the root of the workspace dir.

</specifics>

<deferred>
## Deferred Ideas
- Automated background garbage collection for workspaces.
</deferred>

---
*Phase: 10-file-handling-and-output-artifacts*
*Context gathered: 2026-05-10*
