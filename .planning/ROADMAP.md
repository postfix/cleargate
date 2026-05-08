# ClearGate v1.1 Roadmap

**1 phase** | **4 requirements mapped** | All v1.1 requirements covered ✓

## Phases

### Phase 6: Rootless Execution & ToolSpec Seeding
**Goal:** Implement real Podman rootless isolation and wire up the backend database to replace frontend UI mocks with real ToolSpecs.
**Requirements:** EXEC-06, EXEC-07, TOOL-05, UI-04

**Success Criteria:**
1. The Go backend automatically seeds the DuckDB database with at least one real ToolSpec (e.g., `nmap` or `ffmpeg`) on startup.
2. The React UI reads from `/api/catalog` instead of `/api/admin/drafts` (or vice versa), successfully eliminating the fallback mock FFmpeg schema.
3. Tool executions are routed through the official Podman Go Bindings over the API socket.
4. The backend compiles successfully on macOS using CGO, resolving the `gpgme`/`btrfs` header dependencies (via Homebrew or build tags).
5. Containers are spawned using hardened flags (e.g., `--cap-drop=all`, `--security-opt no-new-privileges`, `--read-only`).

## Traceability Map
- EXEC-06: Phase 6
- EXEC-07: Phase 6
- TOOL-05: Phase 6
- UI-04: Phase 6
