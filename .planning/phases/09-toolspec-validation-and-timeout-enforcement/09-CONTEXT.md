# Phase 09: ToolSpec Validation and Timeout Enforcement - Context

**Gathered:** 2026-05-09
**Status:** Ready for planning

<domain>
## Phase Boundary

Validate ToolSpec YAML against the schema before approval, reject unknown flags at execution time, and enforce container execution timeouts.
</domain>

<decisions>
## Implementation Decisions

### Validation Engine
- **D-01:** Use `go-playground/validator/v10` for struct validation. It provides robust, standard validation rules via struct tags, making the ToolSpec model declarative and easy to validate.

### Timeout Enforcement
- **D-02:** Use Go's `context.WithTimeout` passed directly to the Podman API calls (`containers.Wait` or `containers.Run`). This is the cleanest and most idiomatic way to enforce process timeouts in Go.

### Unknown Flag Rejection
- **D-03:** Explicit validation in the `ExecutionHandler`. Before constructing the `argv[]` array, explicitly check every provided input key against the allowed `ToolSpec.Flags` and `ToolSpec.Inputs`. If any key is not defined, return a 400 error.

### ToolSpec Versioning
- **D-04:** Append-only history. Approving a new version creates a new database row. The old version remains queryable. (This was explicitly required by the success criteria).

### the agent's Discretion
- Validation error message formatting.
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
- `internal/api/execute.go` — Contains the `ExecutionHandler` where unknown flags should be rejected.
- `internal/job/docker_runner.go` or `podman_runner.go` — Contains the container run logic where timeouts need to be enforced.

</code_context>

<specifics>
## Specific Ideas
- The timeout should enforce `runtime.timeoutSeconds` from the ToolSpec.

</specifics>

<deferred>
## Deferred Ideas
None.
</deferred>

---
*Phase: 09-toolspec-validation-and-timeout-enforcement*
*Context gathered: 2026-05-09 (Auto-chained)*
