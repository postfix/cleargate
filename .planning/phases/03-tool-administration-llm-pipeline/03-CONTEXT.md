# Phase 3: Tool Administration & LLM Pipeline - Context

**Gathered:** 2026-05-06
**Status:** Ready for planning
**Mode:** Auto-Chained (Agent-selected defaults)

<domain>
## Phase Boundary

Integration with CloudWeGo Eino for LLM-assisted ToolSpec generation and administrative workflows (draft, review, approve, activate).
</domain>

<decisions>
## Implementation Decisions

### Eino Integration Pattern
- **D-01:** Implement a `TemplateAssistant` interface internally that wraps the Eino API components. This keeps the core execution logic entirely unaware of LLMs.

### Draft ToolSpec Storage
- **D-02:** Store draft ToolSpecs in the metadata database (DuckDB) with explicit status fields (e.g., `draft`, `approved`, `rejected`), ensuring they are completely separated from the active execution path until approved.

### Tool Discovery Engine Input
- **D-03:** For MVP, accept raw `--help` text via an API endpoint for the LLM pipeline rather than auto-executing introspection commands, to maintain a strict security boundary.

### the agent's Discretion
- Agent selected options autonomously via `--chain` flag.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Specification
- `SPEC.md` — Section 8 defines the Eino Integration Model explicitly.
- `SPEC.md` — Section 9 outlines Trust Boundaries for LLM outputs.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/models/toolspec.go` — Target schema the LLM must generate.
- `internal/validation/validator.go` — Reusable for validating LLM drafts.

### Integration Points
- New API package `internal/api/admin` or `internal/api/discovery`.
- DuckDB repository layer for ToolSpecs.

</code_context>
