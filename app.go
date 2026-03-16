package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"orez-crabby/internal/db"
	"orez-crabby/internal/models"
	"orez-crabby/pkg/agent"
	"orez-crabby/pkg/config"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	db            *sql.DB
	agent         *agent.Agent
	mcpManager    *agent.McpManager
	configManager *config.ConfigManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize Config
	cfgMgr, err := config.NewConfigManager("")
	if err != nil {
		fmt.Printf("Error initializing config: %v\n", err)
	} else {
		a.configManager = cfgMgr
		if err := a.configManager.Load(); err != nil {
			fmt.Printf("Error loading config: %v\n", err)
		}
	}

	// Initialize DB
	database, err := db.InitDB()
	if err != nil {
		fmt.Printf("Error initializing DB: %v\n", err)
		return
	}
	a.db = database

	// Initialize McpManager
	a.mcpManager = agent.NewMcpManager()

	// Initialize Agent with configured provider
	cfg := a.configManager.Get()
	var provider agent.Provider
	if cfg.Provider.Name == "ollama" {
		provider = agent.NewOllamaProvider(cfg.Provider.BaseURL, cfg.Provider.Model)
	} else {
		// Fallback to default
		provider = agent.NewOllamaProvider("http://localhost:11434", "llama3")
	}
	a.agent = agent.NewAgent(provider)

	// Auto-connect MCP servers from config
	for _, mcpCfg := range cfg.MCPServers {
		go func(c config.MCPConfig) {
			if err := a.AddMcpServer(c); err != nil {
				fmt.Printf("Failed to auto-connect MCP server %s: %v\n", c.Name, err)
			}
		}(mcpCfg)
	}
}

// RunAgent triggers the agent's orchestration loop for a session
func (a *App) RunAgent(sessionID string, input string) string {
	err := a.agent.Run(a.ctx, sessionID, input, func(step models.Step) {
		// Emit events to frontend as steps occur
		runtime.EventsEmit(a.ctx, "agent:step", step)
	})

	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return "Success"
}

// GetConfig returns the current application configuration
func (a *App) GetConfig() *config.Config {
	return a.configManager.Get()
}

// SetProvider updates the LLM provider settings
func (a *App) SetProvider(name, baseURL, model string) error {
	cfg := a.configManager.Get()
	cfg.Provider.Name = name
	cfg.Provider.BaseURL = baseURL
	cfg.Provider.Model = model

	// Re-initialize agent provider (simplistic for now)
	var provider agent.Provider
	if name == "ollama" {
		provider = agent.NewOllamaProvider(baseURL, model)
	}
	if provider != nil {
		a.agent = agent.NewAgent(provider)
	}

	return a.configManager.Save(cfg)
}

// SelectWorkspace opens a directory picker and sets the active workspace
func (a *App) SelectWorkspace() (models.Workspace, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Workspace Directory",
	})
	if err != nil {
		return models.Workspace{}, err
	}
	if path == "" {
		return models.Workspace{}, fmt.Errorf("no directory selected")
	}

	workspace := models.Workspace{
		ID:        "ws_" + fmt.Sprintf("%d", time.Now().Unix()), // Better ID than length
		Name:      filepath.Base(path),
		Path:      path,
		CreatedAt: time.Now(),
	}

	// Register tools with the new workspace path
	a.agent.RegisterTool(&agent.ReadFileTool{BaseTool: agent.BaseTool{WorkspacePath: path}})
	a.agent.RegisterTool(&agent.ShellTool{BaseTool: agent.BaseTool{WorkspacePath: path}})

	// Save last workspace to config
	cfg := a.configManager.Get()
	cfg.LastWorkspace = path
	a.configManager.Save(cfg)

	return workspace, nil
}

type FileEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
	Path  string `json:"path"`
}

// GetWorkspaceFiles returns a list of files in the given path (relative to workspace)
func (a *App) GetWorkspaceFiles(rootPath string) ([]FileEntry, error) {
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	var files []FileEntry
	for _, entry := range entries {
		if entry.Name()[0] == '.' {
			continue // Skip hidden files
		}
		files = append(files, FileEntry{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Path:  filepath.Join(rootPath, entry.Name()),
		})
	}
	return files, nil
}

// AddMcpServer adds an MCP server and registers its tools
func (a *App) AddMcpServer(cfg config.MCPConfig) error {
	err := a.mcpManager.AddServer(a.ctx, cfg)
	if err != nil {
		return err
	}

	session := a.mcpManager.GetSession(cfg.Name)
	if session == nil {
		return fmt.Errorf("failed to get session for server %s", cfg.Name)
	}

	tools, err := agent.DiscoverTools(a.ctx, session)
	if err != nil {
		return fmt.Errorf("failed to discover tools from %s: %v", cfg.Name, err)
	}

	for _, t := range tools {
		a.agent.RegisterTool(t)
	}

	// Update and save config
	appCfg := a.configManager.Get()
	exists := false
	for i, existing := range appCfg.MCPServers {
		if existing.Name == cfg.Name {
			appCfg.MCPServers[i] = cfg
			exists = true
			break
		}
	}
	if !exists {
		appCfg.MCPServers = append(appCfg.MCPServers, cfg)
	}
	return a.configManager.Save(appCfg)
}

// ListMcpServers returns the list of connected MCP servers
func (a *App) ListMcpServers() []config.MCPConfig {
	return a.mcpManager.ListServers()
}

// RemoveMcpServer removes an MCP server
func (a *App) RemoveMcpServer(name string) error {
	if err := a.mcpManager.RemoveServer(name); err != nil {
		return err
	}

	// Update and save config
	appCfg := a.configManager.Get()
	for i, cfg := range appCfg.MCPServers {
		if cfg.Name == name {
			appCfg.MCPServers = append(appCfg.MCPServers[:i], appCfg.MCPServers[i+1:]...)
			break
		}
	}
	return a.configManager.Save(appCfg)
}
