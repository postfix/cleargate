# Phase 3: Tool Administration & LLM Pipeline - Research

## Objective
Establish the admin workflows for drafting, reviewing, and approving ToolSpecs, integrating with CloudWeGo Eino for LLM-assisted generation.

## Eino Integration Pattern
- **Pattern:** Create a `TemplateAssistant` interface in `internal/llm/assistant.go` that decouples the REST handlers from the Eino library.
- **Implementation:** 
  - `GenerateDraft(ctx context.Context, helpText string) (*models.ToolSpec, error)`
  - Provide a mock implementation for MVP testing (`MockAssistant`) to avoid requiring live LLM keys during core functionality tests.
  - When real Eino is used, the response must be unmarshaled and strictly validated using `internal/validation/validator.go` logic before being accepted as a draft.

## Draft ToolSpec Storage (DuckDB)
- **Pattern:** Use `database/sql` with the `github.com/marcboeker/go-duckdb` driver.
- **Implementation:**
  - Create `internal/repository/toolspec_repo.go`.
  - Schema: `CREATE TABLE IF NOT EXISTS toolspecs (id VARCHAR PRIMARY KEY, name VARCHAR, version VARCHAR, status VARCHAR, content TEXT, created_at TIMESTAMP)`
  - Statuses: `draft`, `approved`.

## API Endpoints
- **Pattern:** Standard Go `net/http` handlers.
- **Implementation:**
  - `POST /api/admin/tools/draft` - Accepts `{"help_text": "..."}` -> Calls `TemplateAssistant` -> Saves to DuckDB as `draft`.
  - `GET /api/admin/tools/drafts` - Lists drafts.
  - `POST /api/admin/tools/{id}/approve` - Changes status to `approved`.

## Validation Architecture
- Eino outputs MUST pass the `yaml.Unmarshal` with `KnownFields(true)` check using `models.ToolSpec` created in Phase 1.
- Use `httptest` to mock the `/api/admin/tools/draft` endpoint.
