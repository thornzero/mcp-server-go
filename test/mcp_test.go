package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"
)

// Test server setup
func setupTestServer(t *testing.T) (*exec.Cmd, io.WriteCloser, io.ReadCloser) {
	// Build server
	buildCmd := exec.Command("sh", "-c", "cd .. && go build -o build/project-manager ./cmd/project-manager")
	if err := buildCmd.Run(); err != nil {
		t.Fatal("Failed to build server:", err)
	}

	// Start server
	cmd := exec.Command("../build/project-manager")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal("Failed to create stdin pipe:", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal("Failed to create stdout pipe:", err)
	}

	if err := cmd.Start(); err != nil {
		t.Fatal("Failed to start server:", err)
	}

	// Give server time to start and verify it's running
	time.Sleep(200 * time.Millisecond)

	// Check if process is still running
	if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
		t.Fatal("Server process exited immediately after start")
	}

	return cmd, stdin, stdout
}

func teardownTestServer(cmd *exec.Cmd) {
	if cmd.Process != nil {
		// Try graceful shutdown first
		cmd.Process.Signal(os.Interrupt)

		// Wait a bit for graceful shutdown
		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		select {
		case <-done:
			// Process exited gracefully
		case <-time.After(3 * time.Second):
			// Force kill if it doesn't exit gracefully
			if cmd.Process != nil {
				cmd.Process.Kill()
				cmd.Wait() // Wait for the kill to take effect
			}
		}
	}
}

// Helper functions
func sendRequest(stdin io.WriteCloser, req MCPRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = stdin.Write(append(data, '\n'))
	if err != nil {
		// Check for broken pipe specifically
		if err.Error() == "write |1: broken pipe" {
			return fmt.Errorf("server process terminated unexpectedly (broken pipe)")
		}
	}
	return err
}

func readResponse(stdout io.ReadCloser) (*MCPResponse, error) {
	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	var resp MCPResponse
	err := json.Unmarshal(scanner.Bytes(), &resp)
	return &resp, err
}

// Test tables
func getMCPProtocolTests() []MCPTestCase {
	return []MCPTestCase{
		{
			Name: "Initialize Server",
			Request: MCPRequest{
				JSONRPC: "2.0",
				ID:      1,
				Method:  "initialize",
				Params: map[string]interface{}{
					"protocolVersion": "2024-11-05",
					"capabilities":    map[string]interface{}{},
					"clientInfo": map[string]interface{}{
						"name":    "test-client",
						"version": "1.0.0",
					},
				},
			},
			ExpectError: false,
			Validate: func(t *testing.T, resp *MCPResponse) {
				if resp.Error != nil {
					t.Errorf("Expected no error, got: %v", resp.Error)
				}
			},
		},
		{
			Name: "List Tools",
			Request: MCPRequest{
				JSONRPC: "2.0",
				ID:      2,
				Method:  "tools/list",
			},
			ExpectError: false,
			Validate: func(t *testing.T, resp *MCPResponse) {
				if resp.Error != nil {
					t.Errorf("Expected no error, got: %v", resp.Error)
				}
				// Validate that tools are returned
				if result, ok := resp.Result.(map[string]interface{}); ok {
					if tools, ok := result["tools"].([]interface{}); ok {
						if len(tools) == 0 {
							t.Error("Expected tools to be returned, got empty list")
						}
					} else {
						t.Error("Expected tools array in response")
					}
				} else {
					t.Error("Expected result object in response")
				}
			},
		},
	}
}

func getGoalsToolTests() []ToolTestCase {
	return []ToolTestCase{
		{
			Name:        "List Goals",
			ToolName:    "goals_list",
			Arguments:   map[string]interface{}{"limit": 10},
			ExpectError: false,
			Validate: func(t *testing.T, resp *MCPResponse) {
				if resp.Error != nil {
					t.Logf("Goals list error (may be expected): %v", resp.Error)
					return
				}
				// Parse and validate goals response
				if result, ok := resp.Result.(map[string]interface{}); ok {
					if content, ok := result["content"].([]interface{}); ok && len(content) > 0 {
						if goalsData, ok := content[0].(map[string]interface{}); ok {
							if goals, ok := goalsData["goals"].([]interface{}); ok {
								t.Logf("Found %d goals", len(goals))
							}
						}
					}
				}
			},
		},
		{
			Name:     "Add Goal",
			ToolName: "goals_add",
			Arguments: map[string]interface{}{
				"title":    "Test Goal from Table Test",
				"priority": 50,
				"notes":    "This goal was created by the table-driven test",
			},
			ExpectError: false,
			Validate: func(t *testing.T, resp *MCPResponse) {
				if resp.Error != nil {
					t.Errorf("Expected no error, got: %v", resp.Error)
					return
				}
				// Validate goal was created
				if result, ok := resp.Result.(map[string]interface{}); ok {
					if content, ok := result["content"].([]interface{}); ok && len(content) > 0 {
						if goalData, ok := content[0].(map[string]interface{}); ok {
							if id, ok := goalData["id"].(float64); ok {
								t.Logf("Goal added with ID: %.0f", id)
							}
						}
					}
				}
			},
		},
		{
			Name:     "Add Goal with Invalid Data",
			ToolName: "goals_add",
			Arguments: map[string]interface{}{
				"title":    "", // Empty title should cause error
				"priority": 50,
			},
			ExpectError: true,
			Validate: func(t *testing.T, resp *MCPResponse) {
				// Check if error is in JSON-RPC error field
				if resp.Error != nil {
					return
				}

				// Check if error is in result with isError flag
				if result, ok := resp.Result.(map[string]interface{}); ok {
					if isError, ok := result["isError"].(bool); ok && isError {
						return
					}
				}

				t.Error("Expected error for empty title, got none")
			},
		},
	}
}

func getADRToolTests() []ToolTestCase {
	return []ToolTestCase{
		{
			Name:        "List ADRs",
			ToolName:    "adrs_list",
			Arguments:   map[string]interface{}{},
			ExpectError: false,
			Validate: func(t *testing.T, resp *MCPResponse) {
				if resp.Error != nil {
					t.Logf("ADRs list error (may be expected): %v", resp.Error)
					return
				}
				// Parse and validate ADRs response
				if result, ok := resp.Result.(map[string]interface{}); ok {
					if content, ok := result["content"].([]interface{}); ok && len(content) > 0 {
						if adrsData, ok := content[0].(map[string]interface{}); ok {
							if adrs, ok := adrsData["adrs"].([]interface{}); ok {
								t.Logf("Found %d ADRs", len(adrs))
							}
						}
					}
				}
			},
		},
	}
}

func getCIToolTests() []ToolTestCase {
	return []ToolTestCase{
		{
			Name:        "Run Tests",
			ToolName:    "ci_run_tests",
			Arguments:   map[string]interface{}{"scope": "./..."},
			ExpectError: false,
			Validate: func(t *testing.T, resp *MCPResponse) {
				if resp.Error != nil {
					t.Logf("CI run tests error (may be expected): %v", resp.Error)
					return
				}
				// Parse and validate CI response
				if result, ok := resp.Result.(map[string]interface{}); ok {
					if content, ok := result["content"].([]interface{}); ok && len(content) > 0 {
						if ciData, ok := content[0].(map[string]interface{}); ok {
							if status, ok := ciData["status"].(string); ok {
								t.Logf("Test status: %s", status)
							}
						}
					}
				}
			},
		},
		{
			Name:        "Get Last Failure",
			ToolName:    "ci_last_failure",
			Arguments:   map[string]interface{}{},
			ExpectError: false,
			Validate: func(t *testing.T, resp *MCPResponse) {
				if resp.Error != nil {
					t.Logf("CI last failure error (may be expected): %v", resp.Error)
					return
				}
				// Parse and validate CI failure response
				if result, ok := resp.Result.(map[string]interface{}); ok {
					if content, ok := result["content"].([]interface{}); ok && len(content) > 0 {
						if failureData, ok := content[0].(map[string]interface{}); ok {
							if status, ok := failureData["status"].(string); ok {
								t.Logf("Last failure status: %s", status)
							}
						}
					}
				}
			},
		},
	}
}

func getMarkdownToolTests() []ToolTestCase {
	return []ToolTestCase{
		{
			Name:        "Lint Markdown Files",
			ToolName:    "markdown_lint",
			Arguments:   map[string]interface{}{"path": "README.md"},
			ExpectError: false,
			Validate: func(t *testing.T, resp *MCPResponse) {
				if resp.Error != nil {
					t.Logf("Markdown lint error (may be expected): %v", resp.Error)
					return
				}
				// Parse and validate markdown lint response
				if result, ok := resp.Result.(map[string]interface{}); ok {
					if content, ok := result["content"].([]interface{}); ok && len(content) > 0 {
						if lintData, ok := content[0].(map[string]interface{}); ok {
							if issues, ok := lintData["issues"].([]interface{}); ok {
								t.Logf("Found %d markdown lint issues", len(issues))
							}
						}
					}
				}
			},
		},
		{
			Name:        "Lint Markdown with Fix",
			ToolName:    "markdown_lint",
			Arguments:   map[string]interface{}{"path": "README.md", "fix": true},
			ExpectError: false,
			Validate: func(t *testing.T, resp *MCPResponse) {
				if resp.Error != nil {
					t.Logf("Markdown lint with fix error (may be expected): %v", resp.Error)
					return
				}
				// Parse and validate markdown lint response
				if result, ok := resp.Result.(map[string]interface{}); ok {
					if content, ok := result["content"].([]interface{}); ok && len(content) > 0 {
						if lintData, ok := content[0].(map[string]interface{}); ok {
							if fixed, ok := lintData["fixed"].(bool); ok {
								t.Logf("Markdown lint fixed: %v", fixed)
							}
						}
					}
				}
			},
		},
	}
}

func getTemplateToolTests() []ToolTestCase {
	return []ToolTestCase{
		{
			Name:        "List Templates",
			ToolName:    "template_list",
			Arguments:   map[string]interface{}{},
			ExpectError: false,
			Validate: func(t *testing.T, resp *MCPResponse) {
				if resp.Error != nil {
					t.Logf("Template list error (may be expected): %v", resp.Error)
					return
				}
				// Parse and validate template list response
				if result, ok := resp.Result.(map[string]interface{}); ok {
					if content, ok := result["content"].([]interface{}); ok && len(content) > 0 {
						if templatesData, ok := content[0].(map[string]interface{}); ok {
							if templates, ok := templatesData["templates"].([]interface{}); ok {
								t.Logf("Found %d templates", len(templates))
							}
						}
					}
				}
			},
		},
		{
			Name:     "Register Template",
			ToolName: "template_register",
			Arguments: map[string]interface{}{
				"id":          "test-template",
				"name":        "Test Template",
				"category":    "test",
				"content":     "# {{title}}\n\n{{description}}",
				"description": "A test template for integration testing",
				"variables": []map[string]interface{}{
					{
						"name":          "title",
						"type":          "string",
						"required":      true,
						"default_value": "",
						"description":   "The title of the document",
					},
					{
						"name":          "description",
						"type":          "string",
						"required":      false,
						"default_value": "Default description",
						"description":   "The description of the document",
					},
				},
			},
			ExpectError: false,
			Validate: func(t *testing.T, resp *MCPResponse) {
				if resp.Error != nil {
					t.Errorf("Expected no error, got: %v", resp.Error)
					return
				}
				// Parse and validate template registration response
				if result, ok := resp.Result.(map[string]interface{}); ok {
					if content, ok := result["content"].([]interface{}); ok && len(content) > 0 {
						if templateData, ok := content[0].(map[string]interface{}); ok {
							if id, ok := templateData["id"].(string); ok {
								t.Logf("Template registered with ID: %s", id)
							}
						}
					}
				}
			},
		},
		{
			Name:        "Get Template",
			ToolName:    "template_get",
			Arguments:   map[string]interface{}{"id": "test-template"},
			ExpectError: false,
			Validate: func(t *testing.T, resp *MCPResponse) {
				if resp.Error != nil {
					t.Logf("Template get error (may be expected): %v", resp.Error)
					return
				}
				// Parse and validate template get response
				if result, ok := resp.Result.(map[string]interface{}); ok {
					if content, ok := result["content"].([]interface{}); ok && len(content) > 0 {
						if templateData, ok := content[0].(map[string]interface{}); ok {
							if name, ok := templateData["name"].(string); ok {
								t.Logf("Retrieved template: %s", name)
							}
						}
					}
				}
			},
		},
		{
			Name:     "Apply Template",
			ToolName: "template_apply",
			Arguments: map[string]interface{}{
				"template_id": "test-template",
				"variables": map[string]interface{}{
					"title":       "Test Document",
					"description": "This is a test document created by the integration test",
				},
			},
			ExpectError: false,
			Validate: func(t *testing.T, resp *MCPResponse) {
				if resp.Error != nil {
					t.Logf("Template apply error (may be expected): %v", resp.Error)
					return
				}
				// Parse and validate template apply response
				if result, ok := resp.Result.(map[string]interface{}); ok {
					if content, ok := result["content"].([]interface{}); ok && len(content) > 0 {
						if applyData, ok := content[0].(map[string]interface{}); ok {
							if output, ok := applyData["output"].(string); ok {
								t.Logf("Template applied, output length: %d characters", len(output))
							}
						}
					}
				}
			},
		},
		{
			Name:     "Update Template",
			ToolName: "template_update",
			Arguments: map[string]interface{}{
				"id":          "test-template",
				"name":        "Updated Test Template",
				"description": "An updated test template for integration testing",
			},
			ExpectError: false,
			Validate: func(t *testing.T, resp *MCPResponse) {
				if resp.Error != nil {
					t.Logf("Template update error (may be expected): %v", resp.Error)
					return
				}
				// Parse and validate template update response
				if result, ok := resp.Result.(map[string]interface{}); ok {
					if content, ok := result["content"].([]interface{}); ok && len(content) > 0 {
						if templateData, ok := content[0].(map[string]interface{}); ok {
							if name, ok := templateData["name"].(string); ok {
								t.Logf("Template updated: %s", name)
							}
						}
					}
				}
			},
		},
		{
			Name:        "Delete Template",
			ToolName:    "template_delete",
			Arguments:   map[string]interface{}{"id": "test-template"},
			ExpectError: false,
			Validate: func(t *testing.T, resp *MCPResponse) {
				if resp.Error != nil {
					t.Logf("Template delete error (may be expected): %v", resp.Error)
					return
				}
				// Parse and validate template delete response
				if result, ok := resp.Result.(map[string]interface{}); ok {
					if content, ok := result["content"].([]interface{}); ok && len(content) > 0 {
						if deleteData, ok := content[0].(map[string]interface{}); ok {
							if deleted, ok := deleteData["deleted"].(bool); ok {
								t.Logf("Template deleted: %v", deleted)
							}
						}
					}
				}
			},
		},
	}
}

// Main test functions
func TestMCPProtocol(t *testing.T) {
	cmd, stdin, stdout := setupTestServer(t)
	defer teardownTestServer(cmd)

	tests := getMCPProtocolTests()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if err := sendRequest(stdin, tt.Request); err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}

			resp, err := readResponse(stdout)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			if tt.ExpectError {
				// Check if error is in JSON-RPC error field
				if resp.Error != nil {
					// Error found in JSON-RPC error field
				} else if result, ok := resp.Result.(map[string]interface{}); ok {
					if isError, ok := result["isError"].(bool); ok && isError {
						// Error found in result with isError flag
					} else {
						t.Error("Expected error but got none")
					}
				} else {
					t.Error("Expected error but got none")
				}
			} else {
				if resp.Error != nil {
					t.Errorf("Unexpected error: %v", resp.Error)
				}
			}

			if tt.Validate != nil {
				tt.Validate(t, resp)
			}
		})
	}
}

func TestGoalsTools(t *testing.T) {
	cmd, stdin, stdout := setupTestServer(t)
	defer teardownTestServer(cmd)

	// Initialize server first
	initReq := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	if err := sendRequest(stdin, initReq); err != nil {
		t.Fatal("Failed to initialize server:", err)
	}

	// Read init response
	_, err := readResponse(stdout)
	if err != nil {
		t.Fatal("Failed to read init response:", err)
	}

	tests := getGoalsToolTests()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := MCPRequest{
				JSONRPC: "2.0",
				ID:      len(tests) + 10, // Ensure unique ID
				Method:  "tools/call",
				Params: ToolCallParams{
					Name:      tt.ToolName,
					Arguments: tt.Arguments,
				},
			}

			if err := sendRequest(stdin, req); err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}

			resp, err := readResponse(stdout)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			if tt.ExpectError {
				// Check if error is in JSON-RPC error field
				if resp.Error != nil {
					// Error found in JSON-RPC error field
				} else if result, ok := resp.Result.(map[string]interface{}); ok {
					if isError, ok := result["isError"].(bool); ok && isError {
						// Error found in result with isError flag
					} else {
						t.Error("Expected error but got none")
					}
				} else {
					t.Error("Expected error but got none")
				}
			} else {
				if resp.Error != nil {
					t.Errorf("Unexpected error: %v", resp.Error)
				}
			}

			if tt.Validate != nil {
				tt.Validate(t, resp)
			}
		})
	}
}

func TestADRTools(t *testing.T) {
	cmd, stdin, stdout := setupTestServer(t)
	defer teardownTestServer(cmd)

	// Initialize server first
	initReq := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	if err := sendRequest(stdin, initReq); err != nil {
		t.Fatal("Failed to initialize server:", err)
	}

	// Read init response
	_, err := readResponse(stdout)
	if err != nil {
		t.Fatal("Failed to read init response:", err)
	}

	tests := getADRToolTests()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := MCPRequest{
				JSONRPC: "2.0",
				ID:      len(tests) + 20, // Ensure unique ID
				Method:  "tools/call",
				Params: ToolCallParams{
					Name:      tt.ToolName,
					Arguments: tt.Arguments,
				},
			}

			if err := sendRequest(stdin, req); err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}

			resp, err := readResponse(stdout)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			if tt.ExpectError {
				// Check if error is in JSON-RPC error field
				if resp.Error != nil {
					// Error found in JSON-RPC error field
				} else if result, ok := resp.Result.(map[string]interface{}); ok {
					if isError, ok := result["isError"].(bool); ok && isError {
						// Error found in result with isError flag
					} else {
						t.Error("Expected error but got none")
					}
				} else {
					t.Error("Expected error but got none")
				}
			} else {
				if resp.Error != nil {
					t.Errorf("Unexpected error: %v", resp.Error)
				}
			}

			if tt.Validate != nil {
				tt.Validate(t, resp)
			}
		})
	}
}

func TestCITools(t *testing.T) {
	cmd, stdin, stdout := setupTestServer(t)
	defer teardownTestServer(cmd)

	// Initialize server first
	initReq := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	if err := sendRequest(stdin, initReq); err != nil {
		t.Fatal("Failed to initialize server:", err)
	}

	// Read init response
	_, err := readResponse(stdout)
	if err != nil {
		t.Fatal("Failed to read init response:", err)
	}

	tests := getCIToolTests()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := MCPRequest{
				JSONRPC: "2.0",
				ID:      len(tests) + 30, // Ensure unique ID
				Method:  "tools/call",
				Params: ToolCallParams{
					Name:      tt.ToolName,
					Arguments: tt.Arguments,
				},
			}

			if err := sendRequest(stdin, req); err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}

			resp, err := readResponse(stdout)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			if tt.ExpectError {
				// Check if error is in JSON-RPC error field
				if resp.Error != nil {
					// Error found in JSON-RPC error field
				} else if result, ok := resp.Result.(map[string]interface{}); ok {
					if isError, ok := result["isError"].(bool); ok && isError {
						// Error found in result with isError flag
					} else {
						t.Error("Expected error but got none")
					}
				} else {
					t.Error("Expected error but got none")
				}
			} else {
				if resp.Error != nil {
					t.Errorf("Unexpected error: %v", resp.Error)
				}
			}

			if tt.Validate != nil {
				tt.Validate(t, resp)
			}
		})
	}
}

func TestMarkdownTools(t *testing.T) {
	cmd, stdin, stdout := setupTestServer(t)
	defer teardownTestServer(cmd)

	// Initialize server first
	initReq := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	if err := sendRequest(stdin, initReq); err != nil {
		t.Fatal("Failed to initialize server:", err)
	}

	// Read init response
	_, err := readResponse(stdout)
	if err != nil {
		t.Fatal("Failed to read init response:", err)
	}

	tests := getMarkdownToolTests()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := MCPRequest{
				JSONRPC: "2.0",
				ID:      len(tests) + 40, // Ensure unique ID
				Method:  "tools/call",
				Params: ToolCallParams{
					Name:      tt.ToolName,
					Arguments: tt.Arguments,
				},
			}

			if err := sendRequest(stdin, req); err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}

			resp, err := readResponse(stdout)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			if tt.ExpectError {
				// Check if error is in JSON-RPC error field
				if resp.Error != nil {
					// Error found in JSON-RPC error field
				} else if result, ok := resp.Result.(map[string]interface{}); ok {
					if isError, ok := result["isError"].(bool); ok && isError {
						// Error found in result with isError flag
					} else {
						t.Error("Expected error but got none")
					}
				} else {
					t.Error("Expected error but got none")
				}
			} else {
				if resp.Error != nil {
					t.Errorf("Unexpected error: %v", resp.Error)
				}
			}

			if tt.Validate != nil {
				tt.Validate(t, resp)
			}
		})
	}
}

func TestTemplateTools(t *testing.T) {
	cmd, stdin, stdout := setupTestServer(t)
	defer teardownTestServer(cmd)

	// Initialize server first
	initReq := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	if err := sendRequest(stdin, initReq); err != nil {
		t.Fatal("Failed to initialize server:", err)
	}

	// Read init response
	_, err := readResponse(stdout)
	if err != nil {
		t.Fatal("Failed to read init response:", err)
	}

	tests := getTemplateToolTests()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := MCPRequest{
				JSONRPC: "2.0",
				ID:      len(tests) + 50, // Ensure unique ID
				Method:  "tools/call",
				Params: ToolCallParams{
					Name:      tt.ToolName,
					Arguments: tt.Arguments,
				},
			}

			if err := sendRequest(stdin, req); err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}

			resp, err := readResponse(stdout)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			if tt.ExpectError {
				// Check if error is in JSON-RPC error field
				if resp.Error != nil {
					// Error found in JSON-RPC error field
				} else if result, ok := resp.Result.(map[string]interface{}); ok {
					if isError, ok := result["isError"].(bool); ok && isError {
						// Error found in result with isError flag
					} else {
						t.Error("Expected error but got none")
					}
				} else {
					t.Error("Expected error but got none")
				}
			} else {
				if resp.Error != nil {
					t.Errorf("Unexpected error: %v", resp.Error)
				}
			}

			if tt.Validate != nil {
				tt.Validate(t, resp)
			}
		})
	}
}

// Integration test that runs all tools in sequence
func TestMCPIntegration(t *testing.T) {
	cmd, stdin, stdout := setupTestServer(t)
	defer teardownTestServer(cmd)

	// Initialize server
	initReq := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	if err := sendRequest(stdin, initReq); err != nil {
		t.Fatal("Failed to initialize server:", err)
	}

	// Read init response
	_, err := readResponse(stdout)
	if err != nil {
		t.Fatal("Failed to read init response:", err)
	}

	// Run all tool tests in sequence
	allTests := []ToolTestCase{}
	allTests = append(allTests, getGoalsToolTests()...)
	allTests = append(allTests, getADRToolTests()...)
	allTests = append(allTests, getCIToolTests()...)
	allTests = append(allTests, getMarkdownToolTests()...)
	allTests = append(allTests, getTemplateToolTests()...)

	for i, tt := range allTests {
		t.Run(tt.Name, func(t *testing.T) {
			req := MCPRequest{
				JSONRPC: "2.0",
				ID:      i + 100, // Ensure unique ID
				Method:  "tools/call",
				Params: ToolCallParams{
					Name:      tt.ToolName,
					Arguments: tt.Arguments,
				},
			}

			if err := sendRequest(stdin, req); err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}

			resp, err := readResponse(stdout)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			if tt.ExpectError {
				// Check if error is in JSON-RPC error field
				if resp.Error != nil {
					// Error found in JSON-RPC error field
				} else if result, ok := resp.Result.(map[string]interface{}); ok {
					if isError, ok := result["isError"].(bool); ok && isError {
						// Error found in result with isError flag
					} else {
						t.Error("Expected error but got none")
					}
				} else {
					t.Error("Expected error but got none")
				}
			} else {
				if resp.Error != nil {
					t.Errorf("Unexpected error: %v", resp.Error)
				}
			}

			if tt.Validate != nil {
				tt.Validate(t, resp)
			}
		})
	}
}
