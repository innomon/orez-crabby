package agent

import (
	"context"
)

// Tool defines the interface for all agent capabilities.
type Tool interface {
	Name() string
	Description() string
	InputSchema() interface{}
	// Execute performs the tool's action. 
	// The return value should be a string (markdown or JSON) representing the result.
	Execute(ctx context.Context, input string) (string, error)
	// RequiresApproval returns true if this tool needs user consent before execution.
	RequiresApproval() bool
}

// ToolRegistry manages available tools.
type ToolRegistry struct {
	tools map[string]Tool
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{tools: make(map[string]Tool)}
}

func (r *ToolRegistry) Register(t Tool) {
	r.tools[t.Name()] = t
}

func (r *ToolRegistry) GetTool(name string) Tool {
	return r.tools[name]
}

func (r *ToolRegistry) ListTools() []ToolDefinition {
	defs := make([]ToolDefinition, 0, len(r.tools))
	for _, t := range r.tools {
		defs = append(defs, ToolDefinition{
			Name:        t.Name(),
			Description: t.Description(),
			InputSchema: t.InputSchema(),
		})
	}
	return defs
}
