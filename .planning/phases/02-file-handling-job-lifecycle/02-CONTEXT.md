# Phase 2: File Handling & Job Lifecycle - Context

**Gathered:** 2026-05-06
**Status:** Ready for planning

<domain>
## Phase Boundary

Support job inputs/outputs, artifacts, and capture job results. Focuses on file uploads, workspace structure, execution tracking, and artifact retrieval.
</domain>

<decisions>
## Implementation Decisions

### Workspace Storage Strategy
- **D-01:** Use a persistent data directory (e.g., `/var/lib/cleargate/jobs/{job_id}`) that maps directly to the `SPEC.md` layout (input, output, logs, metadata).

### File Upload Mechanism
- **D-02:** Stream multipart uploads directly to the job's `input/` folder using safe, generated names to avoid memory crashes on large files.

### Output Artifact Retrieval
- **D-03:** Serve declared artifacts directly from the `output/` path over HTTP after validation.

### Log Capture (Stdout/Stderr)
- **D-04:** Stream to `logs/stdout.log` and `logs/stderr.log` directly in the job workspace alongside execution.

### the agent's Discretion
None.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Specification
- `SPEC.md` — Defines workspace layout and artifact rules.
- `.planning/phases/01-core-execution-backend/01-CONTEXT.md` — Podman execution foundation context.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/runtime/podman/client.go` — Podman execution client.

### Established Patterns
- Podman rootless socket connection.
- `ContainerRuntime` interface abstraction.

### Integration Points
- Backend server HTTP endpoints for file upload and download.
- Job lifecycle manager orchestrating the ContainerRuntime and file system paths.

</code_context>

<specifics>
## Specific Ideas

Follow the directory layout specified in SPEC.md exactly:
`/jobs/{job_id}/input/`, `/jobs/{job_id}/output/`, `/jobs/{job_id}/logs/`, `/jobs/{job_id}/metadata/`

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope

</deferred>

---

*Phase: 02-file-handling-job-lifecycle*
*Context gathered: 2026-05-06*
