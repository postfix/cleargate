# Phase 05: Backend API Wiring & Global Verification - Context

**Gathered:** 2026-05-07
**Status:** Ready for planning

<domain>
## Phase Boundary

Creating the `main.go` entrypoint, wiring up API routes to the existing backend handlers (`internal/api/*`), serving the built React SPA, and formally verifying all orphaned requirements from previous phases.

</domain>

<decisions>
## Implementation Decisions

### Server Framework
- **D-01:** Use the Go standard library (`net/http`) for the router, leveraging Go 1.22+ routing enhancements to keep dependencies minimal.

### Configuration Strategy
- **D-02:** Use a configuration file (JSON or YAML) to provide the server port, database path, and other runtime settings.

### Static File Serving
- **D-03:** The Go backend will serve the built React SPA (`web/dist`) directly to allow for a single-binary deployment.

### the agent's Discretion
- Exact configuration file format (YAML vs JSON) and internal structure.
- CORS policy details if necessary (though single-binary deployment usually avoids CORS issues).
- Graceful shutdown implementation for the HTTP server.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Project Scope
- `.planning/PROJECT.md` — Overall architecture and constraints.
- `.planning/REQUIREMENTS.md` — All v1 requirements that must be verified.
- `.planning/v1.0-MILESTONE-AUDIT.md` — The audit report detailing all orphaned requirements to be verified.

### Implementation References
- `internal/api/upload.go` — Existing upload handler.
- `internal/api/download.go` — Existing download handler.
- `internal/api/admin/draft.go` — Existing admin handler.

</canonical_refs>
