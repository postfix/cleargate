---
wave: 1
depends_on: []
files_modified:
  - internal/llm/assistant.go
  - internal/llm/assistant_test.go
  - internal/repository/toolspec_repo.go
  - internal/repository/toolspec_repo_test.go
  - internal/api/admin/draft.go
  - internal/api/admin/draft_test.go
autonomous: true
---

# Phase 3: Tool Administration & LLM Pipeline

<objective>
Implement the administrative lifecycle for ToolSpecs, including LLM-assisted drafting via an abstracted TemplateAssistant, and persistence in DuckDB.
</objective>

<requirements>
- ADMIN-01: System provides a Versioned ToolSpec Registry stored in DuckDB.
- ADMIN-02: Tool maintainers can review, edit, and submit ToolSpecs for approval.
- ADMIN-03: System uses CloudWeGo Eino (via TemplateAssistant) to draft ToolSpecs from CLI help output.
</requirements>

<tasks>

<task>
<id>1</id>
<title>TemplateAssistant Interface and Mock</title>
<read_first>
- .planning/phases/03-tool-administration-llm-pipeline/03-RESEARCH.md
- internal/models/toolspec.go
</read_first>
<action>
1. Create `internal/llm/assistant.go`.
2. Define `type TemplateAssistant interface { GenerateDraft(ctx context.Context, helpText string) (*models.ToolSpec, error) }`.
3. Implement a `MockAssistant` that returns a dummy `models.ToolSpec` (e.g., parsing a static YAML string) for MVP/Testing purposes without needing actual LLM API keys.
4. Create `internal/llm/assistant_test.go` to verify the mock correctly returns a validated `ToolSpec`.
</action>
<acceptance_criteria>
- `internal/llm/assistant.go` contains the interface.
- `go test ./internal/llm/...` exits 0.
</acceptance_criteria>
</task>

<task>
<id>2</id>
<title>DuckDB ToolSpec Repository</title>
<depends_on>1</depends_on>
<read_first>
- internal/models/toolspec.go
</read_first>
<action>
1. Get the DuckDB driver: `go get github.com/marcboeker/go-duckdb`.
2. Create `internal/repository/toolspec_repo.go`.
3. Implement `NewToolSpecRepository(dbPath string)` that initializes the `toolspecs` table.
4. Implement `SaveDraft(spec *models.ToolSpec) error`. Serialize the struct back to YAML/JSON for the `content` column and set status to `draft`.
5. Implement `ListDrafts()`.
6. Implement `Approve(id string)`.
7. Create `internal/repository/toolspec_repo_test.go`. Use an in-memory db (`?access_mode=READ_WRITE&memory_limit=1GB`) or file-based for testing.
</action>
<acceptance_criteria>
- ToolSpecs can be saved, listed, and approved.
- `go test ./internal/repository/...` exits 0.
</acceptance_criteria>
</task>

<task>
<id>3</id>
<title>Admin API Endpoints</title>
<depends_on>2</depends_on>
<read_first>
- internal/llm/assistant.go
- internal/repository/toolspec_repo.go
</read_first>
<action>
1. Create `internal/api/admin/draft.go`.
2. Implement HTTP handlers:
   - `POST /api/admin/tools/draft` - Parses `{"help_text": "..."}`, calls `TemplateAssistant.GenerateDraft`, then calls `repo.SaveDraft`.
   - `GET /api/admin/tools/drafts` - Calls `repo.ListDrafts`.
   - `POST /api/admin/tools/{id}/approve` - Calls `repo.Approve`.
3. Create `internal/api/admin/draft_test.go` using `httptest` to mock these endpoints and verify the correct DB states and JSON responses.
</action>
<acceptance_criteria>
- API handlers correctly wire the `TemplateAssistant` and `ToolSpecRepository`.
- `go test ./internal/api/admin/...` exits 0.
</acceptance_criteria>
</task>

</tasks>

<verification>
- Automated tests pass for DuckDB repository persistence.
- Eino abstraction correctly generates drafts and API handles them securely.
</verification>

<must_haves>
- The `TemplateAssistant` interface must completely shield the core logic from CloudWeGo Eino details.
- ToolSpecs must be stored with explicit statuses (`draft`, `approved`).
</must_haves>
