package markdown

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/thornzero/mcp-server-go/internal/server"
	"github.com/thornzero/mcp-server-go/internal/types"
)

type MarkdownHandler struct {
	server *server.Server
}

func NewMarkdownHandler(s *server.Server) *MarkdownHandler {
	return &MarkdownHandler{server: s}
}

func (h *MarkdownHandler) MarkdownLint(ctx context.Context, req *mcp.CallToolRequest, input types.MarkdownLintInput) (*mcp.CallToolResult, types.MarkdownLintOutput, error) {
	// Determine the path to lint
	targetPath := h.server.GetRepoRoot()
	if input.Path != nil && *input.Path != "" {
		targetPath = filepath.Join(h.server.GetRepoRoot(), *input.Path)
	}

	// Determine config file path
	configPath := filepath.Join(h.server.GetRepoRoot(), ".markdownlint.json")
	if input.Config != nil && *input.Config != "" {
		configPath = filepath.Join(h.server.GetRepoRoot(), *input.Config)
	}

	// Check if markdownlint is available
	cmd := exec.Command("which", "markdownlint")
	if err := cmd.Run(); err != nil {
		return nil, types.MarkdownLintOutput{}, fmt.Errorf("markdownlint not found. Please install with: npm install -g markdownlint-cli")
	}

	// Build markdownlint command
	args := []string{}

	// Add config file if it exists
	if _, err := os.Stat(configPath); err == nil {
		args = append(args, "--config", configPath)
	}

	// Add fix flag if requested
	fix := false
	if input.Fix != nil && *input.Fix {
		fix = true
		args = append(args, "--fix")
	}

	// Add target path with globs to include .mdc files
	args = append(args, targetPath, "**/*.mdc")

	// Run markdownlint
	cmd = exec.Command("markdownlint", args...)
	output, err := cmd.CombinedOutput()

	var issues []types.LintIssue
	if err != nil {
		// Parse markdownlint output for issues
		lines := strings.Split(string(output), "\n")
		// Updated regex to handle the actual markdownlint output format
		// Format: file:line rule message [Context: "..."]
		re := regexp.MustCompile(`^(.+?):(\d+)\s+(.+?)\s+(.+?)\s+\[Context:.*\]$`)

		for _, line := range lines {
			if matches := re.FindStringSubmatch(line); matches != nil {
				file := strings.TrimPrefix(matches[1], h.server.GetRepoRoot()+"/")
				lineNum, _ := strconv.Atoi(matches[2])
				colNum := 0 // No column info in this format
				rule := matches[3]
				message := matches[4]

				issues = append(issues, types.LintIssue{
					File:    file,
					Line:    lineNum,
					Column:  colNum,
					Rule:    rule,
					Message: message,
				})
			}
		}
	}

	// Apply custom auto-fixes for issues that markdownlint can't fix
	if fix {
		customFixes := h.applyCustomFixes(targetPath, issues)
		if customFixes > 0 {
			// Re-run markdownlint to get updated issues after custom fixes
			cmd = exec.Command("markdownlint", args...)
			output, err = cmd.CombinedOutput()

			// Re-parse issues after custom fixes
			issues = []types.LintIssue{}
			if err != nil {
				lines := strings.Split(string(output), "\n")
				re := regexp.MustCompile(`^(.+?):(\d+)\s+(.+?)\s+(.+?)\s+\[Context:.*\]$`)

				for _, line := range lines {
					if matches := re.FindStringSubmatch(line); matches != nil {
						file := strings.TrimPrefix(matches[1], h.server.GetRepoRoot()+"/")
						lineNum, _ := strconv.Atoi(matches[2])
						colNum := 0 // No column info in this format
						rule := matches[3]
						message := matches[4]

						issues = append(issues, types.LintIssue{
							File:    file,
							Line:    lineNum,
							Column:  colNum,
							Rule:    rule,
							Message: message,
						})
					}
				}
			}
		}
	}

	// Ensure we always return a non-nil slice
	if issues == nil {
		issues = []types.LintIssue{}
	}

	return nil, types.MarkdownLintOutput{
		Issues: issues,
		Fixed:  fix,
		Path:   targetPath,
	}, nil
}

// applyCustomFixes applies custom auto-fixes for issues that markdownlint can't fix automatically
func (h *MarkdownHandler) applyCustomFixes(targetPath string, issues []types.LintIssue) int {
	fixesApplied := 0

	fileInfo, err := os.Stat(targetPath)
	if os.IsNotExist(err) {
		return fixesApplied
	}

	// Group issues by file
	fileIssues := make(map[string][]types.LintIssue)
	for _, issue := range issues {
		fileIssues[issue.File] = append(fileIssues[issue.File], issue)
	}

	// Determine which files to process
	var filesToProcess []string

	if fileInfo.IsDir() {
		// If targetPath is a directory, process all files that have issues
		for file := range fileIssues {
			filePath := filepath.Join(h.server.GetRepoRoot(), file)
			if _, err := os.Stat(filePath); !os.IsNotExist(err) {
				filesToProcess = append(filesToProcess, filePath)
			}
		}
	} else {
		// If targetPath is a single file, only process that file
		// Find issues for this specific file
		relativePath, err := filepath.Rel(h.server.GetRepoRoot(), targetPath)
		if err != nil {
			return fixesApplied
		}

		// Check if there are issues for this file
		if fileIssueList, exists := fileIssues[relativePath]; exists {
			filesToProcess = append(filesToProcess, targetPath)
			// Update the map to only include issues for this file
			fileIssues = map[string][]types.LintIssue{
				relativePath: fileIssueList,
			}
		}
	}

	// Process each file
	for _, filePath := range filesToProcess {
		// Read file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		lines := strings.Split(string(content), "\n")
		modified := false

		// Get the relative path for this file to match with issues
		relativePath, err := filepath.Rel(h.server.GetRepoRoot(), filePath)
		if err != nil {
			continue
		}

		// Apply fixes for each issue in this file
		if fileIssueList, exists := fileIssues[relativePath]; exists {
			for _, issue := range fileIssueList {
				switch issue.Rule {
				case "MD013/line-length":
					// Auto-fix line length by breaking long lines
					if issue.Line <= len(lines) && issue.Line > 0 {
						line := lines[issue.Line-1]
						if len(line) > 80 {
							// Try to break the line at a good spot (before 80 chars)
							newLines := h.breakLongLine(line, 80)
							if len(newLines) > 1 {
								// Replace the long line with broken lines
								newContent := make([]string, 0, len(lines)+len(newLines)-1)
								newContent = append(newContent, lines[:issue.Line-1]...)
								newContent = append(newContent, newLines...)
								newContent = append(newContent, lines[issue.Line:]...)
								lines = newContent
								modified = true
								fixesApplied++
							}
						}
					}
				case "MD040/fenced-code-language":
					// Auto-fix missing language specification in code blocks
					if issue.Line <= len(lines) && issue.Line > 0 {
						line := lines[issue.Line-1]
						if strings.HasPrefix(line, "```") && len(strings.TrimSpace(line)) == 3 {
							// Try to detect language based on context
							language := h.detectCodeBlockLanguage(lines, issue.Line-1)
							if language != "" {
								lines[issue.Line-1] = "```" + language
								modified = true
								fixesApplied++
							}
						}
					}
				}
			}
		}

		// Write back the modified content
		if modified {
			newContent := strings.Join(lines, "\n")
			if !strings.HasSuffix(newContent, "\n") {
				newContent += "\n"
			}
			os.WriteFile(filePath, []byte(newContent), 0644)
		}
	}

	return fixesApplied
}

// breakLongLine breaks a long line at appropriate points (spaces, commas, etc.)
func (h *MarkdownHandler) breakLongLine(line string, maxLength int) []string {
	if len(line) <= maxLength {
		return []string{line}
	}

	// Try to break at good spots
	breakPoints := []string{" ", ", ", " - ", " * ", "  "}

	for _, breakPoint := range breakPoints {
		pos := strings.LastIndex(line[:maxLength], breakPoint)
		if pos > maxLength/2 { // Don't break too early
			firstLine := strings.TrimSpace(line[:pos])
			secondLine := strings.TrimSpace(line[pos+len(breakPoint):])

			// Add appropriate indentation for continuation
			if strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "   ") || strings.HasPrefix(line, "    ") {
				indent := strings.TrimRight(line, strings.TrimLeft(line, " "))
				secondLine = indent + secondLine
			}

			return []string{firstLine, secondLine}
		}
	}

	// If no good break point found, break at maxLength
	firstLine := line[:maxLength-3] + "..."
	secondLine := "... " + line[maxLength-3:]
	return []string{firstLine, secondLine}
}

// detectCodeBlockLanguage tries to detect the language for a code block based on context
func (h *MarkdownHandler) detectCodeBlockLanguage(lines []string, codeBlockLine int) string {
	// Look at the content of the code block to guess the language

	// Check a few lines after the opening ```
	for i := codeBlockLine + 1; i < len(lines) && i < codeBlockLine+5; i++ {
		line := lines[i]
		if strings.HasPrefix(line, "```") {
			break // End of code block
		}

		// Simple heuristics for common languages
		if strings.Contains(line, "package ") && strings.Contains(line, "import ") {
			return "go"
		}
		if strings.Contains(line, "func ") && strings.Contains(line, "(") {
			return "go"
		}
		if strings.Contains(line, "function ") && strings.Contains(line, "(") {
			return "javascript"
		}
		if strings.Contains(line, "const ") && strings.Contains(line, "=") {
			return "javascript"
		}
		if strings.Contains(line, "def ") && strings.Contains(line, "(") {
			return "python"
		}
		if strings.Contains(line, "import ") && strings.Contains(line, "from ") {
			return "python"
		}
		if strings.Contains(line, "SELECT ") && strings.Contains(line, "FROM ") {
			return "sql"
		}
		if strings.Contains(line, "class ") && strings.Contains(line, "{") {
			return "java"
		}
		if strings.Contains(line, "public ") && strings.Contains(line, "static ") {
			return "java"
		}
		if strings.Contains(line, "#include") || strings.Contains(line, "int main") {
			return "cpp"
		}
		if strings.Contains(line, "docker") || strings.Contains(line, "FROM ") {
			return "dockerfile"
		}
		if strings.Contains(line, "yaml") || strings.Contains(line, "---") {
			return "yaml"
		}
		if strings.Contains(line, "{") && strings.Contains(line, "}") {
			return "json"
		}
		if strings.Contains(line, "├") || strings.Contains(line, "└") {
			return "tree"
		}
	}

	// Default fallbacks based on file extension or common patterns
	if strings.Contains(strings.Join(lines[max(0, codeBlockLine-3):codeBlockLine+3], "\n"), "bash") ||
		strings.Contains(strings.Join(lines[max(0, codeBlockLine-3):codeBlockLine+3], "\n"), "$ ") {
		return "bash"
	}

	// If no specific language detected, use "text" as a safe default
	return "text"
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
