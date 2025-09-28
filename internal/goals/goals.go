package goals

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thornzero/mcp-server-go/internal/models"
	"github.com/thornzero/mcp-server-go/internal/server"
	"github.com/thornzero/mcp-server-go/internal/types"
)

type GoalsHandler struct {
	server *server.Server
}

func NewGoalsHandler(s *server.Server) *GoalsHandler {
	return &GoalsHandler{server: s}
}

func (h *GoalsHandler) GoalsList(ctx context.Context, req *mcp.CallToolRequest, input types.GoalsListInput) (*mcp.CallToolResult, types.GoalsListOutput, error) {
	limit := input.Limit
	if limit == 0 {
		limit = 10
	}

	var goals []models.Goal
	err := h.server.GetDB().
		Where("status != ?", "done").
		Order("priority ASC, updated_at DESC").
		Limit(limit).
		Find(&goals).Error
	if err != nil {
		return nil, types.GoalsListOutput{}, err
	}

	// Convert to types.Goal
	var resultGoals []types.Goal
	for _, g := range goals {
		resultGoals = append(resultGoals, types.Goal{
			ID:        int(g.ID),
			Title:     g.Title,
			Priority:  g.Priority,
			Status:    g.Status,
			Notes:     g.Notes,
			UpdatedAt: g.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return nil, types.GoalsListOutput{Goals: resultGoals}, nil
}

func (h *GoalsHandler) GoalsAdd(ctx context.Context, req *mcp.CallToolRequest, input types.GoalsAddInput) (*mcp.CallToolResult, types.GoalsAddOutput, error) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, types.GoalsAddOutput{}, fmt.Errorf("title required")
	}

	prio := 100
	if input.Priority != nil {
		prio = *input.Priority
	}
	notes := ""
	if input.Notes != nil {
		notes = *input.Notes
	}

	goal := models.Goal{
		Title:    input.Title,
		Priority: prio,
		Notes:    notes,
		Status:   "active",
	}

	err := h.server.GetDB().Create(&goal).Error
	if err != nil {
		return nil, types.GoalsAddOutput{}, err
	}

	return nil, types.GoalsAddOutput{ID: int(goal.ID)}, nil
}

func (h *GoalsHandler) GoalsUpdate(ctx context.Context, req *mcp.CallToolRequest, input types.GoalsUpdateInput) (*mcp.CallToolResult, types.GoalsUpdateOutput, error) {
	if input.ID == 0 {
		return nil, types.GoalsUpdateOutput{}, fmt.Errorf("id required")
	}

	updates := make(map[string]interface{})

	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.Notes != nil {
		updates["notes"] = *input.Notes
	}
	if input.Priority != nil {
		updates["priority"] = *input.Priority
	}

	if len(updates) == 0 {
		return nil, types.GoalsUpdateOutput{Updated: 0}, nil
	}

	result := h.server.GetDB().Model(&models.Goal{}).Where("id = ?", input.ID).Updates(updates)
	if result.Error != nil {
		return nil, types.GoalsUpdateOutput{}, result.Error
	}

	return nil, types.GoalsUpdateOutput{Updated: int(result.RowsAffected)}, nil
}
