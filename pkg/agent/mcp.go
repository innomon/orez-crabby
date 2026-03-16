package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"orez-crabby/pkg/config"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// McpManager manages multiple concurrent MCP server connections.
type McpManager struct {
	client   *mcp.Client
	sessions map[string]*mcp.ClientSession
	configs  map[string]config.MCPConfig
	mu       sync.RWMutex
}

func NewMcpManager() *McpManager {
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "Orez-Crabby",
		Version: "0.1.0",
	}, nil)

	return &McpManager{
		client:   client,
		sessions: make(map[string]*mcp.ClientSession),
		configs:  make(map[string]config.MCPConfig),
	}
}

func (m *McpManager) AddServer(ctx context.Context, cfg config.MCPConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.sessions[cfg.Name]; ok {
		return fmt.Errorf("server %s already exists", cfg.Name)
	}

	var transport mcp.Transport
	if cfg.Type == "stdio" {
		cmd := exec.Command(cfg.Command, cfg.Args...)
		cmd.Env = os.Environ()
		for k, v := range cfg.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
		transport = &mcp.CommandTransport{
			Command: cmd,
		}
	} else if cfg.Type == "sse" {
		transport = &mcp.SSEClientTransport{
			Endpoint: cfg.URL,
		}
	} else {
		return fmt.Errorf("unsupported MCP type: %s", cfg.Type)
	}

	session, err := m.client.Connect(ctx, transport, nil)
	if err != nil {
		return err
	}

	m.sessions[cfg.Name] = session
	m.configs[cfg.Name] = cfg
	return nil
}

func (m *McpManager) RemoveServer(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, ok := m.sessions[name]; ok {
		session.Close()
		delete(m.sessions, name)
	}
	delete(m.configs, name)
	return nil
}

func (m *McpManager) ListServers() []config.MCPConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var configs []config.MCPConfig
	for _, cfg := range m.configs {
		configs = append(configs, cfg)
	}
	return configs
}

func (m *McpManager) GetSession(name string) *mcp.ClientSession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[name]
}

// MCPTool is an adapter that makes an MCP tool look like an agent.Tool.
type MCPTool struct {
	session     *mcp.ClientSession
	name        string
	description string
	inputSchema interface{}
}

func (t *MCPTool) Name() string        { return t.name }
func (t *MCPTool) Description() string { return t.description }
func (t *MCPTool) InputSchema() interface{} { return t.inputSchema }
func (t *MCPTool) RequiresApproval() bool {
	// For now, assume execution tools require approval.
	// We could base this on tool name or metadata.
	return strings.HasPrefix(t.name, "write_") || strings.HasPrefix(t.name, "execute_")
}

func (t *MCPTool) Execute(ctx context.Context, input string) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return "", err
	}

	result, err := t.session.CallTool(ctx, &mcp.CallToolParams{
		Name:      t.name,
		Arguments: params,
	})
	if err != nil {
		return "", err
	}

	if result.IsError {
		return "Error: tool execution failed", nil
	}

	var sb strings.Builder
	for _, content := range result.Content {
		switch c := content.(type) {
		case *mcp.TextContent:
			sb.WriteString(c.Text)
		case *mcp.ImageContent:
			sb.WriteString(fmt.Sprintf("[Image: %s]", c.MIMEType))
		case *mcp.EmbeddedResource:
			sb.WriteString(fmt.Sprintf("[Resource: %s]", c.Resource.URI))
		}
	}

	return sb.String(), nil
}

// DiscoverTools fetches available tools from an MCP session and returns them as agent.Tool.
func DiscoverTools(ctx context.Context, session *mcp.ClientSession) ([]Tool, error) {
	res, err := session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		return nil, err
	}

	var tools []Tool
	for _, t := range res.Tools {
		tools = append(tools, &MCPTool{
			session:     session,
			name:        t.Name,
			description: t.Description,
			inputSchema: t.InputSchema,
		})
	}

	return tools, nil
}
