# Plan 03: Wire Frontend Execution and Artifacts

## Objective
Ensure the frontend successfully sends execution requests, receives the live SSE stream, and can upload/download files associated with the job.

## Steps
1. **Frontend Execution Alignment:**
   - In `web/src/pages/ExecutionPage.tsx`, verify the `LogStream` component connects to `/api/jobs/{id}/events` via `EventSource`.
   - Handle the `message` events to append logs dynamically to the terminal window.
   - Parse the final status event to change the `jobStatus` state to `succeeded` or `failed`.

2. **Artifact Downloads:**
   - Update the UI to display a "Download Artifacts" or "View Results" button once the job succeeds.
   - This should hit `/api/download?job_id={id}&file={filename}`.
   - Verify `internal/api/download.go` correctly serves files out of the isolated workspace directory preventing path traversal.

3. **End-to-End Test:**
   - Submit an `nmap` scan via the UI.
   - Verify the form data correctly builds the argv.
   - Verify the logs stream live to the browser.
