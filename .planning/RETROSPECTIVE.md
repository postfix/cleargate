# Retrospective: ClearGate

## Cross-Milestone Trends

| Metric | v1.0 | v1.1 |
|--------|------|------|
| Phases Shipped | 5 | 7 |
| Requirements | ~10 | 19 |
| Bugs/Gaps Found | 0 | 4 |

## Milestone: v1.1 — Security and Sandbox Hardening

**Shipped:** 2026-05-15
**Phases:** 7 | **Plans:** 14

### What Was Built
- Rootless execution with Podman
- Strict schema validation and ToolSpec versioning
- File handling, upload mapping and downloading artifacts
- Preset saving and audit logs
- Live SSE log streaming and backend execution wiring

### What Worked
- Separation of phases based on concerns (e.g. core wiring separate from validation).
- Adding "gap closure" phases directly into the roadmap to quickly resolve testing debts.

### What Was Inefficient
- Not writing summaries and verifications right away for phases 6, 7, and 10 left gaps that had to be retroactively resolved via gap-closure phases (11, 12) or ignored during audit.
- Some structural decisions with dummy interfaces broke during test validations due to forgotten syncs with real implementations.

### Patterns Established
- Use struct tags for validation along with `go-playground/validator/v10`.
- Return structured JSON error maps instead of single string errors from backend to easily integrate with React SPA dynamic form fields.

### Key Lessons
- E2E testing using `go test ./...` with mock db logic needs to thoroughly assert all implementations of an interface, such as `GracefulStop` missing from `DummyRuntime`.

### Cost Observations
- Notable: Fast iterations using `go test` and `make build` with build tags for Podman APIs.
