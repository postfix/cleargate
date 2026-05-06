---
phase: 03
slug: tool-administration-llm-pipeline
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-05-06
---

# Phase 3 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | `go test` |
| **Config file** | `go.mod` |
| **Quick run command** | `go test -short ./...` |
| **Full suite command** | `go test -v -race ./...` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test -short ./...`
- **After every plan wave:** Run `go test -v -race ./...`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 10 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 3-01-01 | 01 | 1 | ADMIN-03 | unit | `go test ./internal/llm/...` | ❌ W0 | ⬜ pending |
| 3-01-02 | 01 | 1 | ADMIN-01 | unit | `go test ./internal/repository/...` | ❌ W0 | ⬜ pending |
| 3-01-03 | 01 | 2 | ADMIN-02 | unit | `go test ./internal/api/admin/...` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/llm/assistant_test.go` — stubs for testing the interface
- [ ] `internal/repository/toolspec_repo_test.go` — stubs for testing DuckDB
- [ ] `internal/api/admin/draft_test.go` — stubs for testing endpoints

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 10s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
