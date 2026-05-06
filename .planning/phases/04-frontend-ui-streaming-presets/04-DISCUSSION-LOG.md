# Phase 4: Frontend UI, Streaming & Presets - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

## Q1: Framework & Styling Foundation
Options:
- React + Vite with Vanilla CSS
- React + Vite with TailwindCSS
Selected: **React + Vite with Vanilla CSS**

## Q2: Dynamic Form Generation
Options:
- Custom Form Builder mapping to ToolSpec
- Library-based (`react-jsonschema-form`)
Selected: **Custom Form Builder**

## Q3: State Management
Options:
- Native React (`useState` / Context)
- External Store (Zustand / Redux)
Selected: **Native React**

## Q4: Log Streaming (Server-Sent Events)
Options:
- Native `EventSource` API
- Custom `fetch()` stream reader
Selected: **Native EventSource API**
