# Phase 06: Rootless Execution & ToolSpec Seeding - Context

**Gathered:** 2026-05-08
**Status:** Ready for planning

<domain>
## Phase Boundary

Implementing real Podman rootless isolation via Go bindings and wiring up the backend database to replace frontend UI mocks with real ToolSpecs synced from the filesystem.

</domain>

<decisions>
## Implementation Decisions

### Podman Socket Connection
- **D-01:** Resolve the Podman connection by checking the `CONTAINER_HOST` environment variable first. If unset, fall back automatically to the macOS default (`~/.local/share/containers/podman/machine/podman.sock`) and then the Linux default (`$XDG_RUNTIME_DIR/podman/podman.sock`).

### Hardening Profile
- **D-02:** Enforce a strict security profile by default (`--cap-drop=all`, `--security-opt no-new-privileges`, `--read-only` rootfs, only mounting the job workspace).
- **D-03:** Allow individual ToolSpecs to explicitly override these restrictions (e.g., requesting a writable rootfs or specific capabilities) if a legacy tool requires it.

### Database Seeder Strategy
- **D-04:** Sync ToolSpecs from the filesystem at startup. The backend will read YAML files from a configurable directory (e.g., passed via config or flag, defaulting to `./tools/`) and upsert them into the DuckDB database.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Architecture
- `.planning/PROJECT.md` — Core architectural constraints and tool execution requirements.
- `.planning/ROADMAP.md` — Phase 6 success criteria.

### Implementation References
- `internal/runtime/dummy.go` (if exists) or `cmd/cleargate/main.go` — Where `DummyRuntime` is currently instantiated.
- `internal/api/catalog.go` — The handler that currently serves the frontend.

</canonical_refs>

<specifics>
## Specific Ideas
- The folder containing the ToolSpec YAMLs must be configurable.

</specifics>

<deferred>
## Deferred Ideas
None — discussion stayed within phase scope.
</deferred>

---

*Phase: 06-rootless-execution-toolspec-seeding*
*Context gathered: 2026-05-08*
