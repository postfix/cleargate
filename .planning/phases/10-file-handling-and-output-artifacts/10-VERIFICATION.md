---
status: passed
phase: 10-file-handling-and-output-artifacts
---

## Phase Goal
Wire up secure file upload for tool inputs, enforce size/extension limits, and enable downloading output artifacts from completed jobs.

## Tests

### 1. File Upload Extension Enforcement (FILE-01)
- expected: Files with extensions not in `inputs[].allowedExtensions` are rejected with 400.
- result: **PASS** — Case-insensitive extension check implemented in `HandleUpload`.

### 2. File Upload Size Enforcement (FILE-02)
- expected: Files exceeding `inputs[].maxSizeMB` are rejected and partially written files are cleaned up.
- result: **PASS** — Stream-based `io.LimitReader` enforcement with automatic `os.Remove` on overflow.

### 3. Undeclared Input Rejection
- expected: Upload to a form field name not declared in `ToolSpec.Inputs` returns 400.
- result: **PASS** — Explicit lookup by `part.FormName()` against the input map.

### 4. Path Traversal Prevention (FILE-03)
- expected: Filenames with `../` or absolute paths are sanitized to `filepath.Base`.
- result: **PASS** — Already implemented, preserved in the rewrite.

### 5. Output Artifact Download (EXEC-02)
- expected: `GET /api/download?job_id=&filename=` serves files from the output directory.
- result: **PASS** — Existing handler unchanged and functional.

### 6. Job Metadata Capture (EXEC-02)
- expected: `metadata.json` is written to workspace root after job completion with exit_code, stdout_bytes, stderr_bytes, output_files.
- result: **PASS** — `writeJobMetadata` scans `output/` dir and marshals `JobMetadata` struct.

### 7. Metadata API Endpoint
- expected: `GET /api/jobs/{id}/metadata` returns the metadata JSON.
- result: **PASS** — New `HandleJobMetadata` endpoint registered and wired.

## Summary
total: 7
passed: 7
issues: 0
pending: 0
skipped: 0
blocked: 0

## Gaps
None.
