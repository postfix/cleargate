# Phase 06 Summary: Rootless Execution & ToolSpec Seeding

## Accomplishments
- Added the official Podman Go bindings with CGO constraints fixed for macOS.
- Enforced a hardened security profile natively via `specgen.SpecGenerator` using `--cap-drop=ALL` and `--read-only`.
- Wired up a DuckDB database seeder that reads `.yaml` ToolSpecs from `./tools/`.
- Removed mock FFmpeg tools from `CatalogPage.tsx` and `ExecutionPage.tsx`, integrating them cleanly with real `/api/catalog` outputs.

## User-facing changes
- `nmap` is successfully seeded into the Catalog on backend start.
- The React SPA correctly dynamically renders the Nmap form.
- The execution system correctly binds to the local Podman socket (`CONTAINER_HOST` / defaults) to sandbox processes.
