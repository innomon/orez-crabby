# Implementation Plan: OpenWork-Go (Wails Edition)

## Phase 1: Foundation (The Shell)
Initialize the project structure, Wails bindings, and the basic UI layout.

### Tasks:
- [x] **1.1. Project Setup:**
    - Initialize Wails with React + TypeScript template.
    - Setup Tailwind CSS.
    - Configure Go project structure (cmd, internal, pkg).
- [x] **1.2. Core Data Layer:**
    - Integrate SQLite (modernc.org).
    - Create schemas for Sessions, Steps, and Config.
- [x] **1.3. Basic Layout:**
    - Create a sidebar (Sessions) and main chat/timeline area.
    - Implement a "Settings" modal (API Keys, Model Selection). (Basic UI implemented, logic deferred to Phase 2)


## Phase 2: Orchestration (The Brain)
Implement the core agent logic and LLM provider integration.

### Tasks:
- [x] **2.1. Provider Interface:**
    - Create a generic Go interface for LLM calls (Prompt, Stream, ToolCall).
    - Implement Ollama provider.
- [x] **2.2. The Planning Loop:**
    - Implement a "Thought" generator (CoT) via system prompt.
    - Create a state machine (Planning -> Executing -> Reflecting -> Completed) in `agent.go`.

- [x] **2.3. Streaming Events:**
    - Setup Wails Event emitters for "agent:step" (incorporates thought, action, and response).


## Phase 3: Tools & Permissions (The Hands)
Build the tool execution engine and the human-in-the-loop interceptor.

### Tasks:
- [x] **3.1. Tool Framework:**
    - Define a `Tool` interface in Go (Name, Description, InputSchema, Execute).
    - Implement "ReadFile" and "ShellCommand" tools.
- [x] **3.2. Permission Interceptor:**
    - Implement a Go mechanism to pause execution.
    - Setup backend logic for "needs-approval" events and response handling.
- [x] **3.3. File Sandbox:**
    - Enforce a "Workspace Root" for all file-related tools using path resolution.


## Phase 4: The Execution Timeline (The UI)
Enhance the UI to support the non-linear "step" visualization.

### Tasks:
- [x] **4.1. Timeline UI Components:**
    - Port the clustering logic (`clusterSteps`) to the React frontend.
    - Implement `StepCard` and `Timeline` components.
    - Set up state management for real-time agent events.
- [x] **4.2. Tool Visualization:**
    - Create specific UIs for "Exploration" vs "Execution" tool outputs.
    - Implement the "Step Grouping" view to collapse exploration sequences.
- [x] **4.3. Log Persistence:**
    - Ensure frontend correctly handles incoming steps and updates based on unique IDs.



## Phase 5: Advanced MCP & Configuration
Implement robust MCP lifecycle management and a centralized configuration system.

### Tasks:
- [ ] **5.1. Advanced MCP Management:**
    - Port MCP server validation and lifecycle logic (Add/Remove/Reload) from `openwork-reference`.
    - Implement a `McpManager` in Go to handle multiple concurrent MCP server connections.
    - Support for both Stdio and (future) Remote MCP transports.
- [ ] **5.2. Centralized Configuration System:**
    - Implement a `ConfigManager` to read/write `opencode.json` (or `orez.json`) using JSONC for human-readability.
    - Migrate hardcoded provider settings (Ollama URL, Model) to the configuration file.
    - Support for workspace-specific configurations (sandboxing, tool overrides).
- [ ] **5.3. MCP & Provider UI:**
    - Create a "Manage MCP Servers" modal to add, remove, and monitor server status.
    - Build a "Provider Settings" interface to configure multiple LLM backends (Ollama, Anthropic, OpenAI).
    - Add a "Status Bar" component to show connected MCP servers and active LLM status.


## Phase 6: Final Polish & UX
Refine the user experience and ensure system stability.

### Tasks:
- [x] **6.1. Workspace Explorer:**
    - Sidebar file tree for the selected Workspace Directory.
    - Backend methods for directory navigation and file listing.
- [x] **6.2. UI/UX Refinement:**
    - Clean, modern UI using Tailwind v4 and Lucide icons.
    - Integrated real-time agent execution timeline.
    - Sandboxed tool execution with permission interceptor.
- [ ] **6.3. System Stability & Testing:**
    - Implement unit tests for the Go orchestration logic and MCP client.
    - Add error handling for lost MCP connections and LLM timeouts.
    - Final audit of the workspace sandboxing and permission interceptor.


---

## Checklist for Immediate Next Steps:
- [ ] Run `wails init -n orez-crabby -t react-ts` (or similar for existing directory).
- [ ] Verify Tailwind CSS installation.
- [ ] Define the `Step` struct in Go for SQLite persistence.
- [ ] Test a simple "Ollama Ping" from the Frontend via Wails Bindings.
