package state

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thornzero/mcp-server-go/internal/server"
	"github.com/thornzero/mcp-server-go/internal/types"
)

type StateHandler struct {
	server *server.Server
}

func NewStateHandler(s *server.Server) *StateHandler {
	return &StateHandler{server: s}
}

func (h *StateHandler) StateLogChange(ctx context.Context, req *mcp.CallToolRequest, input types.StateLogChangeInput) (*mcp.CallToolResult, types.StateLogChangeOutput, error) {
	if strings.TrimSpace(input.Summary) == "" {
		return nil, types.StateLogChangeOutput{}, fmt.Errorf("summary required")
	}

	path := filepath.Join(h.server.GetRepoRoot(), "CHANGELOG_AGENT.md")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, types.StateLogChangeOutput{}, err
	}
	defer f.Close()

	ts := time.Now().Format(time.RFC3339)
	fmt.Fprintf(f, "- %s — %s\n", ts, input.Summary)
	if len(input.Files) > 0 {
		fmt.Fprintf(f, "  - files: %s\n", strings.Join(input.Files, ", "))
	}

	return nil, types.StateLogChangeOutput{OK: true}, nil
}
