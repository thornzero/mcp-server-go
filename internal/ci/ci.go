package ci

import (
	"context"
	"os/exec"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thornzero/mcp-server-go/internal/models"
	"github.com/thornzero/mcp-server-go/internal/server"
	"github.com/thornzero/mcp-server-go/internal/types"
)

type CIHandler struct {
	server *server.Server
}

func NewCIHandler(s *server.Server) *CIHandler {
	return &CIHandler{server: s}
}

func (h *CIHandler) CIRunTests(ctx context.Context, req *mcp.CallToolRequest, input types.CIRunTestsInput) (*mcp.CallToolResult, types.CIRunTestsOutput, error) {
	scope := "./..."
	if input.Scope != nil && *input.Scope != "" {
		scope = *input.Scope
	}

	start := time.Now()
	cmd := exec.Command("go", "test", scope, "-count=1")
	cmd.Dir = h.server.GetRepoRoot()
	output, err := cmd.CombinedOutput()

	status := "pass"
	if err != nil {
		status = "fail"
	}

	// Log to database
	ciRun := models.CIRun{
		Scope:      scope,
		Status:     status,
		StartedAt:  start,
		FinishedAt: &[]time.Time{time.Now()}[0],
	}
	h.server.GetDB().Create(&ciRun)

	return nil, types.CIRunTestsOutput{
		Status: status,
		Output: string(output),
	}, nil
}

func (h *CIHandler) CILastFailure(ctx context.Context, req *mcp.CallToolRequest, input types.CILastFailureInput) (*mcp.CallToolResult, types.CILastFailureOutput, error) {
	var ciRun models.CIRun
	err := h.server.GetDB().
		Where("status = ?", "fail").
		Order("started_at DESC").
		First(&ciRun).Error

	if err != nil {
		return nil, types.CILastFailureOutput{Status: "none"}, nil
	}

	startedAt := ciRun.StartedAt.Format("2006-01-02 15:04:05")
	return nil, types.CILastFailureOutput{
		Status:    ciRun.Status,
		Scope:     &ciRun.Scope,
		StartedAt: &startedAt,
	}, nil
}
