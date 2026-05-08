---
wave: 2
depends_on: ["01-PLAN.md"]
files_modified:
  - internal/repository/toolspec_repo.go
  - cmd/cleargate/main.go
  - web/src/pages/CatalogPage.tsx
  - web/src/pages/ExecutionPage.tsx
autonomous: true
---

# Plan 02: Database Seeder & UI Integration

<objective>
Implement the directory sync logic to seed the DuckDB database with real ToolSpecs on startup, and update the React SPA to read from the real API endpoints instead of falling back to hardcoded mocks.
</objective>

<requirements>
- TOOL-05: Database Seeder
- UI-04: Remove UI Mocks
</requirements>

<tasks>

## 1. Implement Directory Sync in Repository
Add the ability to read YAML files from a directory and upsert them into the database.

<read_first>
- `internal/repository/toolspec_repo.go`
</read_first>

<action>
1. Add a `SyncFromDirectory(dirPath string) error` method to `ToolSpecRepository`.
2. Iterate over all `.yaml` and `.yml` files in `dirPath`.
3. For each file, read the bytes, unmarshal into `models.ToolSpec` using `yaml.v3`.
4. Call `SaveDraft(spec)` to insert the record.
5. Immediately call `Approve(id)` so the tool is approved and available in the catalog.
6. Create an example `tools/nmap.yaml` file in the project root with a valid ToolSpec definition for nmap to serve as the initial seed.
</action>

<acceptance_criteria>
- `internal/repository/toolspec_repo.go` contains `SyncFromDirectory` method.
- The method processes files ending in `.yaml` or `.yml` and marks them as approved.
- `tools/nmap.yaml` exists in the project root.
</acceptance_criteria>

## 2. Wire Seeder and Podman to Main
Update the entrypoint to use the new implementations.

<read_first>
- `cmd/cleargate/main.go`
</read_first>

<action>
1. Define a `-tools-dir` CLI flag defaulting to `"./tools"`.
2. After initializing the `ToolSpecRepository`, call `repo.SyncFromDirectory(toolsDir)`. Ignore errors if the directory doesn't exist, but log a warning.
3. Replace `var runtimeClient runtime.ContainerRuntime = &DummyRuntime{}` with `runtimeClient, err := runtime.NewPodmanRuntime()`. If it returns an error, log a warning and fall back to `DummyRuntime` (so development can continue even if Podman isn't running).
</action>

<acceptance_criteria>
- `cmd/cleargate/main.go` calls `SyncFromDirectory`.
- `cmd/cleargate/main.go` attempts to initialize `NewPodmanRuntime`.
</acceptance_criteria>

## 3. Remove UI Mocks and Fix API Endpoints
Ensure the frontend fetches the catalog from the correct route and stops showing the FFmpeg mock.

<read_first>
- `web/src/pages/CatalogPage.tsx`
- `web/src/pages/ExecutionPage.tsx`
</read_first>

<action>
1. In `CatalogPage.tsx`, change `fetch('/api/admin/tools/drafts')` to `fetch('/api/catalog')`.
2. Remove the `.catch(err => ...)` fallback block that sets the mock `FFmpeg`, `ImageMagick`, and `Nmap` tools.
3. In `ExecutionPage.tsx`, change `fetch('/api/admin/tools/drafts')` to `fetch('/api/catalog')`.
4. Remove the `.catch(err => ...)` fallback block that sets the `mockSpec` for FFmpeg.
5. In both files, if the fetch fails or returns empty, correctly set the `error` state and display "Failed to load tools" or "Tool not found" rather than mocking.
</action>

<acceptance_criteria>
- `web/src/pages/CatalogPage.tsx` no longer contains the string `FFmpeg Video Encoder`.
- `web/src/pages/ExecutionPage.tsx` no longer contains the `mockSpec` definition.
- Both files fetch from `/api/catalog`.
</acceptance_criteria>

</tasks>

<verification>
- Starting the Go backend with a `./tools/` directory containing `nmap.yaml` successfully inserts it into DuckDB.
- Opening the frontend React SPA displays the real Nmap tool.
- Clicking the tool loads the dynamically generated form from the real YAML schema.
</verification>
