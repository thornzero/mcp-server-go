// Package goals provides MCP tools for managing project goals and milestones.
//
// This package implements the goals management functionality for the MCP server,
// allowing users to create, list, update, and track project goals with priorities
// and status tracking.
//
// Example usage:
//
//	handler := NewGoalsHandler(server)
//	result, output, err := handler.GoalsList(ctx, req, types.GoalsListInput{Limit: 10})
package goals

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thornzero/project-manager/internal/models"
	"github.com/thornzero/project-manager/internal/server"
	"github.com/thornzero/project-manager/internal/types"
)

// GoalsHandler handles MCP tool requests for goal management operations.
//
// It provides methods for listing, adding, and updating project goals
// with proper validation and database persistence.
type GoalsHandler struct {
	server *server.Server
}

// NewGoalsHandler creates a new GoalsHandler instance with the provided server.
//
// The server instance is used for database access and configuration.
func NewGoalsHandler(s *server.Server) *GoalsHandler {
	return &GoalsHandler{server: s}
}

// GoalsList retrieves a list of active project goals.
//
// It returns goals that are not marked as "done", ordered by priority (ascending)
// and then by update time (descending). The number of results is limited by the
// input limit parameter, defaulting to 10 if not specified.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - req: MCP tool request (unused but required by interface)
//   - input: GoalsListInput containing optional limit parameter
//
// Returns:
//   - result: MCP call result with JSON response
//   - output: GoalsListOutput containing the list of goals
//   - err: Any error that occurred during retrieval
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
		return nil, types.GoalsAddOutput{}, fmt.Errorf("title is required and cannot be empty")
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
