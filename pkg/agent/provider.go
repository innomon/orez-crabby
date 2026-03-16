package agent

import (
	"context"
)

// Message represents a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ToolDefinition defines a tool that the model can call.
type ToolDefinition struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"input_schema"`
}

// ToolCall represents a request from the model to call a tool.
type ToolCall struct {
	ID        string `json:"id"`
	ToolName  string `json:"tool_name"`
	ToolInput string `json:"tool_input"` // JSON string
}

// Response represents a completion from the LLM.
type Response struct {
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// Provider defines the interface for LLM backends (Ollama, OpenAI, etc.)
type Provider interface {
	Name() string
	Generate(ctx context.Context, messages []Message, tools []ToolDefinition) (*Response, error)
	Stream(ctx context.Context, messages []Message, tools []ToolDefinition, callback func(*Response)) error
}
