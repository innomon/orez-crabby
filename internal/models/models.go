package models

import "time"

// Workspace represents a user-selected project directory.
type Workspace struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

// Session represents a single chat/execution session.
type Session struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	Title       string    `json:"title"`
	Model       string    `json:"model"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StepKind represents the type of a step in the timeline.
type StepKind string

const (
	StepKindThought  StepKind = "thought"
	StepKindToolCall StepKind = "tool_call"
	StepKindResponse StepKind = "response"
	StepKindError    StepKind = "error"
)

// Step represents an individual action or thought by the agent.
type Step struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Kind      StepKind  `json:"kind"`
	Content   string    `json:"content"`    // Markdown content or tool output
	ToolName  string    `json:"tool_name"`  // Optional: for tool_call
	ToolInput string    `json:"tool_input"` // Optional: JSON input
	CreatedAt time.Time `json:"created_at"`
}

// Config represents application-wide settings.
type Config struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
