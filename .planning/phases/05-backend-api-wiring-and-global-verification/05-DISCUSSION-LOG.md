# Phase 05: Backend API Wiring & Global Verification - Discussion Log

**Date:** 2026-05-07

## Server Framework
**Question:** Standard library `net/http` vs router library?
**Options:**
1. Standard library (net/http)
2. Chi router
3. Gin or Echo
4. Other
**User Selected:** 1 (Standard library net/http)

## Configuration Strategy
**Question:** How does the server get its port, DB paths, and other runtime settings?
**Options:**
1. Environment Variables only
2. Environment Variables + .env file support
3. Configuration file (JSON/YAML)
4. CLI Flags
5. Other
**User Selected:** 3 (Configuration file (JSON/YAML))

## Static File Serving
**Question:** Does the Go backend serve the built React SPA (`web/dist`)?
**Options:**
1. Yes, Go serves the SPA
2. No, separate web server
3. Other
**User Selected:** 1 (Yes, Go serves the SPA)
