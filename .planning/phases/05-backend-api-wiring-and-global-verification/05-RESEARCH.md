# Phase 05: Backend API Wiring & Global Verification - Research

## Domain Findings
The goal of this phase is to wire up the backend API, serve the built Vite SPA, and globally verify all requirements to close out the v1.0 milestone.

### Current State
- The backend has handlers for `upload.go`, `download.go`, and `admin/draft.go`.
- There are NO handlers for execution (`EXEC-01`, `EXEC-03`), catalog (`UI-02`), or presets (`PRESET-01`).
- The project has no `cmd/cleargate/main.go` entrypoint.

### Architecture Decisions (from CONTEXT.md)
- **Router**: Go standard library `net/http` using `http.NewServeMux()` (Go 1.22+ routing features).
- **Config**: JSON or YAML file for port and DB paths.
- **Serving SPA**: Go will serve `web/dist` directly for a single-binary deployment.

## Implementation Details

### Missing Handlers
To fulfill the orphaned requirements, the following API endpoints must be created and wired:
1. `GET /api/tools` -> Returns approved ToolSpecs from `repository.ToolSpecRepository`
2. `POST /api/jobs` -> Accepts execution payload, builds command using `command.Builder`, executes via `runtime.Runtime`, and handles SSE streaming.
3. `POST /api/presets` -> Saves job preset.
4. `GET /api/presets` -> Lists saved job presets.

### `main.go` Bootstrapping
1. Load configuration file.
2. Initialize dependencies:
   - `repository.ToolSpecRepository`
   - `workspace.Manager`
   - `runtime.PodmanClient`
   - `llm.TemplateAssistant` (CloudWeGo Eino)
3. Instantiate handlers:
   - `api.UploadHandler`
   - `api.DownloadHandler`
   - `admin.AdminHandler`
   - New handlers for execution, catalog, and presets.
4. Register routes on `http.ServeMux`.
5. Add wildcard route to serve static files from `web/dist` (and fallback to index.html for SPA routing).
6. Start HTTP server.

## Validation Architecture
To verify this phase, the application must be built and run.
Verification requires confirming that the `main.go` binary can serve the React SPA, that the frontend can successfully communicate with the backend API, and that all 15 orphaned requirements from previous phases are now demonstrable end-to-end.

## RESEARCH COMPLETE
