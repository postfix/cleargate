# Phase 2: File Handling & Job Lifecycle - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

## Q1: Workspace Storage Strategy
Options:
- Use a persistent data directory mapping to SPEC.md layout
- Use ephemeral `/tmp` storage
Selected: **Use a persistent data directory mapping to SPEC.md layout**

## Q2: File Upload Mechanism
Options:
- Stream multipart uploads directly to disk
- Buffer in memory first
Selected: **Stream multipart uploads directly to disk**

## Q3: Output Artifact Retrieval
Options:
- Serve declared artifacts directly from the output/ path
- Dynamically zip all contents of output/
Selected: **Serve declared artifacts directly from the output/ path**

## Q4: Log Capture (Stdout/Stderr)
Options:
- Stream to `logs/stdout.log` and `logs/stderr.log` directly
- Store in DuckDB metadata database
Selected: **Stream to `logs/stdout.log` and `logs/stderr.log` directly**
