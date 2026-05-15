---
phase: 09
slug: toolspec-validation-and-timeout-enforcement
status: draft
nyquist_compliant: true
wave_0_complete: false
created: 2026-05-13
---

# Phase 09 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none |
| **Quick run command** | `go test ./internal/models ./internal/api -run TestValidation` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** Run quick tests related to the modified package
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 5 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 09-01-01 | 01 | 1 | TOOL-04 | unit | `go test ./internal/models -run TestToolSpecValidation` | ✅ | ✅ green |
| 09-01-02 | 01 | 1 | TOOL-04 | unit | `go test ./internal/api/admin -run TestDraftValidation` | ✅ | ✅ green |
| 09-02-01 | 02 | 2 | TOOL-04 | unit | `go test ./internal/api -run TestUnknownFlagRejection` | ✅ | ✅ green |
| 09-02-02 | 02 | 2 | TOOL-03 | unit | `go test ./internal/api -run TestGracefulTimeout` | ✅ | ✅ green |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [x] `internal/models/toolspec_validation_test.go` — stubs for struct validation testing
- [x] `internal/api/admin/draft_validation_test.go` — stubs for validation API testing
- [x] `internal/api/execute_validation_test.go` — stubs for unknown flag rejection
- [x] `internal/api/execute_timeout_test.go` — stubs for timeout behavior

*If none: "Existing infrastructure covers all phase requirements."*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Validation UI | TOOL-04 | React form component testing requires browser interaction. | Run React dev server, trigger validation error, verify inline red highlighting. |

---

## Validation Sign-Off

- [x] All tasks have `<automated>` verify or Wave 0 dependencies
- [x] Sampling continuity: no 3 consecutive tasks without automated verify
- [x] Wave 0 covers all MISSING references
- [x] No watch-mode flags
- [x] Feedback latency < 5s
- [x] `nyquist_compliant: true` set in frontmatter

**Approval:** approved 2026-05-13

---

## Validation Audit 2026-05-15
| Metric | Count |
|--------|-------|
| Gaps found | 4 |
| Resolved | 4 |
| Escalated | 0 |
