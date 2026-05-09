# Plan 01: Podman Log Streaming and Attachments

## Objective
Implement real-time log capturing for the `PodmanRuntime` to stream output directly from the executing container.

## Steps
1. **Update PodmanRuntime `Logs` method:**
   - In `internal/runtime/podman.go`, refactor the `Logs` method to actually attach to the running container using `containers.Logs` or `containers.Attach`.
   - Pipe the container's standard output and standard error into Go channels (`chan runtime.LogEvent`).
   - Ensure the channel properly closes when the container process exits or the context is cancelled.
   - Use `containers.Wait` to fetch the final exit code and emit a final `LogEvent` containing the exit status.

2. **Integration with `job.Logger`:**
   - Update `internal/job/logger.go` if necessary to properly multiplex or handle the channels coming from the container runtime.
   - Ensure the SSE endpoint in `internal/api/execute.go` correctly consumes the `LogEvent` stream and formats it into standard SSE chunks (`data: {...}\n\n`).
