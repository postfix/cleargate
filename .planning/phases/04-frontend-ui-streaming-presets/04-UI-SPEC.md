---
phase: 04
slug: frontend-ui-streaming-presets
status: draft
shadcn_initialized: false
preset: none
created: 2026-05-06
---

# Phase 4 — UI Design Contract

> Visual and interaction contract for the ClearGate frontend. Covers the React SPA that renders dynamic tool forms, streams execution logs, and manages presets.

---

## Design System

| Property | Value |
|----------|-------|
| Tool | none (Vanilla CSS) |
| Preset | not applicable |
| Component library | none (custom React components) |
| Icon library | Lucide React |
| Font | Inter (Google Fonts) |

---

## Spacing Scale

Declared values (must be multiples of 4):

| Token | Value | Usage |
|-------|-------|-------|
| xs | 4px | Icon gaps, inline padding |
| sm | 8px | Compact element spacing, form field internal padding |
| md | 16px | Default element spacing, card padding |
| lg | 24px | Section padding, form group gaps |
| xl | 32px | Layout gaps between major sections |
| 2xl | 48px | Page header to content gap |
| 3xl | 64px | Page-level vertical margins |

Exceptions: none

---

## Typography

| Role | Size | Weight | Line Height |
|------|------|--------|-------------|
| Body | 14px | 400 | 1.6 |
| Label | 13px | 500 | 1.4 |
| Heading | 20px | 600 | 1.3 |
| Display | 28px | 700 | 1.2 |
| Code/Log | 13px (monospace: JetBrains Mono) | 400 | 1.5 |

---

## Color

| Role | Value | Usage |
|------|-------|-------|
| Dominant (60%) | #0F1117 | Page background, main surfaces |
| Secondary (30%) | #1A1D27 | Cards, sidebar, input backgrounds |
| Accent (10%) | #6C63FF | Primary CTA buttons, active tab indicators, selected preset border |
| Destructive | #E5484D | Cancel job button, error badges |
| Success | #30A46C | Job succeeded badge, connected status dot |
| Warning | #F5A623 | Job timeout badge, draft status pill |
| Text Primary | #EDEEF0 | Headings, body text |
| Text Secondary | #8B8D98 | Labels, descriptions, timestamps |
| Border | #2E3039 | Card borders, input borders, dividers |
| Surface Elevated | #22252F | Dropdown menus, modals, tooltip backgrounds |

Accent reserved for: primary action buttons ("Run Tool", "Save Preset"), active navigation tab underline, selected preset card left-border highlight, SSE connection active indicator.

---

## Component Specifications

### Tool Catalog Page
- Grid layout: 3 columns on desktop (min 280px per card), 1 column on mobile
- Tool card: 1A1D27 background, 2E3039 border, 8px border-radius
- Card contents: tool icon (32x32), display name (Heading), description truncated to 2 lines (Body), tags as pills (xs padding, 6C63FF/15% background)
- Hover: border transitions to 6C63FF at 40% opacity, 200ms ease

### Dynamic Form (Tool Execution Page)
- Form fields grouped by `ui.category` from ToolSpec
- Category headers: Label weight, Text Secondary color, uppercase, 1px bottom border
- Input types mapped from ToolSpec flag types:
  - `boolean` → toggle switch (44px wide, 24px tall, accent when on)
  - `enum` → select dropdown (custom styled, Surface Elevated dropdown)
  - `string` → text input (Secondary background, Border border, 40px height)
  - `file` → drag-and-drop zone (dashed Border border, 120px min-height)
- Required fields: accent-colored asterisk after label
- Validation errors: Destructive color text below input, 1px Destructive left border on input

### Preset Selector
- Horizontal pill bar above the form
- Each preset: pill shape, Secondary background, Border border
- Selected preset: accent left border (3px), accent/15% background
- "Custom" option always present as first pill

### Log Stream Panel
- Full-width panel below form, collapsible (chevron toggle)
- Dark terminal aesthetic: #0A0C10 background, monospace font
- stdout lines: Text Primary color
- stderr lines: Destructive color with "ERR" prefix badge
- Status events: accent color with "STATUS" prefix badge
- Auto-scroll with "scroll lock" toggle button (pin icon, top-right)
- Max visible height: 400px with overflow scroll

### Job Status Bar
- Sticky bar between form and log panel
- Status pill: color-coded (accent=running, Success=succeeded, Destructive=failed, Warning=timeout)
- Exit code display: monospace, right-aligned
- Elapsed time: live counter while running, final value when complete

---

## Copywriting Contract

| Element | Copy |
|---------|------|
| Primary CTA | "Run Tool" |
| Empty state heading (no tools) | "No tools available" |
| Empty state body (no tools) | "No approved tools found. Contact your administrator to add tools." |
| Empty state heading (no jobs) | "No recent jobs" |
| Empty state body (no jobs) | "Run a tool to see execution history here." |
| Error state (job failed) | "Job failed with exit code {N}. Check stderr output below for details." |
| Error state (connection lost) | "Connection lost. Reconnecting..." |
| Destructive confirmation | "Cancel Job: This will terminate the running process. Continue?" |
| Save preset CTA | "Save as Preset" |
| Preset saved toast | "Preset saved successfully" |
| SSE connected | "Live" (green dot + text) |
| SSE disconnected | "Disconnected" (grey dot + text) |

---

## Layout Structure

```
┌─────────────────────────────────────────────────────┐
│  Header: Logo + "ClearGate" + nav tabs              │
├─────────────────────────────────────────────────────┤
│                                                     │
│  [Tool Catalog]  or  [Tool Execution Page]          │
│                                                     │
│  Tool Execution Page:                               │
│  ┌───────────────────────────────────────────────┐  │
│  │ Tool Header: icon + name + version            │  │
│  ├───────────────────────────────────────────────┤  │
│  │ Preset Bar: [Custom] [Preset1] [Preset2]      │  │
│  ├───────────────────────────────────────────────┤  │
│  │ Dynamic Form (grouped by category)            │  │
│  │   Category: Output                            │  │
│  │     [Output Format ▼]  [Standalone ◉]         │  │
│  │   Category: Input                             │  │
│  │     [Drop files here]                         │  │
│  ├───────────────────────────────────────────────┤  │
│  │ [Run Tool]                        [Save Preset]│ │
│  ├───────────────────────────────────────────────┤  │
│  │ Status Bar: ● Running  |  00:04  |  exit: —   │  │
│  ├───────────────────────────────────────────────┤  │
│  │ Log Stream (collapsible)                      │  │
│  │ > stdout line 1                               │  │
│  │ > stdout line 2                               │  │
│  │ ERR stderr line 1                             │  │
│  └───────────────────────────────────────────────┘  │
│                                                     │
└─────────────────────────────────────────────────────┘
```

---

## Responsive Breakpoints

| Breakpoint | Width | Layout Changes |
|------------|-------|----------------|
| Desktop | ≥1024px | 3-column tool grid, full sidebar |
| Tablet | 768–1023px | 2-column tool grid, collapsed sidebar |
| Mobile | <768px | 1-column tool grid, hamburger nav |

---

## Animations & Transitions

| Element | Property | Duration | Easing |
|---------|----------|----------|--------|
| Card hover | border-color | 200ms | ease |
| Toggle switch | transform | 150ms | ease-in-out |
| Log panel collapse | max-height | 300ms | ease |
| Status pill change | background-color | 200ms | ease |
| Page transitions | opacity | 200ms | ease |
| Toast notification | transform + opacity | 300ms | cubic-bezier(0.16, 1, 0.3, 1) |

---

## Registry Safety

| Registry | Blocks Used | Safety Gate |
|----------|-------------|-------------|
| Lucide React (official) | Icons only | not required |
| No third-party registries | — | — |

---

## Checker Sign-Off

- [ ] Dimension 1 Copywriting: PASS
- [ ] Dimension 2 Visuals: PASS
- [ ] Dimension 3 Color: PASS
- [ ] Dimension 4 Typography: PASS
- [ ] Dimension 5 Spacing: PASS
- [ ] Dimension 6 Registry Safety: PASS

**Approval:** pending
