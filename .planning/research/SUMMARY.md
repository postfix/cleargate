# Research Summary: Rootless Podman Hardening & First Application

## Stack Additions
- **Podman CLI / Bindings**: To execute rootless containers, `os/exec` calling the `podman` binary is highly predictable and easier to configure for rootless subuid/subgid mapping than the Go bindings, which can suffer from CGO header issues (as encountered with `gpgme` on macOS).
- **ToolSpec Seeding**: The backend needs a mechanism to seed real YAML ToolSpecs into DuckDB on startup, removing the need for UI mock fallbacks.

## Feature Table Stakes
- **Rootless Execution**: Jobs must execute under the user's namespace without requiring root daemon privileges.
- **Security Profile (Hardening)**:
  - `--cap-drop=all`: Drop all Linux capabilities.
  - `--security-opt no-new-privileges`: Prevent privilege escalation within the container.
  - `--read-only`: Root filesystem must be read-only, with only the workspace mounted as writable.
  - `--network=none`: Disconnect networking by default (unless specifically requested in the ToolSpec).
- **First Real Application**: Introduce a fully working `nmap` or `ffmpeg` ToolSpec in the backend database.

## Watch Out For (Pitfalls)
- **UID/GID Mapping**: Rootless Podman maps the host user to root (UID 0) inside the container. The workspace directory mounted into the container must have permissions that align with the subuid mapping, otherwise the tool will fail with "Permission denied" when writing outputs.
- **Mock Fallback Confusion**: The UI currently falls back to a hardcoded FFmpeg schema if the backend route fails or returns empty. This MUST be removed or explicitly labeled to prevent user confusion.
- **Mac Development**: Mac uses a Podman machine (VM). The workspace paths mounted from the Mac host must be accessible to the Podman VM, which requires proper volume mounting setup.
