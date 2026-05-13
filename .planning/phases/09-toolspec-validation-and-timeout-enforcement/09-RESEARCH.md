# Phase 09: Technical Research

## Objective
Research how to implement Phase 09: ToolSpec Validation and Timeout Enforcement, fulfilling the context decisions.

## 1. Validation Engine (`go-playground/validator/v10`)
**Pattern Identified:** 
- The `validator/v10` package uses struct tags (e.g., `validate:"required"`) to define rules.
- To provide inline form validation errors (D-05), the backend needs to catch `validator.ValidationErrors`, map each `FieldError.StructNamespace()` or `FieldError.Field()` to the corresponding JSON key, and return a structured JSON response (e.g., `{"errors": {"field_name": "error message"}}`).
- We can add standard tags like `validate:"required,gt=0"` to the `ToolSpec` struct fields in `internal/models/toolspec.go`.
- The `SyncFromDirectory` seeder and `POST /api/admin/drafts` endpoints should both call `validator.Struct(spec)`.

## 2. Timeout Enforcement Behavior (Podman API)
**Pattern Identified:**
- `context.WithTimeout` works well for HTTP requests, but when wrapping a Podman container lifecycle, we want graceful termination (D-03).
- Instead of relying solely on `context.WithTimeout` around `containers.Wait` (which doesn't automatically kill the container gracefully), we need a select loop or a dedicated goroutine that monitors a `context.Context`.
- If the timeout context is `Done()`, the system should call `containers.Kill` with `signal="SIGTERM"`.
- It should then wait up to 5 seconds. If the container hasn't stopped, it should call `containers.Kill` with `signal="SIGKILL"`.
- *Constraint:* Podman bindings `containers.Kill` function signature: `containers.Kill(ctx, nameOrID, &containers.KillOptions{Signal: pointer("SIGTERM")})`.

## 3. Unknown Flag Rejection
**Pattern Identified:**
- The execution endpoint `POST /api/execute` receives a JSON map of inputs.
- To fail fast (D-04), the `ExecutionHandler` must decode the JSON into a `map[string]interface{}` (or similar) and iterate over all provided keys.
- For each key, it checks if it exists in `ToolSpec.Flags` or `ToolSpec.Inputs`. If not, it immediately returns a `400 Bad Request` with an explicit message: `"Unknown flag provided: {key}"`.

## 4. UI Validation Feedback
**Pattern Identified:**
- The React form components (likely using standard state or `react-hook-form`) need to parse the `400 Bad Request` JSON payload.
- The UI should maintain an `errors` state object mapping field IDs to error messages.
- The inputs should conditionally render a red border or an error span below them if their ID exists in the `errors` object.

## Validation Architecture
- **Validation Stage 1 (Static):** Validate `ToolSpec` structs using `validator/v10` before saving.
- **Validation Stage 2 (Dynamic):** Validate execution payloads against the specific `ToolSpec` allowed keys, rejecting unknown flags.
- **Execution Stage (Enforcement):** Manage process timeouts with `SIGTERM` -> `5s` -> `SIGKILL` sequence.

## RESEARCH COMPLETE
