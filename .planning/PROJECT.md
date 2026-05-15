# ClearGate

## What This Is

ClearGate is a secure CLI application gateway that converts approved command-line tools into generated web applications. It allows users to run approved tools in a sandboxed execution environment by configuring command-line flags through a generated SPA instead of a terminal.

## Core Value

Expose CLI tools safely through generated web interfaces using strict, deterministic ToolSpec validation to prevent arbitrary command execution.

## Current Milestone: v2.0 LLM Assistant & WebSocket Streaming

**Goal:** Implement CloudWeGo Eino integration to automatically draft ToolSpecs from raw CLI manuals and upgrade log streaming to WebSockets for broader browser compatibility.

## Requirements

### Validated

<!-- Shipped and confirmed valuable. -->

- ✓ Expose CLI tools safely through generated web interfaces. — v1.0
- ✓ Prevent arbitrary command execution by using structured ToolSpecs and strict validation. — v1.1
- ✓ Support file upload and artifact download for tools that process files. — v1.0
- ✓ Capture stdout, stderr, exit code, metadata, and output files for every run. — v1.0
- ✓ Allow users and teams to save presets as reusable execution profiles. — v1.0
- ✓ Provide strong auditability (who, what, when, inputs, versions, sandbox policy). — v1.0
- ✓ Generate UI automatically from the ToolSpec. — v1.0
- ✓ Implement secure execution backend with sandbox policies (Docker/Podman). — v1.0
- ✓ Rootless Podman hardened sandbox profile. — v1.1
- ✓ Provide a fully working first application (ToolSpec). — v1.1
- ✓ Support local-first and enterprise deployments. — v1.0

### Active

<!-- Current scope. Building toward these. -->

- [ ] Use LLMs to accelerate ToolSpec template creation (via CloudWeGo Eino).
- [ ] WebSocket-based streaming (v1 uses SSE).
- [ ] Maintainer can run discovery on a binary or container to extract help/man/docs.

### Out of Scope

<!-- Explicit boundaries. Includes reasoning to prevent re-adding. -->

- [Browser-based terminal or generic remote shell] — ClearGate is for structured execution, not free-form shell access.
- [Replacement for Kubernetes jobs or CI/CD systems] — Designed for interactive, user-driven CLI tool execution.
- [LLM agent deciding commands on its own] — LLMs are only used for drafting templates, not runtime command generation.
- [Exposing all CLIs automatically] — Requires human review and explicit ToolSpec approval for security.

## Context

Shipped v1.1 with ~3,500 LOC Go and React.
Tech stack: Go, React, DuckDB, Podman bindings.
Current application enforces strict schema validation and runs applications in rootless Podman sandboxes with streaming.

## Constraints

- **Language**: Go (Backend), React (Frontend) — Required by spec.
- **LLM Orchestration**: CloudWeGo Eino — Required for LLM integration.
- **Database**: DuckDB — For metadata storage.
- **Sandbox**: Docker/Podman — Required for isolated job execution.
- **Execution**: Backend must never invoke `sh -c`, commands must be constructed as `argv[]` arrays.

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| ToolSpec as Source of Truth | Ensures deterministic behavior, security, and UI consistency | ✓ Good |
| Untrusted LLM Pipeline | LLMs hallucinate; restricting them to drafting templates ensures safety | — Pending |
| Append-only ToolSpec Versioning | Preserves historical auditability without breaking old presets | ✓ Good |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-05-15 after v1.1 milestone*
