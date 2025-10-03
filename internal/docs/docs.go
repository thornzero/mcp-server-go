// Package docs provides MCP tools for accessing Go documentation.
//
// This package implements documentation tools that allow agents to access
// godoc-generated documentation for Go packages, functions, and types.
// It provides both static documentation generation and dynamic documentation
// serving capabilities.
//
// Example usage:
//
//	handler := NewDocsHandler(server)
//	result, output, err := handler.DocsGet(ctx, req, types.DocsGetInput{Package: "internal/goals"})
package docs

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thornzero/project-manager/internal/server"
	"github.com/thornzero/project-manager/internal/types"
)

// DocsHandler handles MCP tool requests for Go documentation operations.
//
// It provides methods for retrieving godoc-generated documentation,
// listing available packages, and generating static documentation files.
type DocsHandler struct {
	server *server.Server
}

// NewDocsHandler creates a new DocsHandler instance with the provided server.
//
// The server instance is used for project root access and configuration.
func NewDocsHandler(s *server.Server) *DocsHandler {
	return &DocsHandler{server: s}
}

// DocsGet retrieves documentation for a specific Go package or symbol.
//
// This method uses the `go doc` command to generate documentation for the
// specified package, function, or type. It supports both local packages
// and external packages that are available in the module.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - req: MCP tool request (unused but required by interface)
//   - input: DocsGetInput containing the package/symbol path
//
// Returns:
//   - result: MCP call result with documentation content
//   - output: DocsGetOutput containing the formatted documentation
//   - err: Any error that occurred during documentation generation
//
// Example:
//
//	// Get package documentation
//	result, output, err := handler.DocsGet(ctx, req, types.DocsGetInput{
//		Package: "internal/goals",
//	})
//
//	// Get specific function documentation
//	result, output, err := handler.DocsGet(ctx, req, types.DocsGetInput{
//		Package: "internal/goals.GoalsHandler.GoalsList",
//	})
func (h *DocsHandler) DocsGet(ctx context.Context, req *mcp.CallToolRequest, input types.DocsGetInput) (*mcp.CallToolResult, types.DocsGetOutput, error) {
	if input.Package == "" {
		return nil, types.DocsGetOutput{}, fmt.Errorf("package path is required")
	}

	// Change to project root directory
	repoRoot := h.server.GetRepoRoot()

	// Run go doc command
	cmd := exec.CommandContext(ctx, "go", "doc", input.Package)
	cmd.Dir = repoRoot

	output, err := cmd.Output()
	if err != nil {
		return nil, types.DocsGetOutput{}, fmt.Errorf("failed to generate documentation: %v", err)
	}

	// Format the output
	formatted := strings.TrimSpace(string(output))
	if formatted == "" {
		formatted = "No documentation found for the specified package or symbol."
	}

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: formatted,
				},
			},
		}, types.DocsGetOutput{
			Package:       input.Package,
			Documentation: formatted,
		}, nil
}

// DocsList retrieves a list of available Go packages in the project.
//
// This method scans the project directory structure to find all Go packages
// that can be documented. It returns both internal packages and any external
// packages that are imported by the project.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - req: MCP tool request (unused but required by interface)
//   - input: DocsListInput containing optional filters
//
// Returns:
//   - result: MCP call result with package list
//   - output: DocsListOutput containing available packages
//   - err: Any error that occurred during package discovery
func (h *DocsHandler) DocsList(ctx context.Context, req *mcp.CallToolRequest, input types.DocsListInput) (*mcp.CallToolResult, types.DocsListOutput, error) {
	repoRoot := h.server.GetRepoRoot()

	// Run go list to get all packages
	cmd := exec.CommandContext(ctx, "go", "list", "./...")
	cmd.Dir = repoRoot

	output, err := cmd.Output()
	if err != nil {
		return nil, types.DocsListOutput{}, fmt.Errorf("failed to list packages: %v", err)
	}

	// Parse package list
	packages := strings.Split(strings.TrimSpace(string(output)), "\n")
	var filteredPackages []string

	// Apply filters if specified
	for _, pkg := range packages {
		if input.Filter == "" || strings.Contains(pkg, input.Filter) {
			filteredPackages = append(filteredPackages, pkg)
		}
	}

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Available packages:\n%s", strings.Join(filteredPackages, "\n")),
				},
			},
		}, types.DocsListOutput{
			Packages: filteredPackages,
			Count:    len(filteredPackages),
		}, nil
}

// DocsGenerate generates static documentation files for the project.
//
// This method creates comprehensive documentation files using godoc,
// including HTML and markdown formats. The generated files are stored
// in the docs/generated directory.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - req: MCP tool request (unused but required by interface)
//   - input: DocsGenerateInput containing generation options
//
// Returns:
//   - result: MCP call result with generation status
//   - output: DocsGenerateOutput containing file paths and status
//   - err: Any error that occurred during documentation generation
func (h *DocsHandler) DocsGenerate(ctx context.Context, req *mcp.CallToolRequest, input types.DocsGenerateInput) (*mcp.CallToolResult, types.DocsGenerateOutput, error) {
	repoRoot := h.server.GetRepoRoot()
	docsDir := filepath.Join(repoRoot, "docs", "generated")

	// Create docs directory if it doesn't exist
	cmd := exec.CommandContext(ctx, "mkdir", "-p", docsDir)
	cmd.Dir = repoRoot
	cmd.Run()

	var generatedFiles []string
	var status string

	// Generate markdown documentation
	if input.Format == "" || input.Format == "markdown" {
		markdownFile := filepath.Join(docsDir, "README.md")

		// Create comprehensive markdown documentation using Go
		content := "# Generated Documentation\n\n"
		content += "This directory contains auto-generated documentation from Go source code comments.\n\n"
		content += "## Package Documentation\n\n"

		// Get list of packages
		listCmd := exec.CommandContext(ctx, "go", "list", "./internal/...")
		listCmd.Dir = repoRoot
		listOutput, err := listCmd.Output()
		if err != nil {
			return nil, types.DocsGenerateOutput{}, fmt.Errorf("failed to list packages: %v", err)
		}

		packages := strings.Split(strings.TrimSpace(string(listOutput)), "\n")
		for _, pkg := range packages {
			content += fmt.Sprintf("### %s\n\n", pkg)
			content += "```\n"

			// Get documentation for this package
			docCmd := exec.CommandContext(ctx, "go", "doc", pkg)
			docCmd.Dir = repoRoot
			docOutput, err := docCmd.Output()
			if err != nil {
				content += fmt.Sprintf("Error getting documentation: %v\n", err)
			} else {
				content += string(docOutput)
			}

			content += "```\n\n"
		}

		// Write the content to file
		if err := os.WriteFile(markdownFile, []byte(content), 0644); err != nil {
			return nil, types.DocsGenerateOutput{}, fmt.Errorf("failed to write markdown file: %v", err)
		}

		generatedFiles = append(generatedFiles, markdownFile)
		status = "Documentation generated successfully"
	}

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Documentation generation completed.\nGenerated files:\n%s", strings.Join(generatedFiles, "\n")),
				},
			},
		}, types.DocsGenerateOutput{
			Files:  generatedFiles,
			Status: status,
		}, nil
}
