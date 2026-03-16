# Project Mandates: Orez-Crabby (Wails Edition)

This file serves as the source of truth for all AI agent interactions within this workspace. Instructions here take precedence over general defaults.
## 1. Architectural Patterns
- **Frontend-Backend Bridge:** Use Wails `Bindings` for all system-level requests. Prefer Go for heavy logic, using React only for UI state and rendering.
- **Event-Driven UI:** All long-running agent tasks (LLM calls, tool execution) MUST communicate progress via Wails `Events`. Do not use long-polling or simple request-response for streaming data.
- **Tool Clustering:** Explicitly differentiate "Exploration" (read-only) tools from "Execution" (modifying) tools to enable appropriate UI grouping and security gating.
- **CGO-Free:** Prioritize pure Go implementations where possible to simplify cross-compilation. Specifically, use `modernc.org/sqlite` for the database.

## 2. Engineering Standards
- **Naming:** Follow idiomatic Go naming conventions (PascalCase for exports, camelCase for internal).
- **Benchmark:** Use `openwork-reference` as the gold standard for UI/UX patterns (Timeline grouping, Step Cards).
- **Tooling:** Always check for `go fmt` and `npm run lint` before finalizing changes.
- **Security:** 
...
    - Never hardcode API keys. Use an encrypted or local-only configuration file managed by the Go backend.
    - All file operations MUST be scoped to the user-selected "Workspace Directory."
    - Shell tool execution MUST require explicit user confirmation through the Wails UI.

## 3. Tech Stack Preferences
- **Backend:** Go 1.22+
- **Frontend:** React (TypeScript), Tailwind CSS, Lucide Icons.
- **Database:** SQLite (via modernc.org).
- **LLM Integration:** Custom provider interfaces or `langchaingo`.
- **MCP Protocol:** Always use [The official Go SDK for Model Context Protocol servers and clients](https://github.com/modelcontextprotocol/go-sdk)

## 4. Development Workflow
- **Validation:** Every backend change must be accompanied by a Go unit test (`_test.go`). 
- **Conductor:** Refer to `conductor/index.md` for the current roadmap and implementation status.
- **Wails v2/v3:** Adhere to the specific Wails version patterns detected in the `wails.json` configuration.

## 5. Security Protocol (HITL)
The "Human-in-the-Loop" (HITL) system is non-negotiable. No agent action involving File I/O (Write/Delete) or Network access (outside of LLM providers) should execute without a registered `ApprovalEvent` being sent to and confirmed by the frontend.
