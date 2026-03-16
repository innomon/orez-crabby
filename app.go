package main

import (
	"context"
	"database/sql"
	"fmt"
	"orez-crabby/internal/db"
	"orez-crabby/internal/models"
	"orez-crabby/pkg/agent"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx   context.Context
	db    *sql.DB
	agent *agent.Agent
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize DB
	database, err := db.InitDB()
	if err != nil {
		fmt.Printf("Error initializing DB: %v\n", err)
		return
	}
	a.db = database

	// Initialize Agent with Ollama (Llama3 as default)
	provider := agent.NewOllamaProvider("http://localhost:11434", "llama3")
	a.agent = agent.NewAgent(provider)
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
		ID:        "ws_" + fmt.Sprintf("%d", len(path)), // Simple ID for now
		Name:      filepath.Base(path),
		Path:      path,
		CreatedAt: time.Now(),
	}

	// Register tools with the new workspace path
	a.agent.RegisterTool(&agent.ReadFileTool{BaseTool: agent.BaseTool{WorkspacePath: path}})
	a.agent.RegisterTool(&agent.ShellTool{BaseTool: agent.BaseTool{WorkspacePath: path}})

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
