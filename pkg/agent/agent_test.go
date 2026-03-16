package agent

import (
	"context"
	"testing"
	"orez-crabby/internal/models"
)

type MultiResponseMockProvider struct {
	Responses []*Response
	CallCount int
}

func (p *MultiResponseMockProvider) Name() string { return "multi-mock" }
func (p *MultiResponseMockProvider) Generate(ctx context.Context, messages []Message, tools []ToolDefinition) (*Response, error) {
	if p.CallCount >= len(p.Responses) {
		return &Response{Content: "Done"}, nil
	}
	resp := p.Responses[p.CallCount]
	p.CallCount++
	return resp, nil
}
func (p *MultiResponseMockProvider) Stream(ctx context.Context, messages []Message, tools []ToolDefinition, callback func(*Response)) error {
	resp, _ := p.Generate(ctx, messages, tools)
	callback(resp)
	return nil
}

type MockTool struct {
	Executed bool
}

func (t *MockTool) Name() string        { return "test_tool" }
func (t *MockTool) Description() string { return "A test tool" }
func (t *MockTool) InputSchema() interface{} { return nil }
func (t *MockTool) Execute(ctx context.Context, input string) (string, error) {
	t.Executed = true
	return "Tool executed successfully", nil
}
func (t *MockTool) RequiresApproval() bool { return false }

func TestAgentToolExecution(t *testing.T) {
	mockTool := &MockTool{}
	mockProvider := &MultiResponseMockProvider{
		Responses: []*Response{
			{
				Content: "I will call the tool.",
				ToolCalls: []ToolCall{
					{ID: "1", ToolName: "test_tool", ToolInput: "{}"},
				},
			},
			{
				Content: "Tool completed.",
			},
		},
	}
	
	a := NewAgent(mockProvider)
	a.RegisterTool(mockTool)

	sessionID := "test-session"
	userInput := "Execute the test tool"
	
	steps := []models.Step{}
	err := a.Run(context.Background(), sessionID, userInput, func(s models.Step) {
		steps = append(steps, s)
	})

	if err != nil {
		t.Fatalf("Agent run failed: %v", err)
	}

	if !mockTool.Executed {
		t.Fatal("Expected tool to be executed, but it wasn't")
	}

	foundToolCall := false
	for _, s := range steps {
		if s.Kind == models.StepKindToolCall && s.ToolName == "test_tool" {
			foundToolCall = true
			break
		}
	}

	if !foundToolCall {
		t.Fatal("Expected to find a tool call step in the timeline")
	}
}
