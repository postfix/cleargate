# Phase 4: Frontend UI, Streaming & Presets - Research

## Objective
Establish the implementation strategy for the React SPA, dynamic form rendering from `ToolSpec`, server-sent events (SSE) log streaming, and preset management.

## Tech Stack & Tooling
- **Framework**: React 18+ with Vite (`npm create vite@latest frontend -- --template react-ts`).
- **Styling**: Vanilla CSS with CSS Variables mapped directly to the design tokens in `04-UI-SPEC.md`.
- **Icons**: `lucide-react` for standard UI icons.
- **Routing**: `react-router-dom` for navigating between the Tool Catalog and the execution page.
- **Proxy**: Configure `vite.config.ts` to proxy API requests to the Go backend (`http://localhost:8080/api/`).

## Dynamic Form Generation
- **Pattern**: Custom React Component tree.
- **Mapping Strategy**:
  - The backend provides a `ToolSpec` object via `/api/tools/{id}`.
  - The `flags` array dictates the form fields.
  - Group fields visually based on `ui.category`.
  - Field mapping:
    - `type: boolean` → Custom Toggle component.
    - `type: enum` → Custom Select component.
    - `type: string` → Custom TextInput component.
- **State Handling**: Maintain a flat form state object (e.g., `Record<string, string | boolean>`) that updates `onChange`.

## Server-Sent Events (SSE) Stream
- **Pattern**: `EventSource` browser API.
- **Implementation**:
  - Endpoint: `GET /api/jobs/{job_id}/events`
  - React `useEffect` to manage the `EventSource` lifecycle.
  - State: Append incoming `LogEvent` objects to an array for the terminal UI.
  - Auto-scroll: Use a ref on the log container and `scrollIntoView()` on new messages unless the user toggles "scroll lock".
  - Status updates: Listen for `{type: "status"}` and `{type: "complete"}` to update the top sticky status bar.

## Preset Management
- **Pattern**: Local Storage (MVP) or Backend DB. The requirements specify saving successful job parameters. For MVP, we will hit a backend endpoint to save the preset.
- **Backend API (Mock/Implementation)**:
  - `POST /api/presets`
  - `GET /api/presets`
- **UI Interaction**: A pill bar above the form allows selecting a preset, which instantly overwrites the form state object.

## Testing & Validation Strategy (Nyquist)
- Use standard Vite `vitest` for unit tests if needed, but primary validation will be end-to-end functionality of form rendering and SSE consumption.
