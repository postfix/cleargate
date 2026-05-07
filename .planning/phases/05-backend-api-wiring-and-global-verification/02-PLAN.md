---
wave: 2
depends_on: [01-PLAN.md]
files_modified:
  - cmd/cleargate/main.go
autonomous: true
requirements_addressed: [TOOL-01, TOOL-02, TOOL-03, TOOL-04, FILE-01, FILE-02, FILE-03, UI-01]
---

# Wave 2: Server Entrypoint & API Wiring

## Objective
Create the main Go entrypoint, load configuration, initialize all dependencies, register the HTTP routes, and serve the React SPA.

## Tasks

```xml
<task id="05-02-01">
  <read_first>
    - internal/api/upload.go
    - internal/api/download.go
    - internal/api/execute.go
    - internal/api/catalog.go
    - internal/api/preset.go
    - internal/api/admin/draft.go
  </read_first>
  <action>
    Create `cmd/cleargate/main.go`.
    In `main()`, initialize all managers and repositories (workspace, podman client, toolspec repo, eino assistant).
    Instantiate all API handlers.
    Use `http.NewServeMux()` to register the following routes:
    - `/api/upload` -> `uploadHandler.HandleUpload`
    - `/api/download` -> `downloadHandler.HandleDownload`
    - `/api/execute` -> `executionHandler.HandleExecute`
    - `/api/catalog` -> `catalogHandler.HandleListCatalog`
    - `/api/presets` -> `presetHandler` methods
    - `/api/admin/drafts` -> `adminHandler` methods
  </action>
  <acceptance_criteria>
    - `cmd/cleargate/main.go` exists and contains `func main()`
    - `cmd/cleargate/main.go` contains `http.NewServeMux()` and registers API routes
  </acceptance_criteria>
</task>

<task id="05-02-02">
  <read_first>
    - cmd/cleargate/main.go
  </read_first>
  <action>
    Modify `cmd/cleargate/main.go` to serve the static SPA files.
    Use `http.FileServer(http.Dir("./web/dist"))`.
    Ensure that any route not matching `/api/*` falls back to serving the SPA (or explicitly registers `/` to serve the static directory).
    Start the server with `http.ListenAndServe(":8080", mux)`.
  </action>
  <acceptance_criteria>
    - `cmd/cleargate/main.go` contains `http.FileServer` serving `./web/dist`
    - `cmd/cleargate/main.go` contains `http.ListenAndServe`
  </acceptance_criteria>
</task>
```
