package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaProvider struct {
	BaseURL string
	Model   string
}

func NewOllamaProvider(baseURL, model string) *OllamaProvider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &OllamaProvider{BaseURL: baseURL, Model: model}
}

func (p *OllamaProvider) Name() string {
	return "ollama"
}

type ollamaRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ollamaResponse struct {
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

func (p *OllamaProvider) Generate(ctx context.Context, messages []Message, tools []ToolDefinition) (*Response, error) {
	url := fmt.Sprintf("%s/api/chat", p.BaseURL)
	
	reqBody, _ := json.Marshal(ollamaRequest{
		Model:    p.Model,
		Messages: messages,
		Stream:   false,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama error (status %d): %s", resp.StatusCode, string(body))
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, err
	}

	return &Response{
		Content: ollamaResp.Message.Content,
	}, nil
}

func (p *OllamaProvider) Stream(ctx context.Context, messages []Message, tools []ToolDefinition, callback func(*Response)) error {
	// Basic implementation of streaming
	url := fmt.Sprintf("%s/api/chat", p.BaseURL)
	
	reqBody, _ := json.Marshal(ollamaRequest{
		Model:    p.Model,
		Messages: messages,
		Stream:   true,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	for {
		var ollamaResp ollamaResponse
		if err := decoder.Decode(&ollamaResp); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		callback(&Response{
			Content: ollamaResp.Message.Content,
		})

		if ollamaResp.Done {
			break
		}
	}

	return nil
}
