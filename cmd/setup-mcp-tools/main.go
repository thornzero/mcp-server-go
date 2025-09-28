// /home/thornzero/Repositories/mcp-server-go/cmd/setup-mcp-tools/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/thornzero/mcp-server-go/internal/markdown"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: setup-mcp-tools <project-root>")
		fmt.Println("Example: setup-mcp-tools /path/to/your/project")
		os.Exit(1)
	}

	projectRoot := os.Args[1]
	if err := setupMCPTools(projectRoot); err != nil {
		fmt.Printf("Error setting up MCP tools: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ MCP tools setup complete!")
	fmt.Println("📋 Next steps:")
	fmt.Println("1. Restart Cursor to load the new rules")
	fmt.Println("2. Test with: mcp_mcp-server-go_goals_list()")
	fmt.Println("3. Add initial data: mcp_mcp-server-go_goals_add({title: 'Project Goal'})")
}

func setupMCPTools(projectRoot string) error {
	// Create .cursor/rules directory if it doesn't exist
	rulesDir := filepath.Join(projectRoot, ".cursor", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		return fmt.Errorf("failed to create rules directory: %v", err)
	}

	// Generate MCP tools rule using markdown builder
	mcpServerPath := "/home/thornzero/Repositories/mcp-server-go"
	ruleBuilder := markdown.MCPToolsRuleBuilder(mcpServerPath)
	rulePath := filepath.Join(rulesDir, "mcp-tools.mdc")
	if err := ruleBuilder.WriteToFile(rulePath); err != nil {
		return fmt.Errorf("failed to write MCP tools rule: %v", err)
	}

	// Generate MCP usage guide using markdown builder
	usageBuilder := markdown.MCPUsageGuideBuilder()
	usagePath := filepath.Join(rulesDir, "mcp-tools-usage.mdc")
	if err := usageBuilder.WriteToFile(usagePath); err != nil {
		return fmt.Errorf("failed to write MCP usage guide: %v", err)
	}

	// Generate MCP troubleshooting guide using markdown builder
	troubleshootingBuilder := markdown.MCPTroubleshootingGuideBuilder(mcpServerPath)
	troubleshootingPath := filepath.Join(rulesDir, "mcp-tools-troubleshooting.mdc")
	if err := troubleshootingBuilder.WriteToFile(troubleshootingPath); err != nil {
		return fmt.Errorf("failed to write MCP troubleshooting guide: %v", err)
	}

	return nil
}
