# Phase 1: Core Execution Backend - Research

## Objective
Establish the foundational backend architecture for launching validated CLI jobs without shell execution. This research provides implementation patterns for ToolSpec parsing, validation, command generation, and sandboxing using the official Podman Go bindings.

## ToolSpec Parsing
- **Pattern:** Unmarshal YAML into strictly-typed Go structs using `gopkg.in/yaml.v3`.
- **Implementation:** Define a `models.ToolSpec` struct. Any unknown fields should trigger an error (e.g., `yaml.Unmarshal` with `KnownFields(true)` if configured via a custom unmarshaler, or manual validation).

## Validation Strategy
- **Pattern:** Validate at the API boundary before any backend processing starts.
- **Implementation:** Create a `validator` package. Given an incoming JSON request containing job values, cross-reference the keys and types against the loaded `models.ToolSpec`. Reject unknown fields, invalid types, and disabled flags with a 400 Bad Request.

## Sandboxing Interface (Podman Go Bindings)
- **Pattern:** Define a `ContainerRuntime` interface to abstract the execution details, avoiding tight coupling to Podman's specific API across the entire codebase.
- **Implementation:**
  - Import `go.podman.io/podman/v6/pkg/bindings`
  - Connect to rootless socket via `bindings.NewConnection(ctx, "unix:///run/user/$UID/podman/podman.sock")`
  - Define `ContainerRuntime` interface:
    ```go
    type ContainerRuntime interface {
      Create(ctx context.Context, req CreateContainerRequest) (ContainerID, error)
      Start(ctx context.Context, id ContainerID) error
      Wait(ctx context.Context, id ContainerID) error
      Logs(ctx context.Context, id ContainerID) (<-chan LogEvent, error)
    }
    ```
  - Create `internal/runtime/podman/client.go` to implement this using `images.Pull`, `containers.CreateWithSpec`, `containers.Start`, `containers.Wait`.
  - Use `go build -tags remote` to minimize binary size.

## Command Construction Mapping
- **Pattern:** Isolate `argv[]` builder into a pure function `BuildCommand(spec ToolSpec, values map[string]interface{}) ([]string, error)`.
- **Implementation:** 
  - Iterate through `spec.Flags` and `spec.Positionals`.
  - Compile the `cli.render` rules.
  - Return an array of strings, ensuring no shell interpolation occurs.

## Validation Architecture (Nyquist)
- Code must use standard Go testing framework.
- Validation should test the `ContainerRuntime` interface with a mock implementation.
- Real e2e tests should run against an actual rootless Podman socket if available.
