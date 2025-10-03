package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thornzero/project-manager/internal/adrs"
	"github.com/thornzero/project-manager/internal/ci"
	"github.com/thornzero/project-manager/internal/cursorrules"
	"github.com/thornzero/project-manager/internal/docs"
	"github.com/thornzero/project-manager/internal/goals"
	"github.com/thornzero/project-manager/internal/logparser"
	"github.com/thornzero/project-manager/internal/markdown"
	"github.com/thornzero/project-manager/internal/preferredtools"
	"github.com/thornzero/project-manager/internal/search"
	"github.com/thornzero/project-manager/internal/server"
	"github.com/thornzero/project-manager/internal/setup"
	"github.com/thornzero/project-manager/internal/state"
	"github.com/thornzero/project-manager/internal/templates"
)

func main() {
	// Get the directory where the executable is located
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	// Determine project root - try multiple strategies
	repoRoot := findProjectRoot(execPath)

	// Initialize server
	srv, err := server.NewServer(repoRoot)
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Close()

	// Create handlers
	goalsHandler := goals.NewGoalsHandler(srv)
	adrsHandler := adrs.NewADRsHandler(srv)
	ciHandler := ci.NewCIHandler(srv)
	searchHandler := search.NewSearchHandler(srv)
	stateHandler := state.NewStateHandler(srv)
	markdownHandler := markdown.NewMarkdownHandler(srv)
	templatesHandler := templates.NewTemplatesHandler(srv)
	preferredToolsHandler := preferredtools.NewPreferredToolsHandler(srv)
	cursorRulesHandler := cursorrules.NewCursorRulesHandler(srv)
	setupHandler := setup.NewSetupHandler(srv)
	logParserHandler := logparser.NewLogParserHandler(srv)
	docsHandler := docs.NewDocsHandler(srv)

	// Create MCP server
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "project-manager",
		Version: "1.0.0",
	}, nil)

	// Add tools
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "goals_list",
		Description: "List active goals from the project",
	}, goalsHandler.GoalsList)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "goals_add",
		Description: "Add a new goal to the project",
	}, goalsHandler.GoalsAdd)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "goals_update",
		Description: "Update an existing goal",
	}, goalsHandler.GoalsUpdate)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "adrs_list",
		Description: "List Architecture Decision Records (ADRs)",
	}, adrsHandler.ADRsList)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "adrs_get",
		Description: "Get the content of a specific ADR",
	}, adrsHandler.ADRsGet)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "ci_run_tests",
		Description: "Run tests for the project",
	}, ciHandler.CIRunTests)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "ci_last_failure",
		Description: "Get information about the last test failure",
	}, ciHandler.CILastFailure)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "repo_search",
		Description: "Search the repository for text patterns",
	}, searchHandler.RepoSearch)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "state_log_change",
		Description: "Log a change to the project changelog",
	}, stateHandler.StateLogChange)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "changelog_generate",
		Description: "Generate/update a proper changelog file in the root directory",
	}, stateHandler.ChangelogGenerate)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "markdown_lint",
		Description: "Lint markdown files for formatting issues",
	}, markdownHandler.MarkdownLint)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "template_list",
		Description: "List available markdown templates",
	}, templatesHandler.TemplateList)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "template_register",
		Description: "Register a new markdown template",
	}, templatesHandler.TemplateRegister)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "template_get",
		Description: "Get template details by ID",
	}, templatesHandler.TemplateGet)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "template_update",
		Description: "Update an existing markdown template",
	}, templatesHandler.TemplateUpdate)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "template_delete",
		Description: "Delete a markdown template",
	}, templatesHandler.TemplateDelete)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "template_apply",
		Description: "Apply a template to generate markdown content",
	}, templatesHandler.TemplateApply)

	// Preferred Tools tools
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "preferred_tools_list",
		Description: "List preferred tools for specific categories and languages",
	}, preferredToolsHandler.PreferredToolsList)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "preferred_tools_add",
		Description: "Add a new preferred tool",
	}, preferredToolsHandler.PreferredToolsAdd)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "preferred_tools_update",
		Description: "Update an existing preferred tool",
	}, preferredToolsHandler.PreferredToolsUpdate)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "preferred_tools_delete",
		Description: "Delete a preferred tool",
	}, preferredToolsHandler.PreferredToolsDelete)

	// Cursor Rules tools
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "cursor_rules_list",
		Description: "List Cursor rules with optional filtering",
	}, cursorRulesHandler.CursorRulesList)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "cursor_rules_add",
		Description: "Add a new Cursor rule",
	}, cursorRulesHandler.CursorRulesAdd)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "cursor_rules_update",
		Description: "Update an existing Cursor rule",
	}, cursorRulesHandler.CursorRulesUpdate)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "cursor_rules_delete",
		Description: "Delete a Cursor rule",
	}, cursorRulesHandler.CursorRulesDelete)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "cursor_rules_suggest",
		Description: "Suggest community Cursor rules based on criteria",
	}, cursorRulesHandler.CursorRulesSuggest)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "cursor_rules_install",
		Description: "Install a Cursor rule from community repository",
	}, cursorRulesHandler.CursorRulesInstall)

	// Setup tools
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "setup_project_manager",
		Description: "Set up Project Manager tools for a project by creating cursor rules",
	}, setupHandler.SetupProjectManager)

	// Log parsing tools
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "log_parse",
		Description: "Parse and analyze Cursor/VS Code log files with AI-optimized output",
	}, logParserHandler.ParseLog)

	// Documentation tools
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "docs_get",
		Description: "Get documentation for a specific Go package or symbol using godoc",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Get Go Documentation",
			ReadOnlyHint:    true,
			OpenWorldHint:   &[]bool{false}[0],
			DestructiveHint: &[]bool{false}[0],
			IdempotentHint:  true,
		},
	}, docsHandler.DocsGet)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "docs_list",
		Description: "List available Go packages in the project for documentation",
		Annotations: &mcp.ToolAnnotations{
			Title:           "List Go Packages",
			ReadOnlyHint:    true,
			OpenWorldHint:   &[]bool{false}[0],
			DestructiveHint: &[]bool{false}[0],
			IdempotentHint:  true,
		},
	}, docsHandler.DocsList)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "docs_generate",
		Description: "Generate static documentation files for the project using godoc",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Generate Documentation",
			ReadOnlyHint:    false,
			OpenWorldHint:   &[]bool{false}[0],
			DestructiveHint: &[]bool{true}[0],
			IdempotentHint:  true,
		},
	}, docsHandler.DocsGenerate)

	// Run the server over stdin/stdout
	if err := mcpServer.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

// findProjectRoot attempts to find the project root directory using multiple strategies
func findProjectRoot(execPath string) string {
	// Strategy 1: If executable is in build directory, go up one level
	execDir := filepath.Dir(execPath)
	if filepath.Base(execDir) == "build" {
		candidate := filepath.Dir(execDir)
		if isProjectRoot(candidate) {
			return candidate
		}
	}

	// Strategy 2: Look for go.mod file starting from executable directory and going up
	current := execDir
	for {
		if isProjectRoot(current) {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			break // Reached filesystem root
		}
		current = parent
	}

	// Strategy 3: Fallback to executable directory
	return execDir
}

// isProjectRoot checks if a directory is the project root by looking for go.mod
func isProjectRoot(dir string) bool {
	goModPath := filepath.Join(dir, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		// Also check if it contains our module name
		content, err := os.ReadFile(goModPath)
		if err == nil && strings.Contains(string(content), "github.com/thornzero/project-manager") {
			return true
		}
	}
	return false
}
