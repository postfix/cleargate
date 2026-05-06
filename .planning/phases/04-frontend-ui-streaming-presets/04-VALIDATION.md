---
phase: 04
slug: frontend-ui-streaming-presets
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-05-06
---

# Phase 4 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | `vitest` & `react-testing-library` |
| **Config file** | `vite.config.ts` |
| **Quick run command** | `npm run test` |
| **Full suite command** | `npm run test --run` |
| **Estimated runtime** | ~10 seconds |

---

## Sampling Rate

- **After every task commit:** Run `npm run lint` and `npm run test`
- **After every plan wave:** Verify dev server boots cleanly
- **Before `/gsd-verify-work`:** Full test suite green
- **Max feedback latency:** 10 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 4-01-01 | 01 | 1 | UI-01 | unit | `npm run test -- Component` | ❌ W0 | ⬜ pending |
| 4-01-02 | 01 | 2 | UI-02 | unit | `npm run test -- Catalog` | ❌ W0 | ⬜ pending |
| 4-01-03 | 01 | 2 | EXEC-03 | e2e | Manual browser test | ✅ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] Vite project scaffolded successfully.
- [ ] Base CSS variables established from `04-UI-SPEC.md`.

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| SSE Log Streaming | EXEC-03 | Requires active backend connection | Run backend, trigger job, observe frontend log panel updates and auto-scrolls. |
| Form Preset loading | UI-03 | Visual interaction | Click preset pill, verify form inputs update instantly. |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 10s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
