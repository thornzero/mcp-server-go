// /home/thornzero/Repositories/mcp-server-go/internal/logparser/logparser.go
package logparser

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thornzero/mcp-server-go/internal/types"
)

// LogParserHandler handles log parsing operations
type LogParserHandler struct {
	server interface {
		GetRepoRoot() string
	}
}

// NewLogParserHandler creates a new log parser handler
func NewLogParserHandler(server interface {
	GetRepoRoot() string
}) *LogParserHandler {
	return &LogParserHandler{
		server: server,
	}
}

// ErrorPatternConfig defines error patterns to look for
type ErrorPatternConfig struct {
	Pattern     string
	Description string
	Severity    string
}

var errorPatterns = []ErrorPatternConfig{
	{"ERR_INTERNET_DISCONNECTED", "Internet connectivity failures", "low"},
	{"ENOENT.*no such file", "Missing file errors", "medium"},
	{"EntryNotFound.*FileSystemError", "File system errors", "medium"},
	{"potential listener LEAK", "Memory leak warnings", "high"},
	{"mcp.*error|MCP.*error", "MCP-related errors", "high"},
	{"server.*error|Server.*error", "Server-related errors", "high"},
	{"composer.*Error", "Composer context errors", "medium"},
	{"TypeError.*undefined", "JavaScript type errors", "medium"},
	{"ReferenceError.*not defined", "JavaScript reference errors", "medium"},
	{"SyntaxError", "JavaScript syntax errors", "high"},
	{"Permission denied", "Permission errors", "high"},
	{"EACCES", "Access denied errors", "high"},
	{"ENOSPC.*No space left", "Disk space errors", "critical"},
	{"ECONNREFUSED", "Connection refused errors", "medium"},
	{"ETIMEDOUT", "Connection timeout errors", "medium"},
}

// ParseLog analyzes a log file and returns structured analysis
func (h *LogParserHandler) ParseLog(ctx context.Context, req *mcp.CallToolRequest, input types.LogParseInput) (*mcp.CallToolResult, types.LogParseOutput, error) {
	// Determine log file path
	filePath := input.FilePath
	if filePath == "" {
		// Auto-detect log files
		defaultLogs := []string{
			filepath.Join(h.server.GetRepoRoot(), "docs/logs/vscode-app-1759085846096.log"),
			filepath.Join(h.server.GetRepoRoot(), "docs/logs/vscode-app-1759084599756.log"),
		}

		for _, logPath := range defaultLogs {
			if _, err := os.Stat(logPath); err == nil {
				filePath = logPath
				break
			}
		}

		if filePath == "" {
			return nil, types.LogParseOutput{}, fmt.Errorf("no log file specified and no default log found")
		}
	}

	// Check if log file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, types.LogParseOutput{}, fmt.Errorf("log file not found: %s", filePath)
	}

	// Analyze the log file
	analysis, err := h.analyzeLogFile(filePath)
	if err != nil {
		return nil, types.LogParseOutput{}, fmt.Errorf("error analyzing log file: %w", err)
	}

	return &mcp.CallToolResult{}, *analysis, nil
}

func (h *LogParserHandler) analyzeLogFile(filePath string) (*types.LogParseOutput, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get file statistics
	stats, err := h.getFileStatistics(filePath)
	if err != nil {
		return nil, err
	}

	// Read all lines
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Analyze error patterns
	errorPatterns := h.analyzeErrorPatterns(lines)

	// Count different error types
	errorCounts := h.countErrorTypes(lines)

	// Get critical issues
	criticalIssues := h.getCriticalIssues(lines)

	// Get recent errors
	recentErrors := h.getRecentErrors(lines, 10)

	// Get missing files
	missingFiles := h.getMissingFiles(lines)

	// Generate recommendations
	recommendations := h.generateRecommendations(errorCounts, criticalIssues)

	// Generate context
	context := h.generateContext(errorCounts, criticalIssues)

	return &types.LogParseOutput{
		File:            filepath.Base(filePath),
		AnalysisTime:    time.Now().Format(time.RFC3339),
		Statistics:      *stats,
		ErrorCounts:     *errorCounts,
		CriticalIssues:  *criticalIssues,
		ErrorPatterns:   errorPatterns,
		RecentErrors:    recentErrors,
		MissingFiles:    missingFiles,
		Recommendations: recommendations,
		Context:         *context,
	}, nil
}

func (h *LogParserHandler) getFileStatistics(filePath string) (*types.FileStatistics, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Count lines
	lines := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines++
	}

	// Get file info
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &types.FileStatistics{
		Lines:    lines,
		Size:     h.formatBytes(info.Size()),
		Modified: info.ModTime().Format("2006-01-02"),
	}, nil
}

func (h *LogParserHandler) analyzeErrorPatterns(lines []string) []types.ErrorPattern {
	var patterns []types.ErrorPattern

	for _, config := range errorPatterns {
		regex, err := regexp.Compile("(?i)" + config.Pattern)
		if err != nil {
			continue
		}

		var matches []string
		count := 0

		for _, line := range lines {
			if regex.MatchString(line) {
				count++
				if len(matches) < 5 { // Keep only recent 5
					matches = append(matches, line)
				}
			}
		}

		if count > 0 {
			patterns = append(patterns, types.ErrorPattern{
				Pattern:     config.Pattern,
				Description: config.Description,
				Severity:    config.Severity,
				Count:       count,
				Recent:      matches,
			})
		}
	}

	// Sort by count descending
	sort.Slice(patterns, func(i, j int) bool {
		return patterns[i].Count > patterns[j].Count
	})

	return patterns
}

func (h *LogParserHandler) countErrorTypes(lines []string) *types.ErrorCounts {
	counts := &types.ErrorCounts{}

	for _, line := range lines {
		lowerLine := strings.ToLower(line)

		if strings.Contains(line, "ERR_INTERNET_DISCONNECTED") {
			counts.NetworkIssues++
		}
		if strings.Contains(line, "ENOENT") || strings.Contains(line, "EntryNotFound") {
			counts.FileErrors++
		}
		if strings.Contains(line, "listener LEAK") {
			counts.MemoryLeaks++
		}
		if strings.Contains(lowerLine, "composer") && strings.Contains(lowerLine, "error") {
			counts.ComposerErrors++
		}
		if strings.Contains(lowerLine, "mcp") && strings.Contains(lowerLine, "error") {
			counts.MCPErrors++
		}
		if strings.Contains(line, "TypeError") || strings.Contains(line, "ReferenceError") || strings.Contains(line, "SyntaxError") {
			counts.JavaScriptErrors++
		}
		if strings.Contains(line, "Permission denied") || strings.Contains(line, "EACCES") {
			counts.PermissionErrors++
		}
		if strings.Contains(line, "ECONNREFUSED") || strings.Contains(line, "ETIMEDOUT") {
			counts.ConnectionErrors++
		}
	}

	return counts
}

func (h *LogParserHandler) getCriticalIssues(lines []string) *types.CriticalIssues {
	issues := &types.CriticalIssues{}

	for _, line := range lines {
		if strings.Contains(line, "ENOSPC") && strings.Contains(line, "No space left") {
			issues.DiskSpace++
		}
		if strings.Contains(line, "SyntaxError") {
			issues.SyntaxErrors++
		}
		if strings.Contains(line, "Permission denied") || strings.Contains(line, "EACCES") {
			issues.PermissionDenied++
		}
	}

	return issues
}

func (h *LogParserHandler) getRecentErrors(lines []string, limit int) []string {
	var errors []string
	errorRegex := regexp.MustCompile("(?i)(error|warn|fail)")

	for _, line := range lines {
		if errorRegex.MatchString(line) {
			errors = append(errors, line)
		}
	}

	// Return last N errors
	if len(errors) > limit {
		return errors[len(errors)-limit:]
	}
	return errors
}

func (h *LogParserHandler) getMissingFiles(lines []string) []string {
	var files []string
	fileSet := make(map[string]bool)

	regex := regexp.MustCompile(`stat '([^']+)'`)

	for _, line := range lines {
		if strings.Contains(line, "ENOENT") && strings.Contains(line, "no such file") {
			matches := regex.FindStringSubmatch(line)
			if len(matches) > 1 {
				file := matches[1]
				if !fileSet[file] {
					files = append(files, file)
					fileSet[file] = true
				}
			}
		}
	}

	sort.Strings(files)
	return files
}

func (h *LogParserHandler) generateRecommendations(counts *types.ErrorCounts, critical *types.CriticalIssues) []string {
	var recommendations []string

	if critical.DiskSpace > 0 {
		recommendations = append(recommendations, fmt.Sprintf("🚨 CRITICAL: Free up disk space (%d errors)", critical.DiskSpace))
	}
	if critical.SyntaxErrors > 0 {
		recommendations = append(recommendations, fmt.Sprintf("🚨 HIGH: Fix JavaScript syntax errors (%d errors)", critical.SyntaxErrors))
	}
	if critical.PermissionDenied > 0 {
		recommendations = append(recommendations, fmt.Sprintf("🚨 HIGH: Fix permission issues (%d errors)", critical.PermissionDenied))
	}
	if counts.MCPErrors > 0 {
		recommendations = append(recommendations, fmt.Sprintf("🔧 HIGH: Investigate MCP server issues (%d errors)", counts.MCPErrors))
	}
	if counts.MemoryLeaks > 0 {
		recommendations = append(recommendations, fmt.Sprintf("🧠 MEDIUM: Restart Cursor to clear memory leaks (%d warnings)", counts.MemoryLeaks))
	}
	if counts.FileErrors > 0 {
		recommendations = append(recommendations, fmt.Sprintf("📁 MEDIUM: Clean up missing files (%d errors)", counts.FileErrors))
	}
	if counts.NetworkIssues > 0 {
		recommendations = append(recommendations, fmt.Sprintf("🌐 LOW: Check internet connection (%d errors - usually harmless)", counts.NetworkIssues))
	}

	recommendations = append(recommendations, "Most errors are harmless for local development")

	return recommendations
}

func (h *LogParserHandler) generateContext(counts *types.ErrorCounts, critical *types.CriticalIssues) *types.AnalysisContext {
	totalErrors := counts.NetworkIssues + counts.FileErrors + counts.MemoryLeaks +
		counts.ComposerErrors + counts.MCPErrors + counts.JavaScriptErrors +
		counts.PermissionErrors + counts.ConnectionErrors

	var mostCommonIssue string
	if counts.NetworkIssues >= counts.FileErrors && counts.NetworkIssues >= counts.MemoryLeaks {
		mostCommonIssue = "Network connectivity (usually harmless)"
	} else if counts.FileErrors >= counts.MemoryLeaks {
		mostCommonIssue = "Missing files"
	} else {
		mostCommonIssue = "Memory leaks"
	}

	var severityLevel string
	if critical.DiskSpace > 0 || critical.SyntaxErrors > 0 || critical.PermissionDenied > 0 {
		severityLevel = "HIGH - Immediate action required"
	} else if counts.MCPErrors > 0 || counts.MemoryLeaks > 10 {
		severityLevel = "MEDIUM - Investigation needed"
	} else {
		severityLevel = "LOW - Mostly harmless development errors"
	}

	var environment string
	if counts.NetworkIssues > 0 {
		environment = "Local development with occasional connectivity issues"
	} else {
		environment = "Stable local development environment"
	}

	return &types.AnalysisContext{
		TotalErrors:     totalErrors,
		MostCommonIssue: mostCommonIssue,
		SeverityLevel:   severityLevel,
		Environment:     environment,
	}
}

func (h *LogParserHandler) formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
