# Phase 11: Gap Closure - Validate Rootless Execution & ToolSpec Seeding

**Gathered:** 2026-05-10
**Status:** Ready for planning

<domain>
## Phase Boundary
This is an audit gap closure phase. The code for Rootless Execution and ToolSpec Seeding was implemented in Phase 06, but no formal verification was captured. The boundary of this phase is strictly to write the tests and/or perform the manual verification steps required to prove that EXEC-06, EXEC-07, TOOL-05, and UI-04 are fully satisfied, and capture the results in `11-VERIFICATION.md`.
</domain>

<decisions>
## Implementation Decisions

### Test Coverage
- Write missing test cases in Go for the `PodmanRuntime` to verify it passes the correct security flags (`--cap-drop=all`, `--read-only`, rootless execution).
- Verify the DB seeder loads `.yaml` files from `./tools/` on startup.
- Manually run the frontend and verify it fetches from the real catalog endpoint instead of using hardcoded mocks.
</decisions>

<canonical_refs>
## Canonical References
- `.planning/ROADMAP.md`
- `.planning/REQUIREMENTS.md`
- `internal/runtime/podman.go`
- `internal/repository/toolspec_repo.go` (SyncFromDirectory)
- `web/src/pages/CatalogPage.tsx`
</canonical_refs>

<code_context>
## Existing Code Insights
The implementation already exists. The `SpecGenerator` in `podman.go` handles EXEC-06/07. `SyncFromDirectory` handles TOOL-05. The React frontend handles UI-04.
</code_context>
