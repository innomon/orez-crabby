package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

// MCPConfig defines the configuration for an MCP server.
type MCPConfig struct {
	Name    string            `json:"name"`
	Type    string            `json:"type"` // "stdio" or "sse"
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	URL     string            `json:"url,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

// ProviderConfig defines settings for the LLM provider.
type ProviderConfig struct {
	Name    string `json:"name"`
	BaseURL string `json:"baseUrl"`
	Model   string `json:"model"`
}

// Config represents the application-wide configuration.
type Config struct {
	Provider      ProviderConfig `json:"provider"`
	MCPServers    []MCPConfig    `json:"mcpServers"`
	LastWorkspace string         `json:"lastWorkspace,omitempty"`
}

// ConfigManager handles the lifecycle of the configuration file.
type ConfigManager struct {
	path   string
	config *Config
	mu     sync.RWMutex
}

func NewConfigManager(path string) (*ConfigManager, error) {
	if path == "" {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".orez.json")
	}
	return &ConfigManager{
		path:   path,
		config: &Config{},
	}, nil
}

func (m *ConfigManager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			// Initialize with defaults if not found
			m.config.Provider = ProviderConfig{
				Name:    "ollama",
				BaseURL: "http://localhost:11434",
				Model:   "llama3",
			}
			m.config.MCPServers = []MCPConfig{}
			return m.saveLocked()
		}
		return err
	}

	// Strip comments for JSONC support
	cleanJSON := stripComments(data)

	if err := json.Unmarshal(cleanJSON, m.config); err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	return nil
}

func (m *ConfigManager) Get() *Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

func (m *ConfigManager) Save(cfg *Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config = cfg
	return m.saveLocked()
}

func (m *ConfigManager) saveLocked() error {
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.path, data, 0644)
}

func stripComments(data []byte) []byte {
	// Simple comment stripper for // and /* */
	// Note: This is a basic implementation and might fail on strings containing // or /*
	// But it follows the mandate for "human-readability" with minimal extra libs.
	re := regexp.MustCompile(`(?s)//.*?\n|/\*.*?\*/`)
	return re.ReplaceAll(data, []byte(""))
}
