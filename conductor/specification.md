# Specification: OpenWork-Go (Wails Edition)

## 1. Overview
OpenWork-Go is a local-first, open-source AI agent platform built with **Go** and **Wails**. It serves as a desktop-native orchestrator for local and remote LLMs, focusing on "agentic" workflows, a step-by-step execution timeline, and a robust, human-in-the-loop permission system.

## 2. Goals
- **Local-First:** All sensitive data (files, secrets, agent logs) stays on the user's machine.
- **Agentic UX:** A non-linear chat interface focused on "actions" and "steps."
- **Extensible:** A plugin-friendly architecture for "skills" (tools the agent can use).
- **Safe:** Explicit user consent for every tool execution (unless pre-authorized).

## 3. Core Components

### 3.1. Desktop Shell (Wails)
- **Frontend:** React + Tailwind CSS (bundled with Wails).
- **Backend Bridge:** Go methods bound to the frontend for system-level access.
- **System Integration:** Native file pickers, tray support, and window management.

### 3.2. Agent Orchestrator (Go)
- **Session Management:** Persisting agent conversations and execution logs in a local SQLite database.
- **Provider Interface:** Unified API for LLM providers (Ollama, OpenAI, Anthropic).
- **Workflow Engine:** A state machine that manages the "Plan -> Execute -> Reflect" loop, mirroring the logic found in the reference orchestrator.
- **Clustering Logic:** Categorize agent steps into "Exploration" (read, glob, grep) and "Execution" (write, shell) for optimized UI rendering.

### 3.3. Tool & Skill System (Go Plugins / MCP)
- **Built-in Tools:** 
    - *Exploration:* File system (Read/Glob/Grep), Search.
    - *Execution:* File system (Write/Delete), Shell (restricted).
- **External Tools:** Native support for the **Model Context Protocol (MCP)** via stdio or SSE.
- **Permission Interceptor:** A middleware that halts execution to request user approval for specific "Execution" tools.

### 3.4. The Execution Timeline (Frontend)
- **Step Cards:** Instead of simple text bubbles, each "step" is an interactive card.
- **Grouped Views:** Implement the clustering logic from the reference (`groupMessageParts`) to collapse multiple exploration steps into a single "Step Group" to reduce clutter.
- **Real-time Streaming:** Using Wails Events to push agent "thoughts" and "tool outputs" to the UI as they happen.
- **Virtualization:** Use `@tanstack/react-virtual` to handle long execution timelines efficiently.

## 4. Key Features

### 4.1. Workspace Context
- The user selects a "Workspace Directory."
- The agent's `cwd` is locked to this directory.
- All file operations are relative and sandboxed (by convention/configuration).

### 4.2. Human-in-the-Loop (HITL)
- **Permissions:** "Allow Once," "Allow Always for this Session," "Always Allow for this Tool," or "Deny."
- **Diff View:** For file writes, show a diff before applying changes.

### 4.3. Multi-Model Support
- **Ollama Integration:** Auto-detect local Ollama instances.
- **Custom Config:** Manage API keys for remote models (OpenAI, Anthropic).

## 5. Technology Stack
- **Backend:** Go 1.22+
- **Frontend:** React, TypeScript, Tailwind CSS, Lucide Icons.
- **UI Framework:** Wails v2/v3.
- **Database:** SQLite (via `modernc.org/sqlite` for CGO-free portability).
- **LLM Client:** `langchaingo` or custom provider implementations.
- **Communication:** Wails `Events` and `Bindings`.

## 6. Success Metrics
- **Performance:** App launch < 1s.
- **Transparency:** Every agent action is clearly logged in the timeline.
- **Security:** Zero un-prompted file writes to sensitive directories.
