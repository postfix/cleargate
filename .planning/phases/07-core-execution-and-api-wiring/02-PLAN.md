# Plan 02: Argv Construction and Workspace Mapping

## Objective
Safely compile user-provided form values into a rigid `argv[]` command array based on the ToolSpec definition, and map the user workspace into the Podman container.

## Steps
1. **Argv Compilation Logic:**
   - In `internal/api/execute.go`, implement a parsing function that matches the `ExecuteRequest.Values` against the `ToolSpec.Flags` and `ToolSpec.Inputs`.
   - Build a `[]string` containing `argv[0]` (the executable) followed by the correctly formatted flags (e.g., `-sS`, `--target=10.0.0.1`).
   - Ensure malicious input (shell injections like `||`, `&&`, `;`) are neutralized because they are passed strictly as array elements to the container runtime, never executed via `sh -c`.

2. **Workspace Volume Mounting:**
   - Modify the `CreateContainerRequest` generated in `HandleExecute` to specify the local workspace path (`workspaceManager.GetWorkspacePath(jobId)`).
   - In `internal/runtime/podman.go`, update the `specgen` configuration to mount the workspace directory into the container at `/workspace` (or another designated path).
   - Set the container's working directory to the mount point so tool outputs are written there.
