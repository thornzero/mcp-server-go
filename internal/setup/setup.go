// /home/thornzero/Repositories/mcp-server-go/internal/setup/setup.go
package setup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thornzero/mcp-server-go/internal/markdown"
	"github.com/thornzero/mcp-server-go/internal/server"
	"github.com/thornzero/mcp-server-go/internal/types"
)

type SetupHandler struct {
	server *server.Server
}

func NewSetupHandler(s *server.Server) *SetupHandler {
	return &SetupHandler{server: s}
}

func (h *SetupHandler) SetupMCPTools(ctx context.Context, req *mcp.CallToolRequest, input types.SetupMCPToolsInput) (*mcp.CallToolResult, types.SetupMCPToolsOutput, error) {
	// Validate input
	if input.ProjectPath == "" {
		return nil, types.SetupMCPToolsOutput{}, fmt.Errorf("project_path is required")
	}

	// Resolve project path
	projectPath := input.ProjectPath
	if !filepath.IsAbs(projectPath) {
		// If relative path, resolve relative to current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return nil, types.SetupMCPToolsOutput{}, fmt.Errorf("failed to get current working directory: %v", err)
		}
		projectPath = filepath.Join(cwd, projectPath)
	}

	// Check if project path exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return nil, types.SetupMCPToolsOutput{}, fmt.Errorf("project path does not exist: %s", projectPath)
	}

	// Create .cursor/rules directory if it doesn't exist
	rulesDir := filepath.Join(projectPath, ".cursor", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		return nil, types.SetupMCPToolsOutput{}, fmt.Errorf("failed to create rules directory: %v", err)
	}

	// Get MCP server path (where this server is running from)
	mcpServerPath := h.server.GetRepoRoot()

	// Generate MCP tools rule using markdown builder
	ruleBuilder := markdown.MCPToolsRuleBuilder(mcpServerPath)
	rulePath := filepath.Join(rulesDir, "mcp-tools.mdc")
	if err := ruleBuilder.WriteToFile(rulePath); err != nil {
		return nil, types.SetupMCPToolsOutput{}, fmt.Errorf("failed to write MCP tools rule: %v", err)
	}

	// Generate MCP usage guide using markdown builder
	usageBuilder := markdown.MCPUsageGuideBuilder()
	usagePath := filepath.Join(rulesDir, "mcp-tools-usage.mdc")
	if err := usageBuilder.WriteToFile(usagePath); err != nil {
		return nil, types.SetupMCPToolsOutput{}, fmt.Errorf("failed to write MCP usage guide: %v", err)
	}

	// Generate MCP troubleshooting guide using markdown builder
	troubleshootingBuilder := markdown.MCPTroubleshootingGuideBuilder(mcpServerPath)
	troubleshootingPath := filepath.Join(rulesDir, "mcp-tools-troubleshooting.mdc")
	if err := troubleshootingBuilder.WriteToFile(troubleshootingPath); err != nil {
		return nil, types.SetupMCPToolsOutput{}, fmt.Errorf("failed to write MCP troubleshooting guide: %v", err)
	}

	// Return success result
	return nil, types.SetupMCPToolsOutput{
		Success:     true,
		ProjectPath: projectPath,
		RulesDir:    rulesDir,
		FilesCreated: []string{
			"mcp-tools.mdc",
			"mcp-tools-usage.mdc",
			"mcp-tools-troubleshooting.mdc",
		},
		Message: "MCP tools setup complete! Restart Cursor to load the new rules.",
	}, nil
}
