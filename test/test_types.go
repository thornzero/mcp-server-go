package main

import "testing"

// MCP JSON-RPC structures
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// Test data structures
type Goal struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Priority  int    `json:"priority"`
	Status    string `json:"status"`
	Notes     string `json:"notes"`
	UpdatedAt string `json:"updated_at"`
}

type GoalsListOutput struct {
	Goals []Goal `json:"goals"`
}

type GoalsAddOutput struct {
	ID int `json:"id"`
}

// Test case structures
type MCPTestCase struct {
	Name        string
	Request     MCPRequest
	ExpectError bool
	Validate    func(t *testing.T, resp *MCPResponse)
}

type ToolTestCase struct {
	Name        string
	ToolName    string
	Arguments   map[string]interface{}
	ExpectError bool
	Validate    func(t *testing.T, resp *MCPResponse)
}
