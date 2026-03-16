# Orez-Crabby (Wails Edition)

Orez-Crabby is a local-first, desktop-native AI agent platform built with **Go** and **Wails**. It is a Go-based port of the [OpenWork](https://github.com/different-ai/openwork) platform, originally written in Rust, redesigned to provide a lightweight and extensible Go implementation. It serves as an orchestrator for local and remote LLMs, focusing on "agentic" workflows, a step-by-step execution timeline, and a robust human-in-the-loop permission system.

## 🚀 Features

- **Local-First Agentic Loop:** Implements a "Plan -> Execute -> Reflect" state machine for complex task handling.
- **Model Context Protocol (MCP):** Native support for the official [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk), supporting both **Stdio** and **SSE** transports.
- **Dynamic Tool Discovery:** Automatically fetches and registers tools from any connected MCP server.
- **Human-in-the-Loop (HITL):** Explicit user consent required for "Execution" tools (e.g., shell commands, file writes).
- **Interactive Timeline UI:** A non-linear chat interface using Step Cards to visualize the agent's thoughts and actions.
- **Centralized Configuration:** Managed settings via JSONC (`.orez.json`) for LLM providers, MCP servers, and workspace persistence.
- **Integrated File Explorer:** Context-aware workspace navigation for the agent.
- **Status Bar:** Real-time feedback on connected MCP servers and active LLM status.

## 🛠️ Tech Stack

- **Backend:** Go 1.22+ (Wails v2)
- **Frontend:** React (TypeScript), Tailwind CSS, Lucide Icons
- **Database:** SQLite (via `modernc.org/sqlite` for CGO-free portability)
- **MCP SDK:** Official `modelcontextprotocol/go-sdk`
- **LLM Provider:** Ollama (default), extensible to OpenAI/Anthropic

## 🏃 Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.22+
- [Node.js](https://nodejs.org/) & [NPM](https://www.npmjs.com/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)
- [Ollama](https://ollama.com/) (running locally)

### Installation & Development

1. **Clone the repository:**
   ```bash
   git clone https://github.com/innomon/orez-crabby
   cd orez-crabby
   ```

2. **Run in development mode:**
   ```bash
   wails dev
   ```
   This will start the Go backend and the Vite frontend with hot-reload enabled.

3. **Building for Production:**
   ```bash
   wails build
   ```
   The compiled binary will be located in the `build/bin` directory.

## ⚙️ Configuration

Orez-Crabby stores its configuration in `~/.orez.json`. You can manage settings directly in the app via the **Settings Modal** or by editing the JSONC file:

```jsonc
{
  "provider": {
    "name": "ollama",
    "baseUrl": "http://localhost:11434",
    "model": "llama3"
  },
  "mcpServers": [
    {
      "name": "filesystem",
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-everything"]
    }
  ]
}
```

## 🧪 Testing

Run the Go unit tests for the agent and configuration logic:
```bash
go test -v ./...
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
