---
requirements_completed:
  - TOOL-03
---

# Plan 01: ToolSpec Validation & Versioning Summary

## What was done
1. Added `go-playground/validator/v10` struct validation tags to `internal/models/toolspec.go` to enforce requirements (e.g., `TimeoutSeconds` required and `min=1`).
2. Updated `internal/api/admin/draft.go` to return HTTP 400 Bad Request with a structured JSON mapping of field namespaces to exact tags on validation errors.
3. Updated `internal/repository/toolspec_repo.go` to enforce an append-only versioning strategy, using `ON CONFLICT DO NOTHING` to prevent replacing existing, already-approved versions of ToolSpecs.
4. Fixed the `draft_test.go` to inject the proper mock assistant struct and validation test fields.

## Issues encountered
- The mock template generated in tests was initially missing the required struct fields that were added, which caused the tests to fail. This was fixed by expanding the LLM mock string to include required fields like `runtime.timeoutSeconds`.
- `SaveDraft` was previously replacing drafts of the same version due to `ON CONFLICT DO UPDATE`. This was changed to `DO NOTHING` to ensure a true append-only history without primary key panics.

## Next steps
Proceeding to Wave 2 to enforce unknown flag rejection and wire the validation payloads up to the React frontend.
