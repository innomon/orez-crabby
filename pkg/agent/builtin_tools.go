package agent

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// BaseTool provides common functionality for tools.
type BaseTool struct {
	WorkspacePath string
}

func (b *BaseTool) resolve(path string) (string, error) {
	abs, err := filepath.Abs(filepath.Join(b.WorkspacePath, path))
	if err != nil {
		return "", err
	}
	// Basic sandbox check
	rel, err := filepath.Rel(b.WorkspacePath, abs)
	if err != nil || (len(rel) >= 2 && rel[:2] == "..") {
		return "", fmt.Errorf("path %s is outside workspace", path)
	}
	return abs, nil
}

// ReadFileTool reads the content of a file.
type ReadFileTool struct {
	BaseTool
}

func (t *ReadFileTool) Name() string        { return "read_file" }
func (t *ReadFileTool) Description() string { return "Read the content of a file in the workspace" }
func (t *ReadFileTool) InputSchema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]string{"type": "string"},
		},
		"required": []string{"path"},
	}
}
func (t *ReadFileTool) RequiresApproval() bool { return false }
func (t *ReadFileTool) Execute(ctx context.Context, input string) (string, error) {
	// Parse input (simplified for now)
	var args struct{ Path string }
	// We'll need a better JSON parser here in a real implementation
	path, err := t.resolve(args.Path) // Placeholder for real parsing
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(path)
	return string(content), err
}

// ShellTool executes a command in the shell.
type ShellTool struct {
	BaseTool
}

func (t *ShellTool) Name() string        { return "shell" }
func (t *ShellTool) Description() string { return "Execute a command in the shell" }
func (t *ShellTool) InputSchema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": map[string]string{"type": "string"},
		},
		"required": []string{"command"},
	}
}
func (t *ShellTool) RequiresApproval() bool { return true }
func (t *ShellTool) Execute(ctx context.Context, command string) (string, error) {
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	cmd.Dir = t.WorkspacePath
	out, err := cmd.CombinedOutput()
	return string(out), err
}
