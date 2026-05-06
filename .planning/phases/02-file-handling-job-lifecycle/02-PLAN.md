---
wave: 1
depends_on: []
files_modified:
  - internal/workspace/workspace.go
  - internal/workspace/workspace_test.go
  - internal/api/upload.go
  - internal/api/upload_test.go
  - internal/api/download.go
  - internal/api/download_test.go
  - internal/runtime/podman/client.go
autonomous: true
---

# Phase 2: File Handling & Job Lifecycle

<objective>
Implement persistent job workspaces, streaming multipart file uploads, artifact retrieval, and execution log capture.
</objective>

<requirements>
- FILE-01: System supports file uploads streamed directly to the job workspace.
- FILE-02: System allows users to securely download generated artifacts from the `output/` directory.
- FILE-03: Workspace follows the deterministic layout defined in SPEC.md.
- EXEC-02: System captures stdout, stderr, and exit codes for every run.
</requirements>

<tasks>

<task>
<id>1</id>
<title>Job Workspace Initialization</title>
<read_first>
- .planning/phases/02-file-handling-job-lifecycle/02-RESEARCH.md
</read_first>
<action>
1. Create `internal/workspace/workspace.go`.
2. Define a base path constant, e.g., `/var/lib/cleargate/jobs/` (or configurable).
3. Implement `func InitializeWorkspace(jobID string) (string, error)` that uses `os.MkdirAll` to create `input/`, `output/`, `logs/`, and `metadata/` inside the job's directory.
4. Implement `internal/workspace/workspace_test.go` to verify the directories are correctly created using a temp base path.
</action>
<acceptance_criteria>
- `internal/workspace/workspace.go` creates exactly the 4 required subdirectories.
- `go test ./internal/workspace/...` exits 0.
</acceptance_criteria>
</task>

<task>
<id>2</id>
<title>Streaming Multipart Upload API</title>
<depends_on>1</depends_on>
<read_first>
- internal/workspace/workspace.go
</read_first>
<action>
1. Create `internal/api/upload.go`.
2. Implement `func HandleUpload(w http.ResponseWriter, r *http.Request)`.
3. Use `r.MultipartReader()` to avoid memory buffering.
4. For each part, if it's a file, securely generate a random string filename or sanitize the uploaded filename.
5. Create a new file in the job's `input/` workspace directory using `os.Create` and copy the data using `io.Copy`.
6. Create `internal/api/upload_test.go` using `httptest` to mock a multipart upload and verify the file is written to the correct `input/` path.
</action>
<acceptance_criteria>
- File is saved to disk without loading entirely into memory.
- `go test ./internal/api/...` exits 0.
</acceptance_criteria>
</task>

<task>
<id>3</id>
<title>Log File Capture in Podman Runtime</title>
<depends_on>1</depends_on>
<read_first>
- internal/runtime/podman/client.go
</read_first>
<action>
1. Update `internal/runtime/podman/client.go`.
2. Modify the `Logs` function (or similar streaming logic during `Wait`/`Start`) to open `logs/stdout.log` and `logs/stderr.log` in the job's workspace.
3. Pipe the container's output streams (via Podman API) into these files concurrently.
4. *Mocking details*: Given the MVP nature and the remote bindings, if actual logs aren't easily fetchable without an active container, ensure the structure to open the files using `os.O_APPEND|os.O_CREATE|os.O_WRONLY` is in place.
</action>
<acceptance_criteria>
- `client.go` includes logic to write to `stdout.log` and `stderr.log`.
- `go test ./internal/runtime/podman/... -tags "remote containers_image_openpgp"` exits 0.
</acceptance_criteria>
</task>

<task>
<id>4</id>
<title>Artifact Download API</title>
<depends_on>1</depends_on>
<read_first>
- internal/workspace/workspace.go
</read_first>
<action>
1. Create `internal/api/download.go`.
2. Implement `func HandleDownload(w http.ResponseWriter, r *http.Request)`.
3. Parse the `job_id` and `filename` from the request URL (e.g., `/api/jobs/{job_id}/artifacts/{filename}`).
4. Validate the path to strictly prevent directory traversal (e.g., using `filepath.Clean` and `filepath.Rel`).
5. Use `http.ServeFile` to serve the file from the job's `output/` directory.
6. Create `internal/api/download_test.go` testing valid and malicious (`../`) paths.
</action>
<acceptance_criteria>
- Files are served successfully.
- Traversal attempts return 403 or 400.
- `go test ./internal/api/...` exits 0.
</acceptance_criteria>
</task>

</tasks>

<verification>
- API tests confirm streaming uploads and secure downloads function without memory bloat or traversal vulnerabilities.
- Workspace creation constructs correct paths.
</verification>

<must_haves>
- Uses `r.MultipartReader` for uploads (no `ParseMultipartForm`).
- Prevents directory traversal in artifact downloads.
</must_haves>
