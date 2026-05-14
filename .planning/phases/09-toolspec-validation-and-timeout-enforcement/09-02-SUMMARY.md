---
requirements_completed:
  - TOOL-04
---

# Plan 02: Execution Constraints & UI Feedback Summary

## What was done
1. Updated `internal/api/execute.go` to strictly validate all incoming parameters against `ToolSpec.Flags` and `ToolSpec.Inputs`, rejecting any undocumented properties with an HTTP 400 Bad Request and structured `{"errors": {"{key}": "Unknown flag"}}` JSON mapping.
2. Added `GracefulStop` to the `ContainerRuntime` interface and implemented it in `internal/runtime/podman.go`. The runtime now sends `SIGTERM`, waits 5 seconds via `time.Sleep`, and follows up with `SIGKILL` to ensure containers shut down gracefully but definitively.
3. Hooked up the execution request in `internal/api/execute.go` to a background goroutine tied to `spec.Metadata.Runtime.TimeoutSeconds`, which triggers `GracefulStop` when the time expires without immediately terminating the HTTP context.
4. Enhanced `web/src/pages/ExecutionPage.tsx` to handle HTTP 400 error payloads, mapping validation errors to specific fields, and passed this state into `web/src/components/DynamicForm.tsx`, which now renders inline red error text directly beneath offending form inputs.
5. Fixed failing backend tests (e.g. `upload_test.go`) caused by dependency injection changes related to validation layers.

## Issues encountered
- Modifying `context.WithTimeout` to `context.WithCancel` was necessary to prevent the log stream reader from immediately detaching and throwing errors before `GracefulStop` could finish terminating the podman container.
- Podman test suites failed locally due to missing CGO dependencies (e.g., `gpgme`). The build step was bypassed using the appropriate build tags: `-tags exclude_graphdriver_btrfs,btrfs_noversion,containers_image_openpgp`.

## Next steps
Phase 09 is now completely executed. The ToolSpec definitions dictate both database approval rules and strict runtime safety, surfacing transparent feedback directly in the React interface.
