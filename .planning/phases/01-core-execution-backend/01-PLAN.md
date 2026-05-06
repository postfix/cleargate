---
wave: 1
depends_on: []
files_modified:
  - go.mod
  - internal/models/toolspec.go
  - internal/models/toolspec_test.go
  - internal/validation/validator.go
  - internal/validation/validator_test.go
  - internal/command/builder.go
  - internal/command/builder_test.go
  - internal/runtime/runtime.go
  - internal/runtime/podman/client.go
  - internal/runtime/podman/client_test.go
autonomous: true
---

# Phase 1: Core Execution Backend

<objective>
Implement the foundational Go backend to parse ToolSpecs, validate requests, construct arguments safely, and execute containers via Podman bindings.
</objective>

<requirements>
- EXEC-01: User can run approved CLI tools in a sandboxed execution environment (Docker/Podman).
- EXEC-04: Backend compiles validated job values into `argv[]` without using shell execution.
- EXEC-05: System prevents arbitrary command execution by using structured ToolSpecs.
- TOOL-04: ToolSpec schemas are strongly validated.
</requirements>

<tasks>

<task>
<id>1</id>
<title>Initialize Go Module and ToolSpec Models</title>
<read_first>
- .planning/phases/01-core-execution-backend/01-CONTEXT.md
- SPEC.md
</read_first>
<action>
1. Initialize Go module: `go mod init github.com/postfix/cleargate` if not exists.
2. Get yaml dependency: `go get gopkg.in/yaml.v3`
3. Create `internal/models/toolspec.go` with Go structs to represent the ToolSpec schema matching SPEC.md.
   Include strictly-typed structures for `ToolSpec`, `Metadata`, `Runtime`, `Input`, `Flag` (with `CliRender` and `UiRender`), `Positional`, `Output`, `Preset`, `SecurityPolicy`.
4. Create `internal/models/toolspec_test.go` and add a unit test parsing a valid pandoc yaml spec string to ensure fields map correctly.
</action>
<acceptance_criteria>
- `go.mod` contains `module github.com/postfix/cleargate`
- `internal/models/toolspec.go` contains `type ToolSpec struct`
- `go test ./internal/models` exits 0
</acceptance_criteria>
</task>

<task>
<id>2</id>
<title>Implement Input Validation</title>
<depends_on>1</depends_on>
<read_first>
- internal/models/toolspec.go
</read_first>
<action>
1. Create `internal/validation/validator.go`.
2. Implement `func ValidateJobValues(spec *models.ToolSpec, values map[string]interface{}) error`.
3. The function must iterate through `spec.Flags`, check if required flags are present, and validate types for provided values (string, bool, enum).
4. Create `internal/validation/validator_test.go` verifying valid values pass and unknown flags or wrong types return an error.
</action>
<acceptance_criteria>
- `internal/validation/validator.go` contains `func ValidateJobValues`
- `go test ./internal/validation` exits 0
</acceptance_criteria>
</task>

<task>
<id>3</id>
<title>Implement Command Builder</title>
<depends_on>2</depends_on>
<read_first>
- internal/models/toolspec.go
- SPEC.md
</read_first>
<action>
1. Create `internal/command/builder.go`.
2. Implement `func BuildCommand(spec *models.ToolSpec, values map[string]interface{}) ([]string, error)`.
3. The builder must start the array with `spec.Runtime.Argv0`.
4. Then append rendered flags based on `cli.render` rules (`whenTrue`, `sequence`, `keyValue`).
5. Then append positionals based on `order`.
6. Ensure no shell concatenation happens; return a clean `[]string`.
7. Create `internal/command/builder_test.go` covering boolean flag inclusion and sequence formatting.
</action>
<acceptance_criteria>
- `internal/command/builder.go` contains `func BuildCommand` returning `[]string, error`
- `go test ./internal/command` exits 0
</acceptance_criteria>
</task>

<task>
<id>4</id>
<title>Implement Podman ContainerRuntime</title>
<depends_on>1</depends_on>
<read_first>
- .planning/phases/01-core-execution-backend/01-CONTEXT.md
</read_first>
<action>
1. Get podman dependency: `go get go.podman.io/podman/v6/pkg/bindings` and `go.podman.io/podman/v6/pkg/specgen`
2. Create `internal/runtime/runtime.go` with interface `ContainerRuntime` defining `PullImage`, `Create`, `Start`, `Wait`, `Inspect`.
3. Create `internal/runtime/podman/client.go` that implements `ContainerRuntime` connecting to `unix:///run/user/$UID/podman/podman.sock` (or `XDG_RUNTIME_DIR`). Ensure build tag `//go:build remote` is documented or handled.
4. Implement the methods using the official bindings.
5. Create `internal/runtime/podman/client_test.go` with basic structural checks or skip actual podman execution if socket is absent during tests.
</action>
<acceptance_criteria>
- `internal/runtime/runtime.go` contains `type ContainerRuntime interface`
- `internal/runtime/podman/client.go` uses `go.podman.io/podman/v6/pkg/bindings`
- `go test ./internal/runtime/podman` exits 0 or skips cleanly
</acceptance_criteria>
</task>

</tasks>

<verification>
- Automated unit tests for toolspec loading, validation, and command building all pass.
- Podman bindings fetch successfully and compile.
</verification>

<must_haves>
- Uses typed Go structs for ToolSpec parsing.
- Uses exact Podman v6 bindings for the container interface.
- Returns a strict `[]string` for `argv` command construction.
</must_haves>
