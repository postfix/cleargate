# Phase 4: Frontend UI, Streaming & Presets - Context

**Gathered:** 2026-05-06
**Status:** Ready for planning

<domain>
## Phase Boundary

React SPA generation, dynamic form building from ToolSpecs, execution log streaming via Server-Sent Events, and preset management.
</domain>

<decisions>
## Implementation Decisions

### Framework & Styling Foundation
- **D-01:** Use React + Vite with **Vanilla CSS** for styling, adhering to the requirement for premium, high-quality bespoke aesthetics without relying on utility frameworks like Tailwind.

### Dynamic Form Generation
- **D-02:** Build a **Custom Form Builder**. React components will directly map to the `ToolSpec.Flags` and `ToolSpec.Inputs` schema to allow perfect control over the UI rendering rules (e.g., categories, conditionals).

### State Management
- **D-03:** Use **Native React** (`useState` and React Context) for state management. Avoid external stores like Redux or Zustand for this scope to minimize dependencies.

### Log Streaming (Server-Sent Events)
- **D-04:** Use the browser's native **`EventSource` API** to consume the log stream from the backend.

### the agent's Discretion
None.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Specification
- `SPEC.md` — Section 11 (Main User Flows) & Section 18 (Streaming Logs).
- `SPEC.md` — Section 13.4 and 13.9 (Flags and Presets UI metadata).

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/models/toolspec.go` — The JSON payload structure the frontend will receive.

### Integration Points
- Frontend proxy configuration in Vite to route `/api/*` to the Go backend.
- New `web/` or `frontend/` directory in the repository root.

</code_context>
