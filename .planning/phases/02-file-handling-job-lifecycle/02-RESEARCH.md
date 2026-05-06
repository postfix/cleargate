# Phase 2: File Handling & Job Lifecycle - Research

## Objective
Support robust job inputs, outputs, log capturing, and artifact retrieval based on the directory layout specified in SPEC.md.

## Workspace Storage Strategy
- **Pattern:** Create a persistent workspace layout per job ID.
- **Structure:**
  - `/var/lib/cleargate/jobs/{job_id}/input/`
  - `/var/lib/cleargate/jobs/{job_id}/output/`
  - `/var/lib/cleargate/jobs/{job_id}/logs/`
  - `/var/lib/cleargate/jobs/{job_id}/metadata/`
- **Implementation:** Go's `os.MkdirAll` can be used to set up the workspace dynamically. We need a `workspace` package that manages resolving paths for these structures.

## File Upload Mechanism (Streaming)
- **Pattern:** Stream HTTP `multipart/form-data` uploads directly to disk.
- **Implementation:** Avoid `r.ParseMultipartForm` which buffers into memory/tmp space. Instead, use `r.MultipartReader()`, iterate through parts, and stream via `io.Copy` to a newly created file in the `input/` directory with a securely generated name.

## Output Artifact Retrieval
- **Pattern:** Direct file serving over HTTP.
- **Implementation:** Provide a `GET /api/jobs/{job_id}/artifacts/{filename}` endpoint. Use `http.ServeFile`. **Security constraint:** The filename must be strictly validated (no `../` traversal, must only resolve within the `output/` directory) to avoid path traversal vulnerabilities.

## Log Capture (Stdout/Stderr)
- **Pattern:** Concurrent file writers in the execution routine.
- **Implementation:** During Podman container execution, open `logs/stdout.log` and `logs/stderr.log` using `os.OpenFile(..., os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)`. When collecting logs from Podman via `bindings/containers.Logs` or attaching to streams, route the output to these files.

## Validation Architecture (Nyquist)
- Use standard Go `net/http/httptest` to validate upload and download handlers.
- Create temporary directories using `t.TempDir()` during tests to simulate the job workspace.
