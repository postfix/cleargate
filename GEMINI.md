<!-- GSD:project-start source:PROJECT.md -->
## Project

**ClearGate**

ClearGate is a secure CLI application gateway that converts approved command-line tools into generated web applications. It allows users to run approved tools in a sandboxed execution environment by configuring command-line flags through a generated SPA instead of a terminal.

**Core Value:** Expose CLI tools safely through generated web interfaces using strict, deterministic ToolSpec validation to prevent arbitrary command execution.

### Constraints

- **Language**: Go (Backend), React (Frontend) — Required by spec.
- **LLM Orchestration**: CloudWeGo Eino — Required for LLM integration.
- **Database**: DuckDB — For metadata storage.
- **Sandbox**: Docker/Podman — Required for isolated job execution.
- **Execution**: Backend must never invoke `sh -c`, commands must be constructed as `argv[]` arrays.
<!-- GSD:project-end -->

<!-- GSD:stack-start source:STACK.md -->
## Technology Stack

Technology stack not yet documented. Will populate after codebase mapping or first phase.
<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->
## Conventions

Conventions not yet established. Will populate as patterns emerge during development.
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->
## Architecture

Architecture not yet mapped. Follow existing patterns found in the codebase.
<!-- GSD:architecture-end -->

<!-- GSD:workflow-start source:GSD defaults -->
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd-quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd-debug` for investigation and bug fixing
- `/gsd-execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->



<!-- GSD:profile-start -->
## Developer Profile

> Profile not yet configured. Run `/gsd-profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
