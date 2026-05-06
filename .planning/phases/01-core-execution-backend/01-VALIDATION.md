---
phase: 01
slug: core-execution-backend
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-05-06
---

# Phase 1 — Validation Strategy

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
| 1-01-01 | 01 | 1 | EXEC-05 | unit | `go test ./internal/models/...` | ❌ W0 | ⬜ pending |
| 1-01-02 | 01 | 1 | EXEC-04 | unit | `go test ./internal/command/...` | ❌ W0 | ⬜ pending |
| 1-01-03 | 01 | 1 | EXEC-01 | unit | `go test ./internal/runtime/podman/...` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/models/toolspec_test.go` — stubs for testing YAML parsing
- [ ] `internal/command/builder_test.go` — stubs for testing argv generation
- [ ] `internal/runtime/podman/client_test.go` — stubs for ContainerRuntime testing

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Execute container via podman socket | EXEC-01 | Requires active socket / environment | Run a simple test locally ensuring podman starts the container |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 10s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
