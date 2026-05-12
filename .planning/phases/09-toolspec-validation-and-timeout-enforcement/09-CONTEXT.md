# Phase 09: ToolSpec Validation and Timeout Enforcement - Context

**Gathered:** 2026-05-12
**Status:** Ready for planning

<domain>
## Phase Boundary

Validate ToolSpec YAML against the schema before approval, reject unknown flags at execution time, and enforce container execution timeouts.
</domain>

<decisions>
## Implementation Decisions

### Validation Engine & Rules
- **D-01:** Use `go-playground/validator/v10` for struct validation.
- **D-02:** Stick to standard struct tags (`required`, `min`, `max`, etc.) to keep validation simple and declarative, relying on the admin's domain knowledge rather than complex regex checks for flag formats.

### Timeout Enforcement Behavior
- **D-03:** Implement graceful termination. When the `timeoutSeconds` is reached, send `SIGTERM` to the container, wait 5 seconds to allow for partial output/logs to flush, and then issue a hard `SIGKILL`.

### Unknown Flag Rejection
- **D-04:** Fail fast with an HTTP 400 response. Before constructing the `argv[]` array, explicitly check every provided input key against the allowed `ToolSpec.Flags` and `ToolSpec.Inputs`. If any unknown key is submitted, reject the entire execution immediately to ensure deterministic execution and prevent wasted compute.

### Validation Feedback UI
- **D-05:** Provide inline form validation errors in the React frontend. The backend should return a structured JSON response detailing exactly which flags failed validation. The UI will use this payload to highlight the specific offending fields in red, rather than relying on a generic global error banner.

### ToolSpec Versioning
- **D-06:** Append-only history. Approving a new version creates a new database row. The old version remains queryable.
</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Architecture & Specs
- `.planning/ROADMAP.md` — Phase 9 goal and constraints

No other external specs referenced.
</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/models/toolspec.go` — Contains the structs to be validated.
- `internal/api/execute.go` — Contains the `ExecutionHandler` where unknown flags should be rejected and validation responses generated.
- `internal/job/docker_runner.go` or `podman_runner.go` — Contains the container run logic where graceful timeouts (`SIGTERM` -> `SIGKILL`) need to be enforced.

</code_context>

<specifics>
## Specific Ideas
- The timeout should enforce `runtime.timeoutSeconds` from the ToolSpec.
</specifics>

<deferred>
## Deferred Ideas
None.
</deferred>
