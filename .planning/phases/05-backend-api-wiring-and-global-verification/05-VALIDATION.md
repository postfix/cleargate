---
phase: 05
slug: backend-api-wiring-and-global-verification
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-05-07
---

# Phase 05 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none |
| **Quick run command** | `go test ./internal/api...` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/api...`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 5 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 05-01-01 | 01 | 1 | EXEC-01 | unit | `go test ./internal/api` | ❌ W0 | ⬜ pending |
| 05-01-02 | 01 | 1 | UI-02 | unit | `go test ./internal/api` | ❌ W0 | ⬜ pending |
| 05-01-03 | 01 | 1 | PRESET-01 | unit | `go test ./internal/api` | ❌ W0 | ⬜ pending |
| 05-02-01 | 02 | 2 | ALL | build | `go build ./cmd/cleargate` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `cmd/cleargate/main.go` — stub for main server
- [ ] `internal/api/execute_test.go` — stubs for execution
- [ ] `internal/api/catalog_test.go` — stubs for catalog
- [ ] `internal/api/preset_test.go` — stubs for presets

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| SPA Serving | UI-01 | Browser interaction | Run the built binary and access localhost in browser |
| SSE Streaming | EXEC-03 | Async UI events | Run a long job and verify the logs appear in the UI |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 5s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
