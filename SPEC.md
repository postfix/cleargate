# ClearGate — Production Specification

## 1. Product Definition

**ClearGate** is a secure CLI application gateway that converts approved command-line tools into generated web applications.

It allows users to:

* upload files for processing;
* configure command-line flags through a generated SPA instead of a terminal;
* run approved tools in a sandboxed execution environment;
* stream stdout/stderr and job status;
* download generated artifacts;
* save, share, and reuse profiles/presets;
* generate first-draft tool templates from `--help`, man pages, examples, and documentation;
* use a local or remote LLM only to enrich and normalize templates, not to bypass policy or invent executable behavior.

The central artifact of the product is the **ToolSpec**: a deterministic, reviewable, versioned schema that describes a CLI tool, its allowed flags, input/output files, command construction rules, UI metadata, validation rules, sandbox policy, and presets.

ClearGate is not a web shell. It does not expose arbitrary command execution. The frontend never sends shell commands. The backend never invokes `sh -c`. Every execution is compiled from validated structured input into an `argv[]` array using an allowlisted ToolSpec.

---

## 3. Product Goals

ClearGate exists to solve a common enterprise problem:

> Many useful tools are command-line applications, but giving users shell access is risky, hard to audit, hard to standardize, and unsuitable for non-expert users.

ClearGate turns approved CLI tools into secure internal web applications.

Primary goals:

1. **Expose CLI tools safely** through generated web interfaces.
2. **Prevent arbitrary command execution** by using structured ToolSpecs and strict validation.
3. **Support file upload and artifact download** for tools that process files.
4. **Capture stdout, stderr, exit code, metadata, and output files** for every run.
5. **Allow users and teams to save presets** as reusable execution profiles.
6. **Use LLMs to accelerate template creation**, but never as the authority for security or command execution.
7. **Support local-first and enterprise deployments**.
8. **Provide strong auditability**: who ran what, when, with which inputs, which version of the tool, which ToolSpec, and which sandbox policy.

---

## 4. Non-Goals

ClearGate is not:

* a browser-based terminal;
* a generic remote shell;
* a replacement for Kubernetes jobs or CI/CD systems;
* a no-policy command launcher;
* an LLM agent that decides what command to run on its own;
* a system that safely exposes every CLI automatically without human review;
* a sandbox escape prevention silver bullet.

Some tools are too dangerous to expose automatically and must use manual templates, restricted presets, or be rejected entirely.

---

## 5. Core Principle

The **ToolSpec is the product**.

Everything is derived from it:

```text
ToolSpec → generated UI
ToolSpec → input validation
ToolSpec → command construction
ToolSpec → sandbox policy
ToolSpec → file handling
ToolSpec → presets
ToolSpec → audit logs
ToolSpec → tests
ToolSpec → documentation
```

The LLM may help draft or annotate the ToolSpec, but it must not become the source of truth.

---

## 6. High-Level System Architecture

```text
                ┌───────────────────────────┐
                │ CLI binary / documentation │
                └─────────────┬─────────────┘
                              │
                              ▼
             ┌─────────────────────────────────┐
             │ Discovery Engine                │
             │ --help, man, docs, examples     │
             └─────────────┬───────────────────┘
                           │
                           ▼
             ┌─────────────────────────────────┐
             │ Eino LLM Template Pipeline      │
             │ normalize, classify, enrich     │
             └─────────────┬───────────────────┘
                           │
                           ▼
             ┌─────────────────────────────────┐
             │ ToolSpec Validator              │
             │ schema, flags, policy, tests    │
             └─────────────┬───────────────────┘
                           │
                           ▼
             ┌─────────────────────────────────┐
             │ Versioned ToolSpec Registry     │
             └─────────────┬───────────────────┘
                           │
          ┌────────────────┴─────────────────┐
          ▼                                  ▼
┌──────────────────────┐           ┌──────────────────────┐
│ Generated SPA         │           │ Secure Exec Backend   │
│ forms, presets, jobs  │           │ sandbox, jobs, files  │
└──────────┬───────────┘           └──────────┬───────────┘
           │                                  │
           ▼                                  ▼
┌──────────────────────┐           ┌──────────────────────┐
│ User/team presets     │           │ Artifacts/logs/audit  │
└──────────────────────┘           └──────────────────────┘
```

---

## 7. Technology Stack

### Backend

* Language: **Go**
* LLM orchestration: **CloudWeGo Eino**
* API style: HTTP JSON API
* Streaming: Server-Sent Events first; WebSocket optional later
* Metadata storage: DuckDB
* Artifact storage: local filesystem 
* Sandbox execution: Docker/Podman
* Spec format: YAML + JSON Schema validation
* Authentication: pluggable; local admin and users 

### Frontend

* React SPA
* Generated forms from ToolSpec-derived UI schema
* Monaco editor only for advanced schema review, not for raw command entry
* SSE-based job log stream
* Artifact browser
* Preset manager

### LLM Providers

Supported through Eino components:

* Eino-supported chat model providers

LLM use is limited to the template-generation and enrichment lifecycle. Runtime job execution does not require an LLM.

---

## 8. Eino Integration Model

ClearGate uses Eino for the **LLM-assisted ToolSpec generation pipeline**.

Eino is not used as a command runner.

Eino responsibilities:

1. Parse ugly help/man/doc text into structured candidate facts.
2. Classify flags by category.
3. Suggest human-readable labels and descriptions.
4. Detect probable input/output flags.
5. Suggest validation patterns.
6. Suggest mutually exclusive groups.
7. Suggest safe presets.
8. Produce a draft ToolSpec patch.
9. Generate review notes explaining evidence and uncertainty.

Eino pipeline shape:

```text
Raw Evidence
  ↓
Help Parser Component
  ↓
Documentation Retriever Component
  ↓
Flag Normalizer Chain
  ↓
Risk Classifier Chain
  ↓
Preset Suggestion Chain
  ↓
ToolSpec Draft Builder
  ↓
Deterministic Validator
  ↓
Human Review / Approval
```

Eino should be wrapped behind an internal interface so the rest of the product does not depend directly on LLM implementation details.

```go
type TemplateAssistant interface {
    DraftToolSpec(ctx context.Context, input DraftToolSpecInput) (*DraftToolSpecResult, error)
    EnrichToolSpec(ctx context.Context, input EnrichToolSpecInput) (*EnrichToolSpecResult, error)
    SuggestPresets(ctx context.Context, input PresetSuggestionInput) (*PresetSuggestionResult, error)
    ExplainFlag(ctx context.Context, input ExplainFlagInput) (*FlagExplanation, error)
}
```

---

## 9. Trust Boundary

The LLM output is **untrusted**.

LLM output must pass:

* JSON/YAML schema validation;
* known-flag verification;
* command dry-run tests where possible;
* policy validation;
* source evidence checks;
* dangerous flag detection;
* human approval for production activation;
* golden test generation and execution.

The LLM may produce:

* descriptions;
* labels;
* categories;
* draft schema;
* draft presets;
* risk explanations;
* review notes.

The LLM must not:

* authorize dangerous flags;
* bypass validation;
* directly construct shell commands;
* run jobs;
* invent unsupported flags without marking them low-confidence;
* decide final sandbox policy;
* decide final network policy.

---

## 10. User Roles

### Anonymous user

Optional for local-only mode. Usually disabled.

### Standard user

Can:

* view approved tools;
* upload files;
* run allowed tools using allowed presets or allowed custom values;
* view own job history;
* download own artifacts;
* create personal presets if permitted.

### Power user

Can:

* create team presets;
* view team job history;
* use advanced options if permitted.

### Tool maintainer

Can:

* create draft ToolSpecs;
* run discovery;
* run Eino-assisted enrichment;
* review generated templates;
* edit ToolSpecs;
* run validation tests;
* submit ToolSpecs for approval.

### Administrator

Can:

* approve ToolSpecs;
* configure sandbox policies;
* manage users and teams;
* configure LLM providers;
* manage artifact retention;
* view audit logs;
* disable tools or specific versions.

---

## 11. Main User Flows

### 11.1 Run an approved tool

1. User opens tool page.
2. SPA loads ToolSpec-derived UI schema.
3. User uploads input files if required.
4. User selects a preset or manually sets allowed options.
5. Frontend submits structured JSON job request.
6. Backend validates request against ToolSpec.
7. Backend creates isolated job workspace.
8. Backend compiles values into `argv[]`.
9. Backend launches sandboxed job.
10. User sees live stdout/stderr/status.
11. Backend captures output artifacts.
12. User downloads results.
13. Job metadata and audit log are persisted.

### 11.2 Save preset from a successful job

1. User runs a tool.
2. User clicks “Save as preset”.
3. System stores selected validated values.
4. Preset is associated with ToolSpec version.
5. Preset can be personal, team, or global depending on permissions.

### 11.3 Generate new ToolSpec

1. Maintainer registers executable path or container image.
2. Discovery engine runs safe introspection commands.
3. System collects evidence: help output, version, docs, examples.
4. Eino pipeline drafts ToolSpec.
5. Validator checks generated schema.
6. Maintainer reviews warnings and low-confidence fields.
7. Maintainer edits ToolSpec.
8. System runs dry-run/golden tests.
9. Admin approves ToolSpec.
10. Tool becomes visible to users.

---

## 12. ToolSpec Overview

ToolSpec is a versioned YAML document.

Example structure:

```yaml
apiVersion: cleargate.dev/v1
kind: ToolSpec

metadata:
  name: pandoc
  displayName: Pandoc
  description: Convert documents between formats.
  version: "3.1.12"
  owner: docs-platform-team
  tags: [documents, conversion]

runtime:
  executable: /usr/bin/pandoc
  argv0: pandoc
  workingDirectory: job
  timeoutSeconds: 300
  maxMemoryMB: 1024
  maxStdoutBytes: 10485760
  maxStderrBytes: 10485760
  network:
    enabled: false

sandbox:
  profile: default-no-network
  user: nonroot
  readonlyRootFilesystem: true
  seccompProfile: default

inputs:
  - id: input_file
    type: file
    required: true
    destination: input/
    maxSizeMB: 50
    allowedExtensions: [.md, .html, .docx, .txt]

flags:
  - id: output_format
    type: enum
    required: true
    default: html
    values: [html, pdf, docx, markdown]
    cli:
      render:
        sequence: ["-t", "{{value}}"]
    ui:
      label: Output format
      category: Output
      widget: select

  - id: standalone
    type: boolean
    default: true
    cli:
      render:
        whenTrue: ["--standalone"]
    ui:
      label: Standalone document
      category: Output
      widget: checkbox

positionals:
  - id: input_path
    source: input_file.path
    order: 100

outputs:
  - id: result
    type: file
    path: output/result.{{flags.output_format}}
    downloadable: true

presets:
  - id: html_standalone
    name: HTML standalone
    values:
      output_format: html
      standalone: true
```

---

## 13. ToolSpec Sections

### 13.1 Metadata

```yaml
metadata:
  name: string
  displayName: string
  description: string
  version: string
  owner: string
  tags: string[]
  homepage: string
  documentation: string[]
```

### 13.2 Runtime

Defines how to launch the tool.

```yaml
runtime:
  executable: /usr/bin/tool
  argv0: tool
  containerImage: optional/image:tag
  workingDirectory: job
  timeoutSeconds: 600
  maxMemoryMB: 2048
  maxCPUPercent: 200
  maxStdoutBytes: 10485760
  maxStderrBytes: 10485760
  environment:
    allowlist: []
    fixed:
      LC_ALL: C.UTF-8
  network:
    enabled: false
```

### 13.3 Inputs

Supported input types:

* file;
* directory archive;
* text;
* JSON;
* secret reference;
* scalar values;
* generated temporary file.

Example:

```yaml
inputs:
  - id: source_archive
    type: file
    required: true
    destination: input/
    maxSizeMB: 200
    allowedExtensions: [.zip, .tar.gz]
    unpack:
      enabled: true
      maxFiles: 10000
      maxExpandedSizeMB: 1000
      denySymlinks: true
```

### 13.4 Flags

Supported flag types:

* boolean;
* string;
* integer;
* float;
* enum;
* multi-enum;
* path;
* file;
* duration;
* key-value map;
* repeatable value;
* raw JSON object constrained by schema.

Example:

```yaml
flags:
  - id: severity
    type: multi-enum
    values: [LOW, MEDIUM, HIGH, CRITICAL]
    cli:
      render:
        sequence: ["--severity", "{{join(value, ',')}}"]
```

### 13.5 CLI Rendering

Rendering rules are explicit.

```yaml
cli:
  render:
    whenTrue: ["--flag"]
```

```yaml
cli:
  render:
    sequence: ["--output", "{{value}}"]
```

```yaml
cli:
  render:
    keyValue: "--define={{key}}={{value}}"
```

```yaml
cli:
  render:
    repeatEach:
      sequence: ["-I", "{{item}}"]
```

The renderer must never concatenate a full shell command.

### 13.6 Validation

```yaml
validation:
  pattern: "^[0-9,-]+$"
  min: 1
  max: 65535
  maxLength: 128
  requiredWhen:
    flag: mode
    equals: custom
```

### 13.7 Mutually Exclusive Groups

```yaml
groups:
  - id: scan_type
    type: exactly_one
    members: [tcp_connect, syn_scan, udp_scan]
```

### 13.8 Output Artifacts

```yaml
outputs:
  - id: sarif_report
    type: file
    path: output/report.sarif
    downloadable: true
    contentType: application/sarif+json
```

### 13.9 Presets

```yaml
presets:
  - id: quick_scan
    name: Quick scan
    description: Fast default scan for normal users.
    visibility: global
    locked: false
    values:
      severity: [HIGH, CRITICAL]
      format: sarif
```

### 13.10 Security Policy

```yaml
security:
  riskLevel: medium
  requireApprovalForRun: false
  denyUnknownFlags: true
  allowRawArgs: false
  pathPolicy:
    allowAbsolutePaths: false
    allowParentTraversal: false
    allowSymlinks: false
  deniedFlags:
    - --exec
    - --script
  network:
    enabled: false
```

---

## 14. Command Builder

The command builder compiles validated job values into `argv[]`.

Input:

```json
{
  "tool": "pandoc",
  "toolSpecVersion": "sha256:...",
  "values": {
    "input_file": "upload_123",
    "output_format": "html",
    "standalone": true
  }
}
```

Output:

```json
{
  "argv": [
    "/usr/bin/pandoc",
    "-t",
    "html",
    "--standalone",
    "/workspace/input/input.md",
    "-o",
    "/workspace/output/result.html"
  ]
}
```

Rules:

1. Unknown fields are rejected.
2. Unknown flags are rejected.
3. Disabled flags are rejected.
4. Values are validated before rendering.
5. Paths are converted to workspace-relative safe paths.
6. Output paths are synthesized by policy, not user-controlled unless explicitly allowed.
7. The command is launched using `exec.CommandContext` or container runtime equivalent with argv array.
8. Shell execution is forbidden.

---

## 15. Execution Sandbox

### MVP sandbox

* Podman container per job.
* Non-root user.
* Read-only root filesystem where possible.
* Mounted job workspace.
* No host path mounts except controlled input/output workspace.
* CPU/memory/time limits.
* Network disabled by default.

### Hardened sandbox

Future production profiles:
* rootless Podman.

### Sandbox profiles

```yaml
sandboxProfiles:
  - id: default-no-network
    network: false
    rootfs: readonly
    user: nonroot
    seccomp: default
    maxProcesses: 64

  - id: network-restricted
    network: true
    allowedCIDRs:
      - 10.0.0.0/8
      - 192.168.0.0/16
```

---

## 16. File Handling

### Upload rules

* Every upload receives an internal ID.
* Original filenames are stored as metadata only.
* Files are copied into job workspace using safe generated names.
* File extension and MIME checks are applied.
* Size limits are enforced.
* Archive expansion is controlled.
* Symlinks are denied by default.
* Path traversal is denied.

### Workspace layout

```text
/jobs/{job_id}/
  input/
  output/
  logs/
    stdout.log
    stderr.log
  metadata/
    job.json
    argv.json
```

### Download rules

* Only declared output artifacts are downloadable by default.
* Optional artifact discovery can be enabled for selected tools.
* Symlink escapes are denied.
* Absolute paths are denied.
* Files outside workspace are never downloadable.

---

## 17. Job Lifecycle

Job states:

```text
created
validating
queued
preparing_workspace
running
collecting_artifacts
succeeded
failed
cancelled
timeout
policy_denied
sandbox_error
```

Job record:

```json
{
  "id": "job_01H...",
  "tool": "pandoc",
  "toolSpecVersion": "sha256:...",
  "presetId": "html_standalone",
  "status": "succeeded",
  "createdBy": "user_123",
  "createdAt": "2026-05-06T10:00:00Z",
  "startedAt": "2026-05-06T10:00:02Z",
  "finishedAt": "2026-05-06T10:00:08Z",
  "exitCode": 0,
  "limits": {
    "timeoutSeconds": 300,
    "maxMemoryMB": 1024
  }
}
```

---

## 18. Streaming Logs

MVP uses Server-Sent Events.

Endpoint:

```http
GET /api/jobs/{job_id}/events
```

Events:

```json
{"type":"status","status":"running"}
{"type":"stdout","data":"..."}
{"type":"stderr","data":"..."}
{"type":"artifact","artifactId":"result"}
{"type":"complete","status":"succeeded","exitCode":0}
```

Rules:

* stdout/stderr are size-limited;
* logs are stored as artifacts;
* binary output to stdout can be captured but should be handled carefully;
* logs should be redacted if secrets are involved.

---

## 19. Presets and Profiles

Preset types:

* system preset;
* team preset;
* user preset;
* locked compliance preset;
* temporary run configuration.

Preset schema:

```yaml
apiVersion: cleargate.dev/v1
kind: Preset
metadata:
  id: web_mp4
  name: Web MP4
  owner: media-team
  visibility: team
spec:
  tool: ffmpeg
  toolSpecVersionConstraint: ">=1.0.0 <2.0.0"
  values:
    video_codec: libx264
    crf: 23
    output_format: mp4
```

Preset validation:

* validated at creation;
* revalidated after ToolSpec upgrade;
* invalid presets are marked stale;
* stale presets cannot run until migrated or approved.

---

## 20. Discovery Engine

Discovery engine collects evidence.

Supported evidence sources:

* `tool --help`;
* `tool -h`;
* `tool help`;
* `tool subcommand --help`;
* `tool --version`;
* `man tool`;
* README files;
* docs directory;
* examples;
* known command snippets;
* existing shell scripts;
* maintainer-provided notes.

Discovery output:

```json
{
  "tool": "pandoc",
  "version": "3.1.12",
  "helpOutputs": [...],
  "manpageText": "...",
  "examples": [...],
  "detectedFlags": [...],
  "detectedSubcommands": [...],
  "rawEvidence": [...]
}
```

Discovery must run in a safe environment. It should not execute arbitrary examples from documentation.

---

## 21. LLM-Assisted Template Generation

### 21.1 Input to Eino pipeline

```json
{
  "tool": "pandoc",
  "version": "3.1.12",
  "help": "...",
  "manpage": "...",
  "examples": ["pandoc input.md -o output.html"],
  "schemaVersion": "cleargate.dev/v1",
  "policyHints": {
    "networkDefault": false,
    "allowRawArgs": false
  }
}
```

### 21.2 Output from Eino pipeline

```json
{
  "draftToolSpec": {...},
  "confidence": "medium",
  "warnings": [
    "Flag --lua-filter executes user-provided Lua code; mark as dangerous."
  ],
  "evidenceMap": [
    {
      "field": "flags.output_format",
      "source": "help",
      "confidence": "high"
    }
  ]
}
```

### 21.3 Validation after LLM

The generated draft must be checked by deterministic validators:

* schema validity;
* all flags exist in evidence or are explicitly maintainer-added;
* dangerous flags are marked;
* unsupported fields rejected;
* render rules are valid;
* generated example jobs compile to argv;
* dry-run tests pass where possible;
* security policy is present.

---

## 22. Confidence Model

Each discovered item has confidence.

```yaml
confidence:
  sourceHelp: true
  sourceManpage: true
  sourceExamples: false
  llmInferred: false
  level: high
```

Levels:

* high: appears in help/man and examples or has deterministic parse evidence;
* medium: appears in one reliable source;
* low: LLM inferred or ambiguous;
* rejected: failed validation or marked unsafe.

Low-confidence executable fields require human review before activation.

---

## 23. Risk Classification

Flags and tools receive risk classification.

### Flag risk levels

* safe;
* low;
* medium;
* high;
* forbidden.

### Dangerous categories

* executes scripts or commands;
* loads external configuration;
* enables network access;
* writes arbitrary filesystem paths;
* deletes files;
* accepts raw expressions;
* enables plugins;
* disables security checks;
* accesses credentials;
* performs scanning or brute force;
* consumes excessive resources.

Example:

```yaml
flags:
  - id: lua_filter
    cli:
      render:
        sequence: ["--lua-filter", "{{value}}"]
    risk:
      level: high
      reason: Executes user-provided Lua filter.
      defaultEnabled: false
      requiresRole: tool-maintainer
```

---

## 24. API Specification

### 24.1 Tools

```http
GET /api/tools
GET /api/tools/{tool}
GET /api/tools/{tool}/schema
GET /api/tools/{tool}/presets
```

### 24.2 Jobs

```http
POST /api/tools/{tool}/jobs
GET /api/jobs/{job_id}
GET /api/jobs/{job_id}/events
POST /api/jobs/{job_id}/cancel
GET /api/jobs/{job_id}/artifacts
GET /api/jobs/{job_id}/artifacts/{artifact_id}/download
```

### 24.3 Uploads

```http
POST /api/uploads
GET /api/uploads/{upload_id}
DELETE /api/uploads/{upload_id}
```

### 24.4 Presets

```http
POST /api/tools/{tool}/presets
PUT /api/tools/{tool}/presets/{preset_id}
DELETE /api/tools/{tool}/presets/{preset_id}
```

### 24.5 ToolSpec Management

```http
POST /api/admin/tools/discover
POST /api/admin/tools/draft
POST /api/admin/tools/validate
POST /api/admin/tools/approve
GET /api/admin/tools/{tool}/versions
```

---

## 25. Backend Package Layout

```text
cmd/cleargate-server/
  main.go

internal/api/
  handlers.go
  middleware.go
  routes.go

internal/spec/
  model.go
  loader.go
  validator.go
  renderer.go
  migration.go

internal/discovery/
  runner.go
  parser.go
  evidence.go

internal/assistant/
  eino_client.go
  draft_pipeline.go
  prompts.go
  schemas.go

internal/jobs/
  service.go
  queue.go
  lifecycle.go
  events.go

internal/executor/
  command_builder.go
  runner.go
  limits.go

internal/sandbox/
  interface.go
  docker.go
  podman.go
  local_dev.go

internal/artifacts/
  store.go
  local.go
  s3.go

internal/uploads/
  service.go
  validation.go

internal/presets/
  service.go
  migration.go

internal/audit/
  audit.go
  sink.go

internal/auth/
  identity.go
  rbac.go

internal/storage/
  db.go
  migrations/

web/
  src/
```

---

## 26. Core Go Interfaces

### 26.1 Spec Repository

```go
type ToolSpecRepository interface {
    GetActive(ctx context.Context, tool string) (*ToolSpec, error)
    GetVersion(ctx context.Context, tool string, version string) (*ToolSpec, error)
    SaveDraft(ctx context.Context, spec *ToolSpec) error
    Approve(ctx context.Context, tool string, version string, actor Actor) error
    ListTools(ctx context.Context, actor Actor) ([]ToolSummary, error)
}
```

### 26.2 Command Builder

```go
type CommandBuilder interface {
    Build(ctx context.Context, spec *ToolSpec, req JobRequest, ws Workspace) (*CommandPlan, error)
}

type CommandPlan struct {
    Executable string
    Args       []string
    Env        map[string]string
    WorkDir    string
    Outputs    []ExpectedArtifact
}
```

### 26.3 Sandbox Runner

```go
type SandboxRunner interface {
    Run(ctx context.Context, plan CommandPlan, policy SandboxPolicy, sink EventSink) (*RunResult, error)
    Cancel(ctx context.Context, jobID string) error
}
```

### 26.4 Artifact Store

```go
type ArtifactStore interface {
    Put(ctx context.Context, artifact Artifact) error
    Get(ctx context.Context, id string) (*Artifact, error)
    OpenRead(ctx context.Context, id string) (io.ReadCloser, error)
    DeleteExpired(ctx context.Context) error
}
```

### 26.5 Template Assistant

```go
type TemplateAssistant interface {
    DraftToolSpec(ctx context.Context, input DraftToolSpecInput) (*DraftToolSpecResult, error)
    EnrichToolSpec(ctx context.Context, input EnrichToolSpecInput) (*EnrichToolSpecResult, error)
    SuggestPresets(ctx context.Context, input PresetSuggestionInput) (*PresetSuggestionResult, error)
}
```

---

## 27. Frontend Structure

Pages:

* `/tools` — tool catalog;
* `/tools/:tool` — generated tool runner UI;
* `/tools/:tool/presets` — preset manager;
* `/jobs/:jobId` — job detail and live logs;
* `/admin/tools` — ToolSpec management;
* `/admin/tools/:tool/spec` — spec editor/reviewer;
* `/admin/audit` — audit logs;
* `/settings/llm` — LLM provider configuration.

Generated form behavior:

* reads UI schema generated from ToolSpec;
* groups fields by category;
* hides advanced and dangerous fields by default;
* validates client-side for convenience;
* backend remains the source of truth;
* allows “save as preset” after successful validation.

---

## 28. Audit Logging

Every important action must be auditable.

Audit events:

* tool discovered;
* ToolSpec draft generated;
* ToolSpec edited;
* ToolSpec approved;
* job created;
* job started;
* job completed;
* job cancelled;
* artifact downloaded;
* preset created;
* preset changed;
* policy denied execution;
* admin changed sandbox policy;
* LLM provider changed.

Audit event shape:

```json
{
  "id": "audit_01H...",
  "timestamp": "2026-05-06T10:00:00Z",
  "actor": "user_123",
  "action": "job.created",
  "resource": "job_456",
  "tool": "pandoc",
  "toolSpecVersion": "sha256:...",
  "ip": "10.1.2.3",
  "metadata": {}
}
```

---

## 29. Security Requirements

### 29.1 Mandatory

* No shell invocation for jobs.
* No raw arbitrary arguments unless explicitly enabled for a tool and role.
* All executable paths are allowlisted.
* All flags are allowlisted.
* All values are validated.
* All paths are workspace-confined.
* Network is disabled by default.
* Jobs run as non-root.
* Timeouts are mandatory.
* Memory limits are mandatory.
* Output size limits are mandatory.
* Artifact downloads are workspace-confined.
* ToolSpec changes require approval before production activation.

### 29.2 Recommended

* Rootless containers.
* seccomp profile.
* AppArmor/SELinux profile.
* Read-only root filesystem.
* Per-job ephemeral workspace.
* Per-job container.
* Separate artifact storage from execution nodes.
* Malware scanning for uploads in enterprise mode.
* OIDC/SAML integration.
* Signed ToolSpecs.

### 29.3 Forbidden by default

* `sh -c`;
* arbitrary interpreter execution;
* arbitrary host path mounts;
* Docker socket mount;
* Kubernetes service account mount;
* unrestricted outbound network;
* arbitrary output paths;
* symlink-following downloads;
* LLM-generated execution without validation.

---

## 30. Tool Risk Categories

### Good MVP candidates

* pandoc;
* jq;
* yq;
* trivy filesystem/image scans with controlled input;
* syft;
* grype;
* semgrep with controlled rules;
* ffmpeg with network disabled.

### Use with caution

* nmap;
* openssl;
* imagemagick;
* sqlite3;
* git;
* curl;
* wget.

### Manual-template-only or deny by default

* bash;
* sh;
* python;
* node;
* perl;
* ruby;
* docker;
* kubectl;
* aws;
* gcloud;
* az;
* ssh;
* ansible;
* terraform.

---

## 31. MVP Scope

MVP must support:

* one backend server;
* React SPA;
* local admin account;
* ToolSpec YAML loading;
* generated forms;
* file upload;
* file output download;
* boolean/string/integer/enum flags;
* command rendering to argv;
* Docker or Podman sandbox;
* no-network default sandbox;
* job lifecycle;
* SSE log streaming;
* local filesystem artifacts;
* SQLite metadata;
* personal presets;
* manual ToolSpec creation;
* Eino-assisted draft generation for simple tools;
* validation pipeline;
* audit log.

MVP should target these tools:

1. pandoc
2. jq
3. trivy
4. syft
5. grype
6. semgrep
7. ffmpeg, limited safe profile

---

## 32. Phase Plan

### Phase 0 — Foundation

* Define ToolSpec v1.
* Implement schema validation.
* Implement command renderer.
* Implement path safety library.
* Implement local job workspace.
* Implement basic audit log.

### Phase 1 — Manual Tool Runner

* Load manual ToolSpecs.
* Generate simple UI forms.
* Submit jobs.
* Run sandboxed commands.
* Stream logs.
* Download declared artifacts.

### Phase 2 — Presets

* Save personal presets.
* Load presets into UI.
* Re-run jobs from presets.
* Validate presets against ToolSpec version.

### Phase 3 — Discovery

* Run `--help`, `-h`, `--version`.
* Capture raw evidence.
* Parse obvious flags deterministically.
* Store evidence map.

### Phase 4 — Eino Template Assistant

* Configure Ollama/OpenAI-compatible provider.
* Build Eino chain for draft ToolSpec generation.
* Validate generated ToolSpec.
* Show review UI.
* Require human approval.

### Phase 5 — Enterprise Hardening

* OIDC/SAML.
* RBAC.
* Team presets.
* PostgreSQL.
* S3 artifact storage.
* Advanced sandbox profiles.
* Signed ToolSpecs.
* Audit export.

---

## 33. Acceptance Criteria

### Command safety

* Backend never invokes shell for job execution.
* Unknown flags are rejected.
* Unknown values are rejected.
* Path traversal attempts are rejected.
* Absolute user-provided output paths are rejected.

### Sandbox safety

* Job runs as non-root.
* Job workspace is isolated.
* Network is disabled by default.
* Timeout kills the job.
* Memory limit is enforced.

### UI generation

* UI is generated from ToolSpec.
* Required fields are shown.
* Advanced fields can be collapsed.
* Dangerous fields are hidden or permission-gated.

### Presets

* Preset can be saved after valid run.
* Preset can be reloaded.
* Preset is rejected if incompatible with active ToolSpec.

### LLM pipeline

* Eino can draft ToolSpec from help output.
* LLM-generated ToolSpec is not activated without validation.
* Low-confidence flags are marked for review.
* Dangerous flags are surfaced as warnings.

---

## 34. Testing Strategy

### Unit tests

* ToolSpec validation;
* CLI renderer;
* path normalization;
* flag validation;
* preset validation;
* risk classifier rules;
* workspace artifact collection.

### Integration tests

* run pandoc fixture;
* run jq fixture;
* upload file, process, download artifact;
* cancel long-running job;
* timeout job;
* reject dangerous path;
* reject unknown flag;
* reject stale preset.

### Golden tests

Golden ToolSpecs:

* pandoc;
* jq;
* trivy;
* syft;
* semgrep;
* ffmpeg safe subset.

Golden outputs:

* expected UI schema;
* expected argv rendering;
* expected validation errors;
* expected artifact list.

### Security tests

* command injection attempts;
* shell metacharacters;
* path traversal;
* symlink escape;
* archive bomb;
* oversized stdout;
* forbidden network test;
* forbidden flag use;
* raw argument rejection.

---

## 35. Example ToolSpec: jq

```yaml
apiVersion: cleargate.dev/v1
kind: ToolSpec
metadata:
  name: jq
  displayName: jq
  description: Process JSON files with jq filters.
  version: "1.7"

runtime:
  executable: /usr/bin/jq
  workingDirectory: job
  timeoutSeconds: 60
  maxMemoryMB: 256
  network:
    enabled: false

inputs:
  - id: input_json
    type: file
    required: true
    destination: input/
    maxSizeMB: 20
    allowedExtensions: [.json]

flags:
  - id: filter
    type: string
    required: true
    validation:
      maxLength: 2048
    cli:
      render:
        positional: true
        order: 10
    ui:
      label: jq filter
      widget: textarea
      category: Query

  - id: raw_output
    type: boolean
    default: false
    cli:
      render:
        whenTrue: ["-r"]
    ui:
      label: Raw output
      widget: checkbox
      category: Output

positionals:
  - id: input_path
    source: input_json.path
    order: 20

outputs:
  - id: stdout
    type: stdout
    downloadable: true
```

Note: `jq` filters are a controlled expression language but still require resource limits and careful output limits.

---

## 36. Example ToolSpec: pandoc

```yaml
apiVersion: cleargate.dev/v1
kind: ToolSpec
metadata:
  name: pandoc
  displayName: Pandoc
  description: Convert documents between formats.

runtime:
  executable: /usr/bin/pandoc
  workingDirectory: job
  timeoutSeconds: 300
  maxMemoryMB: 1024
  network:
    enabled: false

inputs:
  - id: input_file
    type: file
    required: true
    destination: input/
    maxSizeMB: 50
    allowedExtensions: [.md, .html, .docx, .txt]

flags:
  - id: output_format
    type: enum
    required: true
    default: html
    values: [html, docx, markdown]
    ui:
      label: Output format
      widget: select
      category: Output

  - id: standalone
    type: boolean
    default: true
    cli:
      render:
        whenTrue: ["--standalone"]
    ui:
      label: Standalone document
      widget: checkbox
      category: Output

positionals:
  - id: input_path
    source: input_file.path
    order: 10

outputs:
  - id: result
    type: file
    path: output/result.{{flags.output_format}}
    cli:
      render:
        sequence: ["-o", "{{path}}"]
    downloadable: true
```

---

## 37. Open Design Questions

1. Should ClearGate support multi-step workflows in v1, or only single-command jobs?
2. Should generated UI schema be a separate artifact or derived live from ToolSpec?
3. Should dangerous flags be hidden, disabled, or shown with role-based unlock?
4. Should presets support approval workflows?
5. Should tool execution happen on the same node as API server in MVP, or through worker nodes from day one?
6. Should ToolSpecs be signed with a local key in MVP or only in enterprise phase?
7. How much auto-discovery is acceptable before human review?

Recommended answers for MVP:

1. single-command jobs only;
2. derive UI schema from ToolSpec;
3. hide forbidden flags, permission-gate dangerous flags;
4. personal presets no approval, team/global presets require approval;
5. same node for MVP, worker nodes later;
6. no signing in MVP, but store content hash;
7. auto-discovery allowed only for draft specs.

---

## 38. Final Product Statement

**ClearGate is a secure CLI application gateway for turning approved command-line tools into generated, auditable, policy-controlled web applications. It uses Go for the backend, Eino for LLM-assisted ToolSpec generation, sandboxed execution for safe runtime isolation, and a generated SPA for user-friendly execution, file processing, result download, and reusable presets.**
