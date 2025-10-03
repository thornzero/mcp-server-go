// /home/thornzero/Repositories/project-manager/internal/setup/setup.go
package setup

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thornzero/project-manager/internal/server"
	"github.com/thornzero/project-manager/internal/types"
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
	return true
}

func (h *SetupHandler) SetupProjectManager(ctx context.Context, req *mcp.CallToolRequest, input types.SetupProjectManagerInput) (*mcp.CallToolResult, types.SetupProjectManagerOutput, error) {
	// Validate input
	if input.ProjectPath == "" {
		return nil, types.SetupProjectManagerOutput{}, fmt.Errorf("project_path is required")
	}

	// Resolve project path
	projectPath := input.ProjectPath
	if !filepath.IsAbs(projectPath) {
		// If relative path, resolve relative to current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return nil, types.SetupProjectManagerOutput{}, fmt.Errorf("failed to get current working directory: %v", err)
		}
		projectPath = filepath.Join(cwd, projectPath)
	}

	// Check if project path exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return nil, types.SetupProjectManagerOutput{}, fmt.Errorf("project path does not exist: %s", projectPath)
	}

	// Create .cursor/rules directory if it doesn't exist
	rulesDir := filepath.Join(projectPath, ".cursor", "rules")
	if err := os.MkdirAll(rulesDir, 0755); err != nil {
		return nil, types.SetupProjectManagerOutput{}, fmt.Errorf("failed to create rules directory: %v", err)
	}

	// Detect project type
	projectType := h.detectProjectType(projectPath)

	// Determine what rules to create
	shouldCreateGeneric := h.shouldCreateGenericRules(projectPath)

	var filesCreated []string

	// Copy Project Manager rule files from docs/rules/
	mcpServerPath := h.server.GetRepoRoot()
	sourceDir := filepath.Join(mcpServerPath, "docs", "rules")

	// List of rule files to copy
	ruleFiles, err := os.ReadDir(sourceDir)
	if err != nil {
		return nil, types.SetupProjectManagerOutput{}, fmt.Errorf("failed to read rules directory: %v", err)
	}

	for _, filename := range ruleFiles {
		sourcePath := filepath.Join(sourceDir, filename.Name())
		destPath := filepath.Join(rulesDir, filename.Name())

		// Copy file from source to destination
		if err := h.copyFile(sourcePath, destPath, filename.Name()); err != nil {
			return nil, types.SetupProjectManagerOutput{}, fmt.Errorf("failed to copy %s: %v", filename, err)
		}
		filesCreated = append(filesCreated, filename.Name())
	}

	// Only create generic rules if appropriate
	if shouldCreateGeneric {
		// Create project-specific generic rules based on project type
		switch projectType {
		case "go":
			h.copyFile(sourceDir, rulesDir, "go-guidelines.mdc")
		case "nodejs":
			h.copyFile(sourceDir, rulesDir, "nodejs-guidelines.mdc")
		case "python":
			h.copyFile(sourceDir, rulesDir, "python-guidelines.mdc")
		case "rust":
			h.copyFile(sourceDir, rulesDir, "rust-guidelines.mdc")
		default:
			h.copyFile(sourceDir, rulesDir, "general-guidelines.mdc")
		}
	}

	// Return success result
	return nil, types.SetupProjectManagerOutput{
		Success:      true,
		ProjectPath:  projectPath,
		RulesDir:     rulesDir,
		FilesCreated: filesCreated,
		Message:      fmt.Sprintf("Project Manager tools setup complete! Created %d files. Restart Cursor to load the new rules.", len(filesCreated)),
	}, nil
}

// copyFile copies a file from source to destination
func (h *SetupHandler) copyFile(src, dst, filename string) error {
	sourcePath := filepath.Join(src, filename)
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	destPath := filepath.Join(dst, filename)
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	err = destFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %v", err)
	}

	return nil
}
