---
wave: 1
depends_on: []
files_modified:
  - internal/runtime/podman.go
  - go.mod
  - Makefile
autonomous: true
---

# Plan 01: Podman Bindings & Hardened Sandbox

<objective>
Implement the real Podman runtime using `github.com/containers/podman/v4/pkg/bindings`. Enforce the strict rootless security profile (capabilities dropped, read-only rootfs) via the bindings API, while allowing the environment to auto-resolve the socket connection on both Mac and Linux.
</objective>

<requirements>
- EXEC-06: Real Podman Integration
- EXEC-07: Security Profile
</requirements>

<tasks>

## 1. Add Podman Bindings and Build Constraints
Implement the `ContainerRuntime` interface using the official Podman bindings.

<read_first>
- `cmd/cleargate/main.go` (to see the `DummyRuntime` we are replacing)
- `internal/runtime/runtime.go` (if exists, or where the interface is defined)
</read_first>

<action>
1. Run `go get github.com/containers/podman/v4/pkg/bindings`
2. Create `internal/runtime/podman.go`.
3. Implement `NewPodmanRuntime() (*PodmanRuntime, error)`.
4. In the constructor, resolve the connection URI:
   - Read `os.Getenv("CONTAINER_HOST")`. If set, use it.
   - If unset, check `runtime.GOOS`. If `darwin`, fallback to `unix:///Users/` + `os.Getenv("USER")` + `/.local/share/containers/podman/machine/podman.sock`.
   - If `linux`, fallback to `unix://` + `os.Getenv("XDG_RUNTIME_DIR")` + `/podman/podman.sock`.
   - Call `bindings.NewConnection(ctx, uri)`.
5. Create a `Makefile` (or update existing) with a `build` target: `go build -tags exclude_graphdriver_btrfs,btrfs_noversion,containers_image_openpgp ./cmd/cleargate` to ensure it compiles on macOS without requiring `gpgme` headers natively.
</action>

<acceptance_criteria>
- `go.mod` contains `github.com/containers/podman/v4`
- `internal/runtime/podman.go` contains `NewPodmanRuntime` and the URI fallback logic.
- `Makefile` exists and contains the `-tags` flags to bypass `gpgme`/`btrfs` headers.
- Code compiles cleanly on macOS using the `Makefile` build command.
</acceptance_criteria>

## 2. Implement Hardened Container Execution
Map the ToolSpec execution logic to the Podman bindings, enforcing the security profile.

<read_first>
- `internal/runtime/podman.go`
</read_first>

<action>
1. Implement the `Create` method in `PodmanRuntime`.
2. Construct the `specgen.SpecGenerator` object.
3. Apply default hardening:
   - `CapDrop = []string{"ALL"}`
   - `ReadOnlyFilesystem = true` (using a boolean pointer)
   - `NoNewPrivileges = true`
4. If the ToolSpec has security overrides (e.g., `spec.Runtime.Capabilities`), append them to `CapAdd`.
5. Bind mount the job workspace directory into the container. Ensure the mount has the correct options (e.g., `Z` for SELinux if necessary, or just standard read/write).
6. Implement `Start`, `Wait`, and `Logs` methods using the respective `bindings/containers` packages.
</action>

<acceptance_criteria>
- `internal/runtime/podman.go` contains `containers.CreateWithSpec` using `specgen`.
- The `SpecGenerator` explicitly sets `CapDrop` to `ALL`, `ReadOnlyFilesystem` to true, and `NoNewPrivileges` to true.
</acceptance_criteria>

</tasks>

<verification>
- A test program or manual run using the `Makefile` builds the binary successfully.
- Triggering an execution attempts to connect to the Podman socket and create a container with the strict security flags.
</verification>
