// /home/thornzero/Repositories/mcp-server-go/internal/preferredtools/preferredtools.go
package preferredtools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thornzero/mcp-server-go/internal/models"
	"github.com/thornzero/mcp-server-go/internal/server"
	"github.com/thornzero/mcp-server-go/internal/types"
)

type PreferredToolsHandler struct {
	server *server.Server
}

func NewPreferredToolsHandler(srv *server.Server) *PreferredToolsHandler {
	return &PreferredToolsHandler{server: srv}
}

func (h *PreferredToolsHandler) PreferredToolsList(ctx context.Context, req *mcp.CallToolRequest, input types.PreferredToolsListInput) (*mcp.CallToolResult, types.PreferredToolsListOutput, error) {
	var tools []models.PreferredTool
	query := h.server.GetDB()

	// Apply filters
	if input.Category != "" {
		query = query.Where("category = ?", input.Category)
	}
	if input.Language != "" {
		query = query.Where("language = ?", input.Language)
	}

	// Apply limit
	if input.Limit > 0 {
		query = query.Limit(input.Limit)
	}

	// Order by priority (higher priority first), then by name
	err := query.Order("priority DESC, name ASC").Find(&tools).Error
	if err != nil {
		return nil, types.PreferredToolsListOutput{}, err
	}

	// Convert to output format
	var resultTools []types.PreferredTool
	for _, tool := range tools {
		resultTools = append(resultTools, types.PreferredTool{
			ID:          tool.ID,
			Name:        tool.Name,
			Category:    tool.Category,
			Description: tool.Description,
			Language:    tool.Language,
			UseCase:     tool.UseCase,
			Priority:    tool.Priority,
			CreatedAt:   tool.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   tool.UpdatedAt.Format(time.RFC3339),
		})
	}

	return nil, types.PreferredToolsListOutput{Tools: resultTools}, nil
}

func (h *PreferredToolsHandler) PreferredToolsAdd(ctx context.Context, req *mcp.CallToolRequest, input types.PreferredToolsAddInput) (*mcp.CallToolResult, types.PreferredToolsAddOutput, error) {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Category) == "" {
		return nil, types.PreferredToolsAddOutput{}, fmt.Errorf("name and category are required")
	}

	tool := models.PreferredTool{
		Name:        input.Name,
		Category:    input.Category,
		Description: input.Description,
		Language:    input.Language,
		UseCase:     input.UseCase,
		Priority:    input.Priority,
	}

	err := h.server.GetDB().Create(&tool).Error
	if err != nil {
		return nil, types.PreferredToolsAddOutput{}, err
	}

	return nil, types.PreferredToolsAddOutput{ID: tool.ID, Success: true}, nil
}

func (h *PreferredToolsHandler) PreferredToolsUpdate(ctx context.Context, req *mcp.CallToolRequest, input types.PreferredToolsUpdateInput) (*mcp.CallToolResult, types.PreferredToolsUpdateOutput, error) {
	var tool models.PreferredTool
	err := h.server.GetDB().First(&tool, input.ID).Error
	if err != nil {
		return nil, types.PreferredToolsUpdateOutput{}, fmt.Errorf("tool not found")
	}

	// Update fields if provided
	updates := make(map[string]interface{})
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Category != "" {
		updates["category"] = input.Category
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.Language != "" {
		updates["language"] = input.Language
	}
	if input.UseCase != "" {
		updates["use_case"] = input.UseCase
	}
	if input.Priority != 0 {
		updates["priority"] = input.Priority
	}

	if len(updates) > 0 {
		err = h.server.GetDB().Model(&tool).Updates(updates).Error
		if err != nil {
			return nil, types.PreferredToolsUpdateOutput{}, err
		}
	}

	return nil, types.PreferredToolsUpdateOutput{Success: true}, nil
}

func (h *PreferredToolsHandler) PreferredToolsDelete(ctx context.Context, req *mcp.CallToolRequest, input types.PreferredToolsDeleteInput) (*mcp.CallToolResult, types.PreferredToolsDeleteOutput, error) {
	err := h.server.GetDB().Delete(&models.PreferredTool{}, input.ID).Error
	if err != nil {
		return nil, types.PreferredToolsDeleteOutput{}, err
	}

	return nil, types.PreferredToolsDeleteOutput{Success: true}, nil
}
