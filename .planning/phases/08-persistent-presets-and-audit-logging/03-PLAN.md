---
wave: 3
depends_on: [02]
files_modified:
  - web/src/pages/ExecutionPage.tsx
  - web/src/components/PresetBar.tsx
  - web/src/pages/AuditPage.tsx
  - web/src/App.tsx
autonomous: true
---

# 03-PLAN: Frontend UI

<objective>
Update the frontend to pass `tool_id` to presets, support deleting presets, and implement a new Audit Log tab.
</objective>

<requirements_addressed>
- PRESET-01
- AUDIT-01
</requirements_addressed>

## Tasks

```xml
<task>
  <description>Update Preset Save and List for per-tool scoping</description>
  <read_first>
    - web/src/pages/ExecutionPage.tsx
  </read_first>
  <action>
    In `web/src/pages/ExecutionPage.tsx`:
    1. Update the `fetch('/api/presets')` in the `useEffect` to append the `tool_id`: `fetch('/api/presets?tool_id=' + id)`.
    2. In `handleSavePreset`, add `tool_id: id` to the `newPreset` JSON payload sent to `POST /api/presets`.
  </action>
  <acceptance_criteria>
    `grep "tool_id=" web/src/pages/ExecutionPage.tsx` exits 0.
    `grep "tool_id: id" web/src/pages/ExecutionPage.tsx` exits 0.
  </acceptance_criteria>
</task>
<task>
  <description>Implement Preset Deletion in PresetBar</description>
  <read_first>
    - web/src/components/PresetBar.tsx
    - web/src/pages/ExecutionPage.tsx
  </read_first>
  <action>
    In `web/src/components/PresetBar.tsx`:
    1. Add an optional `onDelete?: (id: string) => void` prop.
    2. Add an `isCustom: boolean` property to presets or differentiate tool-provided from user-provided presets (e.g., if a preset is in `savedPresets`, show a delete button). The easiest way: add `isUserDefined?: boolean` to the `Preset` model in TS, and render a small "X" button next to the preset name if true.
    In `web/src/pages/ExecutionPage.tsx`:
    1. Pass `isUserDefined: true` to the `newPreset` in `handleSavePreset`.
    2. Map over `savedPresets` to ensure they have `isUserDefined = true` when merging with `toolSpec.presets`.
    3. Implement `handleDeletePreset` which calls `DELETE /api/presets?id=...` (or `/api/presets/{id}` based on the backend router), and removes it from `savedPresets` state. Pass this handler to `<PresetBar onDelete={handleDeletePreset} />`.
  </action>
  <acceptance_criteria>
    `grep "onDelete" web/src/components/PresetBar.tsx` exits 0.
    `grep "handleDeletePreset" web/src/pages/ExecutionPage.tsx` exits 0.
  </acceptance_criteria>
</task>
<task>
  <description>Add Audit Log Page</description>
  <read_first>
    - web/src/App.tsx
  </read_first>
  <action>
    Create `web/src/pages/AuditPage.tsx`. It should:
    - Fetch from `GET /api/admin/audit` on mount.
    - Render a simple table of audit logs (`job_id`, `tool_id`, `exit_code`, `created_at`).
    In `web/src/App.tsx`:
    - Add a `<Route path="/audit" element={<AuditPage />} />`.
    In the main navigation (e.g. `web/src/pages/CatalogPage.tsx` or wherever the nav header is), add a link to `/audit` (e.g. "Audit Log").
  </action>
  <acceptance_criteria>
    `cat web/src/pages/AuditPage.tsx` succeeds.
    `grep "/audit" web/src/App.tsx` exits 0.
  </acceptance_criteria>
</task>
```

<verification>
`cd web && npm run build` exits 0.
</verification>

<must_haves>
- Preset delete action triggers a backend deletion.
- Audit page successfully fetches and renders a list.
</must_haves>
