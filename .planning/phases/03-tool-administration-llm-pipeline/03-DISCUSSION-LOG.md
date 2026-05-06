# Phase 3: Tool Administration & LLM Pipeline - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

*Note: This phase was processed via `--chain`. The agent made architectural decisions autonomously based on `SPEC.md` constraints.*

## Q1: Eino Integration Pattern
Options:
- Tightly couple Eino orchestration in the admin handlers.
- Abstract behind a `TemplateAssistant` interface.
Selected: **Abstract behind a `TemplateAssistant` interface.** (Aligns with SPEC.md Section 8)

## Q2: Draft ToolSpec Storage
Options:
- Store as temporary files.
- Store in DuckDB with a `status` field.
Selected: **Store in DuckDB with a `status` field.**

## Q3: Tool Discovery Engine Input
Options:
- Auto-execute arbitrary commands to fetch help.
- Accept raw `--help` text via API payload.
Selected: **Accept raw `--help` text via API payload.** (Safer for MVP)
