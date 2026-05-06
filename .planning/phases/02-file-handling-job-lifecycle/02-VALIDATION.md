---
phase: 02
slug: file-handling-job-lifecycle
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-05-06
---

# Phase 2 — Validation Strategy

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
| 2-01-01 | 01 | 1 | FILE-01 | unit | `go test ./internal/workspace/...` | ❌ W0 | ⬜ pending |
| 2-01-02 | 01 | 1 | FILE-01 | unit | `go test ./internal/api/...` | ❌ W0 | ⬜ pending |
| 2-01-03 | 01 | 2 | EXEC-02 | unit | `go test ./internal/runtime/podman/...` | ✅ W0 | ⬜ pending |
| 2-01-04 | 01 | 2 | FILE-02 | unit | `go test ./internal/api/...` | ✅ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/workspace/workspace_test.go` — stubs for testing workspace creation
- [ ] `internal/api/upload_test.go` — stubs for testing multipart upload
- [ ] `internal/api/download_test.go` — stubs for testing artifact download

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Test API via curl | FILE-01 | Validate end-to-end HTTP request | Start server, upload a file via `curl -F`, verify file on disk. |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 10s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
