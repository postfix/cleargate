# ClearGate Roadmap

**4 phases** | **15 requirements mapped** | All v1 requirements covered ✓

## Phases

### Phase 1: Core Execution Backend
**Goal:** Establish the foundational backend architecture for launching validated CLI jobs without shell execution.
**Requirements:** EXEC-01, EXEC-04, EXEC-05, TOOL-04
**Success Criteria:**
1. A hardcoded, simple ToolSpec can be loaded and validated.
2. An incoming execution request is validated against the ToolSpec.
3. The system constructs a safe `argv[]` array.
4. The system launches a sandboxed process using `exec.CommandContext` or Docker API.

### Phase 2: File Handling & Job Lifecycle
**Goal:** Support inputs/outputs, artifacts, and capture job results.
**Requirements:** FILE-01, FILE-02, FILE-03, EXEC-02
**Success Criteria:**
1. Users can upload an input file to an isolated job workspace.
2. The job runs and writes outputs to a designated folder.
3. System captures stdout, stderr, and exit code.
4. Users can download the resulting artifact securely.

### Phase 3: Tool Administration & LLM Pipeline
**Goal:** Implement the Eino LLM pipeline to draft ToolSpecs from raw help/docs.
**Requirements:** TOOL-01, TOOL-02, TOOL-03
**Success Criteria:**
1. The discovery engine extracts help and version information from a target binary.
2. The CloudWeGo Eino pipeline parses this output and generates a draft ToolSpec YAML.
3. Maintainers can review, edit, and submit the draft ToolSpec.
4. Administrators can approve ToolSpecs for production use.

### Phase 4: Frontend UI, Streaming & Presets
**Goal:** Generate the React SPA, stream logs, and enable presets and auditing.
**Requirements:** UI-01, UI-02, UI-03, EXEC-03, AUDIT-01, PRESET-01
**Success Criteria:**
1. The React SPA dynamically renders a form based on a ToolSpec.
2. Job execution status, stdout, and stderr are streamed to the frontend via SSE.
3. Users can save successful job parameters as presets and reload them.
4. An audit log correctly records execution details.

## Traceability Map
- EXEC-01: Phase 1
- EXEC-02: Phase 2
- EXEC-03: Phase 4
- EXEC-04: Phase 1
- EXEC-05: Phase 1
- FILE-01: Phase 2
- FILE-02: Phase 2
- FILE-03: Phase 2
- UI-01: Phase 4
- UI-02: Phase 4
- UI-03: Phase 4
- TOOL-01: Phase 3
- TOOL-02: Phase 3
- TOOL-03: Phase 3
- TOOL-04: Phase 1
- AUDIT-01: Phase 4
- PRESET-01: Phase 4
