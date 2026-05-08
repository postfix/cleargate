# ClearGate

## What This Is

ClearGate is a secure CLI application gateway that converts approved command-line tools into generated web applications. It allows users to run approved tools in a sandboxed execution environment by configuring command-line flags through a generated SPA instead of a terminal.

## Core Value

Expose CLI tools safely through generated web interfaces using strict, deterministic ToolSpec validation to prevent arbitrary command execution.

## Current Milestone: v1.1 Rootless Podman Hardening & First Application

**Goal:** Implement a fully rootless, hardened Podman sandbox profile and deploy the first fully working CLI application through ClearGate.

**Target features:**
- Implement rootless Podman execution for job isolation
- Create hardened sandbox security profiles (seccomp, network isolation)
- Define and approve a real, fully working ToolSpec application
- Verify end-to-end execution of this first real application

## Requirements

### Validated

<!-- Shipped and confirmed valuable. -->

(None yet — ship to validate)

### Active

<!-- Current scope. Building toward these. -->

- [ ] Expose CLI tools safely through generated web interfaces.
- [ ] Prevent arbitrary command execution by using structured ToolSpecs and strict validation.
- [ ] Support file upload and artifact download for tools that process files.
- [ ] Capture stdout, stderr, exit code, metadata, and output files for every run.
- [ ] Allow users and teams to save presets as reusable execution profiles.
- [ ] Use LLMs to accelerate ToolSpec template creation (via CloudWeGo Eino).
- [ ] Support local-first and enterprise deployments.
- [ ] Provide strong auditability (who, what, when, inputs, versions, sandbox policy).
- [ ] Generate UI automatically from the ToolSpec.
- [ ] Implement secure execution backend with sandbox policies (Docker/Podman).
- [ ] Rootless Podman hardened sandbox profile.
- [ ] Provide a fully working first application (ToolSpec).

### Out of Scope

<!-- Explicit boundaries. Includes reasoning to prevent re-adding. -->

- [Browser-based terminal or generic remote shell] — ClearGate is for structured execution, not free-form shell access.
- [Replacement for Kubernetes jobs or CI/CD systems] — Designed for interactive, user-driven CLI tool execution.
- [LLM agent deciding commands on its own] — LLMs are only used for drafting templates, not runtime command generation.
- [Exposing all CLIs automatically] — Requires human review and explicit ToolSpec approval for security.

## Context

- Enterprise environments have many useful CLI tools, but giving users shell access is risky and hard to audit.
- ClearGate solves this by making the ToolSpec the source of truth for UI, validation, command construction, and security.
- The LLM output is untrusted and must pass schema validation, policy checks, and human approval before becoming an active ToolSpec.

## Constraints

- **Language**: Go (Backend), React (Frontend) — Required by spec.
- **LLM Orchestration**: CloudWeGo Eino — Required for LLM integration.
- **Database**: DuckDB — For metadata storage.
- **Sandbox**: Docker/Podman — Required for isolated job execution.
- **Execution**: Backend must never invoke `sh -c`, commands must be constructed as `argv[]` arrays.

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| ToolSpec as Source of Truth | Ensures deterministic behavior, security, and UI consistency | — Pending |
| Untrusted LLM Pipeline | LLMs hallucinate; restricting them to drafting templates ensures safety | — Pending |

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
*Last updated: 2026-05-06 after initialization*
