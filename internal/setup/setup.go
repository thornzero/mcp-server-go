// /home/thornzero/Repositories/mcp-server-go/internal/setup/setup.go
package setup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

// detectProjectType analyzes the project to determine what type of rules to create
func (h *SetupHandler) detectProjectType(projectPath string) string {
	// Check for Go project
	if _, err := os.Stat(filepath.Join(projectPath, "go.mod")); err == nil {
		return "go"
	}

	// Check for Node.js project
	if _, err := os.Stat(filepath.Join(projectPath, "package.json")); err == nil {
		return "nodejs"
	}

	// Check for Python project
	if _, err := os.Stat(filepath.Join(projectPath, "requirements.txt")); err == nil {
		return "python"
	}
	if _, err := os.Stat(filepath.Join(projectPath, "pyproject.toml")); err == nil {
		return "python"
	}

	// Check for Rust project
	if _, err := os.Stat(filepath.Join(projectPath, "Cargo.toml")); err == nil {
		return "rust"
	}

	// Default to generic
	return "generic"
}

// shouldCreateGenericRules determines if generic rules should be created
func (h *SetupHandler) shouldCreateGenericRules(projectPath string) bool {
	// Don't create generic rules if there are already existing rules
	rulesDir := filepath.Join(projectPath, ".cursor", "rules")
	if _, err := os.Stat(rulesDir); err == nil {
		files, err := os.ReadDir(rulesDir)
		if err == nil && len(files) > 0 {
			// Check if there are already numbered rules (00-, 01-, etc.)
			for _, file := range files {
				if strings.HasPrefix(file.Name(), "0") && strings.HasSuffix(file.Name(), ".md") {
					return false // Don't create generic rules if numbered rules exist
				}
			}
		}
	}

	// Don't create generic rules for MCP server projects
	if strings.Contains(projectPath, "mcp-server") {
		return false
	}

	return true
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

	// Detect project type
	projectType := h.detectProjectType(projectPath)

	// Determine what rules to create
	shouldCreateGeneric := h.shouldCreateGenericRules(projectPath)

	var filesCreated []string

	// Always create MCP-specific rules
	mcpServerPath := h.server.GetRepoRoot()

	// Generate MCP tools rule using markdown builder
	ruleBuilder := markdown.MCPToolsRuleBuilder(mcpServerPath)
	rulePath := filepath.Join(rulesDir, "mcp-tools.mdc")
	if err := ruleBuilder.WriteToFile(rulePath); err != nil {
		return nil, types.SetupMCPToolsOutput{}, fmt.Errorf("failed to write MCP tools rule: %v", err)
	}
	filesCreated = append(filesCreated, "mcp-tools.mdc")

	// Generate MCP usage guide using markdown builder
	usageBuilder := markdown.MCPUsageGuideBuilder()
	usagePath := filepath.Join(rulesDir, "mcp-tools-usage.mdc")
	if err := usageBuilder.WriteToFile(usagePath); err != nil {
		return nil, types.SetupMCPToolsOutput{}, fmt.Errorf("failed to write MCP usage guide: %v", err)
	}
	filesCreated = append(filesCreated, "mcp-tools-usage.mdc")

	// Generate MCP troubleshooting guide using markdown builder
	troubleshootingBuilder := markdown.MCPTroubleshootingGuideBuilder(mcpServerPath)
	troubleshootingPath := filepath.Join(rulesDir, "mcp-tools-troubleshooting.mdc")
	if err := troubleshootingBuilder.WriteToFile(troubleshootingPath); err != nil {
		return nil, types.SetupMCPToolsOutput{}, fmt.Errorf("failed to write MCP troubleshooting guide: %v", err)
	}
	filesCreated = append(filesCreated, "mcp-tools-troubleshooting.mdc")

	// Only create generic rules if appropriate
	if shouldCreateGeneric {
		// Create project-specific generic rules based on project type
		switch projectType {
		case "go":
			h.createGoProjectRules(rulesDir, &filesCreated)
		case "nodejs":
			h.createNodeJSProjectRules(rulesDir, &filesCreated)
		case "python":
			h.createPythonProjectRules(rulesDir, &filesCreated)
		case "rust":
			h.createRustProjectRules(rulesDir, &filesCreated)
		default:
			h.createGenericProjectRules(rulesDir, &filesCreated)
		}
	}

	// Return success result
	return nil, types.SetupMCPToolsOutput{
		Success:      true,
		ProjectPath:  projectPath,
		RulesDir:     rulesDir,
		FilesCreated: filesCreated,
		Message:      fmt.Sprintf("MCP tools setup complete! Created %d files. Restart Cursor to load the new rules.", len(filesCreated)),
	}, nil
}

// createGoProjectRules creates Go-specific project rules
func (h *SetupHandler) createGoProjectRules(rulesDir string, filesCreated *[]string) {
	// For now, just create a basic Go rule
	goRule := `---
title: Go Project Guidelines
description: Basic guidelines for Go projects
globs: ["**/*.go", "go.mod", "go.sum"]
alwaysApply: true
---

# Go Project Guidelines

## Code Style
- Follow Go conventions and idioms
- Use meaningful variable and function names
- Handle errors explicitly
- Use context for cancellation and timeouts

## Testing
- Write tests for new functionality
- Use table-driven tests where appropriate
- Run tests before committing: ` + "`" + `go test ./...` + "`" + `

## Dependencies
- Keep dependencies minimal
- Use ` + "`" + `go mod tidy` + "`" + ` to clean up dependencies
- Pin major versions in go.mod
`

	goRulePath := filepath.Join(rulesDir, "go-guidelines.mdc")
	if err := os.WriteFile(goRulePath, []byte(goRule), 0644); err == nil {
		*filesCreated = append(*filesCreated, "go-guidelines.mdc")
	}
}

// createNodeJSProjectRules creates Node.js-specific project rules
func (h *SetupHandler) createNodeJSProjectRules(rulesDir string, filesCreated *[]string) {
	nodeRule := `---
title: Node.js Project Guidelines
description: Basic guidelines for Node.js projects
globs: ["**/*.js", "**/*.ts", "package.json"]
alwaysApply: true
---

# Node.js Project Guidelines

## Code Style
- Use ESLint and Prettier for code formatting
- Follow consistent naming conventions
- Use TypeScript for better type safety

## Testing
- Write unit tests with Jest or similar
- Use meaningful test descriptions
- Mock external dependencies

## Dependencies
- Keep package.json clean
- Use exact versions for critical dependencies
- Regular security audits with ` + "`" + `npm audit` + "`" + `
`

	nodeRulePath := filepath.Join(rulesDir, "nodejs-guidelines.mdc")
	if err := os.WriteFile(nodeRulePath, []byte(nodeRule), 0644); err == nil {
		*filesCreated = append(*filesCreated, "nodejs-guidelines.mdc")
	}
}

// createPythonProjectRules creates Python-specific project rules
func (h *SetupHandler) createPythonProjectRules(rulesDir string, filesCreated *[]string) {
	pythonRule := `---
title: Python Project Guidelines
description: Basic guidelines for Python projects
globs: ["**/*.py", "requirements.txt", "pyproject.toml"]
alwaysApply: true
---

# Python Project Guidelines

## Code Style
- Follow PEP 8 style guidelines
- Use type hints for better code clarity
- Use meaningful variable and function names

## Testing
- Write tests with pytest or unittest
- Use descriptive test names
- Mock external dependencies

## Dependencies
- Use virtual environments
- Pin dependency versions
- Regular security updates
`

	pythonRulePath := filepath.Join(rulesDir, "python-guidelines.mdc")
	if err := os.WriteFile(pythonRulePath, []byte(pythonRule), 0644); err == nil {
		*filesCreated = append(*filesCreated, "python-guidelines.mdc")
	}
}

// createRustProjectRules creates Rust-specific project rules
func (h *SetupHandler) createRustProjectRules(rulesDir string, filesCreated *[]string) {
	rustRule := `---
title: Rust Project Guidelines
description: Basic guidelines for Rust projects
globs: ["**/*.rs", "Cargo.toml", "Cargo.lock"]
alwaysApply: true
---

# Rust Project Guidelines

## Code Style
- Follow Rust conventions and idioms
- Use meaningful variable and function names
- Handle errors with Result and Option types
- Use clippy for additional linting

## Testing
- Write unit tests in the same file
- Use integration tests in tests/ directory
- Run tests with ` + "`" + `cargo test` + "`" + `

## Dependencies
- Keep Cargo.toml clean
- Use exact versions for critical dependencies
- Regular updates with ` + "`" + `cargo update` + "`" + `
`

	rustRulePath := filepath.Join(rulesDir, "rust-guidelines.mdc")
	if err := os.WriteFile(rustRulePath, []byte(rustRule), 0644); err == nil {
		*filesCreated = append(*filesCreated, "rust-guidelines.mdc")
	}
}

// createGenericProjectRules creates generic project rules
func (h *SetupHandler) createGenericProjectRules(rulesDir string, filesCreated *[]string) {
	genericRule := `---
title: General Project Guidelines
description: Basic guidelines for any project
globs: ["**/*"]
alwaysApply: true
---

# General Project Guidelines

## Code Quality
- Write clean, readable code
- Use meaningful variable and function names
- Add comments for complex logic
- Follow consistent formatting

## Testing
- Write tests for new functionality
- Test edge cases and error conditions
- Keep tests simple and focused

## Documentation
- Keep README up to date
- Document API changes
- Use clear commit messages
`

	genericRulePath := filepath.Join(rulesDir, "general-guidelines.mdc")
	if err := os.WriteFile(genericRulePath, []byte(genericRule), 0644); err == nil {
		*filesCreated = append(*filesCreated, "general-guidelines.mdc")
	}
}
