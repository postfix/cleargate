# Phase 8: Validation Strategy

**Dimension 8: Nyquist Validation Requirement**
All features must be verifiable by the planner's `<acceptance_criteria>`.

## Core Verification Flows

### 1. Preset Persistence
1. Start backend.
2. Call `POST /api/presets` with a `tool_id` and values.
3. Verify it returns 201.
4. Restart backend.
5. Call `GET /api/presets?tool_id=...` and verify the preset is still present.

### 2. Preset CRUD
1. Call `DELETE /api/presets/{id}` for a saved preset.
2. Verify 200 or 204.
3. Call `GET /api/presets?tool_id=...` and verify it is gone.

### 3. Audit Logging
1. Execute a job via `POST /api/execute`.
2. Wait for completion.
3. Call `GET /api/admin/audit`.
4. Verify the `job_id` exists in the audit list, with the correct `exit_code`.
