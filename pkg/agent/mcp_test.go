package agent

import (
	"testing"
	"orez-crabby/pkg/config"
)

func TestMcpManager(t *testing.T) {
	m := NewMcpManager()
	if m == nil {
		t.Fatal("Expected NewMcpManager to return a manager, got nil")
	}

	// Test adding a server (mocking transport logic is hard without actual MCP server)
	// But we can test the internal state management.
	
	m.mu.Lock()
	m.configs["test-server"] = config.MCPConfig{Name: "test-server", Type: "stdio"}
	m.mu.Unlock()

	servers := m.ListServers()
	if len(servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(servers))
	}

	if servers[0].Name != "test-server" {
		t.Errorf("Expected server name 'test-server', got '%s'", servers[0].Name)
	}

	m.RemoveServer("test-server")
	servers = m.ListServers()
	if len(servers) != 0 {
		t.Errorf("Expected 0 servers after removal, got %d", len(servers))
	}
}

func TestMCPToolApproval(t *testing.T) {
	tool := &MCPTool{name: "write_file"}
	if !tool.RequiresApproval() {
		t.Error("Expected write_file tool to require approval")
	}

	tool2 := &MCPTool{name: "read_file"}
	if tool2.RequiresApproval() {
		t.Error("Expected read_file tool not to require approval")
	}
}
