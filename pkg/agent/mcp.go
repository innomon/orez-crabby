package agent

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

// MCPClient handles communication with an external Model Context Protocol server via stdio.
type MCPClient struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	mu     sync.Mutex
	id     int
}

type jsonRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   interface{}     `json:"error,omitempty"`
	ID      int             `json:"id"`
}

func NewMCPClient(command string, args ...string) (*MCPClient, error) {
	cmd := exec.Command(command, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &MCPClient{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
	}, nil
}

func (c *MCPClient) Call(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	c.mu.Lock()
	c.id++
	id := c.id
	c.mu.Unlock()

	req := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      id,
	}

	data, _ := json.Marshal(req)
	_, err := fmt.Fprintln(c.stdin, string(data))
	if err != nil {
		return nil, err
	}

	// Read response (line by line for simplicity in stdio)
	scanner := bufio.NewScanner(c.stdout)
	if scanner.Scan() {
		var resp jsonRPCResponse
		if err := json.Unmarshal(scanner.Bytes(), &resp); err != nil {
			return nil, err
		}
		if resp.Error != nil {
			return nil, fmt.Errorf("MCP error: %v", resp.Error)
		}
		return resp.Result, nil
	}

	return nil, fmt.Errorf("no response from MCP server")
}

func (c *MCPClient) Close() error {
	c.stdin.Close()
	return c.cmd.Wait()
}
