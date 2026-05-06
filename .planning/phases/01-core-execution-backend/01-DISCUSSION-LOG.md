# Phase 1: Core Execution Backend - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

## Q1: ToolSpec Parsing
Options:
- Use Go structs with `yaml.v3` for strict typed parsing
- Load into unstructured `map[string]interface{}`
Selected: **[auto] Use Go structs with `yaml.v3` for strict typed parsing**

## Q2: Validation Strategy
Options:
- Validate inputs against tool spec immediately at the API boundary
- Defer validation to command construction phase
Selected: **[auto] Validate inputs against tool spec immediately at the API boundary**

## Q3: Sandboxing Interface
Options:
- Use `os/exec` calling the podman/docker CLI
- Use the Go Docker SDK
Selected: **[manual override] Use `os/exec` calling the rootless podman CLI (no Docker SDK)**

## Q4: Command Construction Mapping
Options:
- Isolate the `argv[]` builder into a pure, testable function
- Mix command builder tightly with execution logic
Selected: **[auto] Isolate the `argv[]` builder into a pure, testable function**
