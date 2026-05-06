---
wave: 1
depends_on: []
files_modified:
  - web/package.json
  - web/index.css
  - web/src/App.tsx
  - web/src/components/ToolCatalog.tsx
  - web/src/components/DynamicForm.tsx
  - web/src/components/PresetBar.tsx
  - web/src/components/LogStream.tsx
autonomous: true
---

# Phase 4: Frontend UI, Streaming & Presets

<objective>
Generate the React SPA with a Tool Catalog, Dynamic Form Builder derived from ToolSpecs, execution log streaming via SSE, and preset management functionality.
</objective>

<requirements>
- UI-01: System dynamically generates forms and SPA UI from ToolSpec schemas.
- UI-02: User can view a catalog of approved tools.
- UI-03: User can select and apply presets to fill out tool forms.
- EXEC-03: System streams job status, stdout, and stderr live to the UI.
- PRESET-01: Users can save and load job runs as reusable presets.
- AUDIT-01: (Handled implicitly by backend log capture, frontend just streams).
</requirements>

<tasks>

<task>
<id>1</id>
<title>Vite SPA Scaffold & Design System</title>
<read_first>
- .planning/phases/04-frontend-ui-streaming-presets/04-UI-SPEC.md
</read_first>
<action>
1. Run `npm create vite@latest web -- --template react-ts`.
2. Configure `web/vite.config.ts` to proxy `/api` requests to `http://localhost:8080`.
3. In `web/index.css`, establish CSS Variables derived precisely from `04-UI-SPEC.md` (e.g., `--color-dominant: #0F1117;`, `--color-accent: #6C63FF;`, typography scale, spacing).
4. Install dependencies: `lucide-react`, `react-router-dom`.
</action>
<acceptance_criteria>
- `web/package.json` exists.
- `web/index.css` contains the design tokens specified in UI-SPEC.
</acceptance_criteria>
</task>

<task>
<id>2</id>
<title>Tool Catalog Page</title>
<depends_on>1</depends_on>
<read_first>
- web/index.css
</read_first>
<action>
1. Create `web/src/components/ToolCatalog.tsx`.
2. Fetch tools from `/api/admin/tools/drafts` (or a new approved tools endpoint, for MVP drafts is fine if approved endpoint isn't wired fully for listing yet).
3. Render a grid of cards using the CSS classes derived from the UI-SPEC (3-column desktop, border `#2E3039`).
4. Set up `react-router-dom` in `App.tsx` to route `/` to Catalog and `/tool/:id` to execution page.
</action>
<acceptance_criteria>
- `web/src/components/ToolCatalog.tsx` implements the 3-column grid specified in UI-SPEC.
- Cards have `#2E3039` border.
</acceptance_criteria>
</task>

<task>
<id>3</id>
<title>Dynamic Form Builder</title>
<depends_on>2</depends_on>
<read_first>
- internal/models/toolspec.go
</read_first>
<action>
1. Create `web/src/components/DynamicForm.tsx`.
2. Component receives a `ToolSpec` JSON object.
3. Iterate over `spec.flags` and `spec.inputs`.
4. Group fields visually by `ui.category` (if present).
5. Render specific inputs based on `type`: `boolean` (toggle), `enum` (select), `string` (text), `file` (file input).
6. Manage form state using React `useState` (`Record<string, any>`).
7. Implement a "Run Tool" button (styled with accent color `#6C63FF`).
</action>
<acceptance_criteria>
- `DynamicForm.tsx` correctly renders inputs based on type.
- State is successfully aggregated on submission.
</acceptance_criteria>
</task>

<task>
<id>4</id>
<title>SSE Log Stream Component</title>
<depends_on>3</depends_on>
<read_first>
- web/index.css
</read_first>
<action>
1. Create `web/src/components/LogStream.tsx`.
2. When a job starts, connect to `GET /api/jobs/{job_id}/events` using the native `EventSource` API.
3. Render a terminal-style panel (`#0A0C10` background, monospace).
4. Parse incoming events (`{type: "stdout", data: "..."}`). Color code stderr red (`#E5484D`).
5. Auto-scroll to bottom as new messages arrive.
</action>
<acceptance_criteria>
- `LogStream.tsx` uses `new EventSource()`.
- Terminal panel has `#0A0C10` background.
</acceptance_criteria>
</task>

<task>
<id>5</id>
<title>Preset Management UI</title>
<depends_on>3</depends_on>
<read_first>
- web/index.css
</read_first>
<action>
1. Create `web/src/components/PresetBar.tsx`.
2. Render a row of pill buttons. "Custom" is always first.
3. Clicking a preset pill invokes an `onSelect` callback that populates the `DynamicForm` state with `preset.values`.
4. Implement a "Save Preset" button that serializes current form state.
</action>
<acceptance_criteria>
- Clicking a preset updates the form state.
- Preset pills follow UI-SPEC styling.
</acceptance_criteria>
</task>

</tasks>

<verification>
- Vite server starts successfully.
- ToolSpec JSON translates correctly to interactive form elements.
- SSE connection properly receives and renders log chunks.
</verification>

<must_haves>
- Uses Vanilla CSS mapping to `04-UI-SPEC.md` colors and spacing.
- Forms are strictly derived from the ToolSpec payload.
- Log panel uses `EventSource` and is styled like a terminal.
</must_haves>
