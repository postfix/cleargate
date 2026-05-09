# Phase 8: Persistent Presets and Audit Logging - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-05-09
**Phase:** 08-persistent-presets-and-audit-logging
**Areas discussed:** Preset Scoping, Audit Data Shape, Audit Queryability, Preset Lifecycle

---

## Preset Scoping

| Option | Description | Selected |
|--------|-------------|----------|
| Per-tool | Presets stored with a tool_id column. The nmap page only shows nmap presets. | ✓ |
| Global | Keep current behavior. Simpler schema, but you'd see irrelevant presets on every tool page. | |

**User's choice:** Per-tool
**Notes:** 

---

## Audit Data Shape

| Option | Description | Selected |
|--------|-------------|----------|
| Minimal | Just job_id, tool_id, timestamp, and exit_code. | ✓ |
| Detailed | Adds input_values, toolspec_version, and duration_ms. | |

**User's choice:** Minimal
**Notes:** 

---

## Audit Queryability

| Option | Description | Selected |
|--------|-------------|----------|
| Backend only | Just write-only to DuckDB. | |
| API endpoint only | Expose GET /api/admin/audit, no UI. | |
| Full UI | Add an "Audit Log" tab to the frontend React app. | ✓ |

**User's choice:** Full UI
**Notes:** 

---

## Preset Lifecycle

| Option | Description | Selected |
|--------|-------------|----------|
| Save-only | Good enough for MVP. | |
| Full CRUD | Allow deleting and updating existing presets from the UI. | ✓ |

**User's choice:** Full CRUD
**Notes:** 

---
