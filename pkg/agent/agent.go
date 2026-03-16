package agent

import (
	"context"
	"fmt"
	"sync"

	"orez-crabby/internal/models"
)

const systemPrompt = `You are a helpful desktop assistant. 
Follow this cycle for every task:
1. **Plan**: Analyze the user's request and think about what tools you need.
2. **Execute**: Call the necessary tools.
3. **Reflect**: Look at the tool outputs and decide if the task is complete or if more steps are needed.

Always explain your reasoning in the "thought" phase before calling any tools.
`

type AgentState string

const (
	StateIdle      AgentState = "idle"
	StatePlanning  AgentState = "planning"
	StateExecuting AgentState = "executing"
	StateReflecting AgentState = "reflecting"
	StateCompleted  AgentState = "completed"
)

// Agent manages the interaction between the LLM and the tools.
type Agent struct {
	provider    Provider
	registry    *ToolRegistry
	sessions    map[string]*models.Session
	workspaceID string
	mu          sync.RWMutex

	// PendingApproval stores channels for tools waiting for user consent.
	pendingApprovals map[string]chan bool
}

func NewAgent(provider Provider) *Agent {
	return &Agent{
		provider:         provider,
		registry:         NewToolRegistry(),
		sessions:         make(map[string]*models.Session),
		pendingApprovals: make(map[string]chan bool),
	}
}

func (a *Agent) RegisterTool(t Tool) {
	a.registry.Register(t)
}

func (a *Agent) HandleApproval(id string, approved bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if ch, ok := a.pendingApprovals[id]; ok {
		ch <- approved
		close(ch)
		delete(a.pendingApprovals, id)
	}
}

// Run executes the agentic loop with explicit state transitions.
func (a *Agent) Run(ctx context.Context, sessionID string, userInput string, onStep func(models.Step)) error {
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userInput},
	}
	
	state := StatePlanning

	for state != StateCompleted {
		onStep(models.Step{
			Kind:    models.StepKindThought,
			Content: fmt.Sprintf("Transitioning to state: %s", state),
		})

		// 1. Generate next action (Thought or Tool Call)
		resp, err := a.provider.Generate(ctx, messages, a.registry.ListTools())
		if err != nil {
			state = StateCompleted
			return err
		}

		// Handle text response as a "Thought" or "Final Response"
		if resp.Content != "" {
			kind := models.StepKindThought
			if len(resp.ToolCalls) == 0 {
				kind = models.StepKindResponse
				state = StateCompleted
			} else {
				state = StateExecuting
			}

			onStep(models.Step{
				Kind:    kind,
				Content: resp.Content,
			})
			messages = append(messages, Message{Role: "assistant", Content: resp.Content})
		}

		// 2. Handle Tool Execution
		if len(resp.ToolCalls) > 0 {
			state = StateExecuting
			for _, tc := range resp.ToolCalls {
				tool := a.registry.GetTool(tc.ToolName)
				if tool == nil {
					errResult := fmt.Sprintf("Error: Tool %s not found", tc.ToolName)
					messages = append(messages, Message{Role: "user", Content: errResult})
					continue
				}

				// Permission Interceptor
				if tool.RequiresApproval() {
					onStep(models.Step{
						ID:        tc.ID,
						Kind:      models.StepKindToolCall,
						ToolName:  tc.ToolName,
						ToolInput: tc.ToolInput,
						Content:   "WAITING_FOR_APPROVAL",
					})

					approvalID := fmt.Sprintf("%s_%s", sessionID, tc.ID)
					ch := make(chan bool)
					a.mu.Lock()
					a.pendingApprovals[approvalID] = ch
					a.mu.Unlock()

					approved := <-ch
					if !approved {
						messages = append(messages, Message{Role: "user", Content: "Tool execution denied by user."})
						continue
					}
				}

				onStep(models.Step{
					ID:        tc.ID,
					Kind:      models.StepKindToolCall,
					ToolName:  tc.ToolName,
					ToolInput: tc.ToolInput,
					Content:   "EXECUTING",
				})

				result, err := tool.Execute(ctx, tc.ToolInput)
				if err != nil {
					result = fmt.Sprintf("Error: %v", err)
				}

				onStep(models.Step{
					ID:        tc.ID,
					Kind:      models.StepKindToolCall,
					ToolName:  tc.ToolName,
					ToolInput: tc.ToolInput,
					Content:   result,
				})

				messages = append(messages, Message{Role: "user", Content: fmt.Sprintf("Tool %s result: %s", tc.ToolName, result)})
			}
			
			// After tool execution, transition to reflection
			state = StateReflecting
		} else if state != StateCompleted {
			// No tools and not completed? Force reflection or completion
			state = StateCompleted
		}

		// Limit loop to prevent infinite runs (basic safety)
		if len(messages) > 20 {
			onStep(models.Step{
				Kind:    models.StepKindError,
				Content: "Maximum conversation depth reached. Stopping for safety.",
			})
			state = StateCompleted
		}
	}
	
	return nil
}
