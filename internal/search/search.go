package search

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thornzero/project-manager/internal/server"
	"github.com/thornzero/project-manager/internal/types"
)

type SearchHandler struct {
	server *server.Server
}

func NewSearchHandler(s *server.Server) *SearchHandler {
	return &SearchHandler{server: s}
}

func (h *SearchHandler) RepoSearch(ctx context.Context, req *mcp.CallToolRequest, input types.RepoSearchInput) (*mcp.CallToolResult, types.RepoSearchOutput, error) {
	if strings.TrimSpace(input.Q) == "" {
		return nil, types.RepoSearchOutput{}, fmt.Errorf("query required")
	}

	max := 50
	if input.Max != nil {
		max = *input.Max
	}

	target := h.server.GetRepoRoot()
	if input.Path != nil && *input.Path != "" {
		target = filepath.Join(h.server.GetRepoRoot(), *input.Path)
	}

	// Use ripgrep if available, else fallback to grep
	cmd := exec.Command("rg", "--line-number", "--no-heading", "--max-count", fmt.Sprint(max), input.Q, target)
	output, err := cmd.CombinedOutput()
	if err != nil && len(output) == 0 {
		cmd = exec.Command("grep", "-Rn", input.Q, target)
		output, _ = cmd.CombinedOutput()
	}

	lines := strings.Split(string(output), "\n")
	results := []types.SearchResult{} // Initialize as empty slice, not nil
	re := regexp.MustCompile(`^(.+?):(\d+):(.*)$`)

	for _, line := range lines {
		if matches := re.FindStringSubmatch(line); matches != nil {
			results = append(results, types.SearchResult{
				File:  strings.TrimPrefix(matches[1], h.server.GetRepoRoot()+"/"),
				Line:  matches[2],
				Match: strings.TrimSpace(matches[3]),
			})
			if len(results) >= max {
				break
			}
		}
	}

	return nil, types.RepoSearchOutput{Results: results}, nil
}
