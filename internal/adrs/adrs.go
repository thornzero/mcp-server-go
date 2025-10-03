package adrs

import (
	"context"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thornzero/mcp-server-go/internal/models"
	"github.com/thornzero/mcp-server-go/internal/server"
	"github.com/thornzero/mcp-server-go/internal/types"
)

type ADRsHandler struct {
	server *server.Server
}

func NewADRsHandler(s *server.Server) *ADRsHandler {
	return &ADRsHandler{server: s}
}

func (h *ADRsHandler) ADRsList(ctx context.Context, req *mcp.CallToolRequest, input types.ADRsListInput) (*mcp.CallToolResult, types.ADRsListOutput, error) {
	var adrs []models.ADR
	query := h.server.GetDB()

	if input.Query != nil && strings.TrimSpace(*input.Query) != "" {
		searchTerm := "%" + *input.Query + "%"
		query = query.Where("title LIKE ? OR id LIKE ? OR content LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	err := query.Find(&adrs).Error
	if err != nil {
		return nil, types.ADRsListOutput{}, err
	}

	// Convert to types.ADR
	var resultADRs []types.ADR
	for _, adr := range adrs {
		resultADRs = append(resultADRs, types.ADR{
			ID:        adr.ID,
			Title:     adr.Title,
			Content:   adr.Content,
			UpdatedAt: adr.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return nil, types.ADRsListOutput{ADRs: resultADRs}, nil
}

func (h *ADRsHandler) ADRsGet(ctx context.Context, req *mcp.CallToolRequest, input types.ADRsGetInput) (*mcp.CallToolResult, types.ADRsGetOutput, error) {
	var adr models.ADR
	err := h.server.GetDB().Where("id = ?", input.ID).First(&adr).Error
	if err != nil {
		return nil, types.ADRsGetOutput{}, err
	}

	return nil, types.ADRsGetOutput{
		ID:      adr.ID,
		Title:   adr.Title,
		Content: adr.Content,
	}, nil
}
