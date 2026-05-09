# Phase 7: Core Execution and API Wiring Context

## Context
ClearGate v1.1 has integrated the official Podman rootless bindings and successfully pulls/runs containers with a hardened security profile (Phase 6). However, the API routing and execution pipeline still lacks full end-to-end functionality to be considered a "fully working product." Specifically, input/output file mapping, `argv[]` safe compilation, and live streaming of logs to the frontend via SSE are either stubbed or incomplete. 

## Goal
Wire the complete lifecycle of a tool execution from the React UI to the Podman backend and back, ensuring safe `argv` assembly, correct workspace volume mounting, and live stdout/stderr streaming. 

## Requirements in Scope
- **EXEC-02**: Capture stdout, stderr, exit code, metadata, and output files.
- **EXEC-03**: Stream job status and logs live to the UI via Server-Sent Events (SSE) or WebSockets.
- **EXEC-04**: Backend compiles validated job values into a strict `argv[]` array (preventing arbitrary execution).
- **FILE-01/FILE-02**: Support input file uploads and output artifact retrieval.

## Implementation Guidelines
- Modify `ExecuteRequest` to correctly pass files and map them into the `ContainerRuntime.Create` configuration.
- Implement the `Logs` streaming method in `PodmanRuntime` by attaching to the container output streams and piping them to a Go channel.
- Complete the HTTP SSE implementation in `internal/api/execute.go` (the `/api/jobs/{id}/events` route) to push logs to the UI dynamically.
- Update `web/src/pages/ExecutionPage.tsx` to read the log stream and output status cleanly.
