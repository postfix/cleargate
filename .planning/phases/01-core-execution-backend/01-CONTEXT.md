# Phase 1: Core Execution Backend - Context

**Gathered:** 2026-05-06
**Status:** Ready for planning

<domain>
## Phase Boundary

Establish the foundational backend architecture for launching validated CLI jobs without shell execution. Focuses purely on parsing the ToolSpec, validating job values, creating the `argv[]` array, and starting a sandboxed process.
</domain>

<decisions>
## Implementation Decisions

### ToolSpec Parsing
- **D-01:** Use Go structs with `yaml.v3` for strict typed parsing.

### Validation Strategy
- **D-02:** Validate inputs against tool spec immediately at the API boundary, returning 400 Bad Request if invalid.

### Sandboxing Interface
- **D-03:** Define a `ContainerRuntime` interface and use the official Podman Go bindings (`go.podman.io/podman/v6/pkg/bindings`) communicating via the rootless socket to ensure programmatic container management for strict execution isolation.

### Command Construction Mapping
- **D-04:** Isolate the `argv[]` builder into a pure, easily testable function that accepts a ToolSpec and Job Values, and returns `[]string`.

### the agent's Discretion
None — execution strategy strictly follows auto-selected defaults.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Specification
- `SPEC.md` — Core specifications detailing `argv[]` builder mapping and overall execution constraints.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- None (Greenfield phase)

### Established Patterns
- Go standard project layout required.

### Integration Points
- Backend server entrypoint.

</code_context>

<specifics>
## Specific Ideas

- Must use official Podman Go bindings (`go.podman.io/podman/v6/pkg/bindings`).
- Abstraction: Define a `ContainerRuntime` interface to hide Podman details.
- Use rootless socket (`/run/user/$UID/podman/podman.sock` or via `XDG_RUNTIME_DIR`).
- Build with `go build -tags remote` to reduce binary size.

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 01-core-execution-backend*
*Context gathered: 2026-05-06*
