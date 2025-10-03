package templates

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thornzero/project-manager/internal/markdown"
	"github.com/thornzero/project-manager/internal/models"
	"github.com/thornzero/project-manager/internal/server"
	"github.com/thornzero/project-manager/internal/types"
)

type TemplatesHandler struct {
	server *server.Server
}

func NewTemplatesHandler(s *server.Server) *TemplatesHandler {
	return &TemplatesHandler{server: s}
}

func (h *TemplatesHandler) TemplateList(ctx context.Context, req *mcp.CallToolRequest, input types.TemplateListInput) (*mcp.CallToolResult, types.TemplateListOutput, error) {
	var templates []models.MarkdownTemplate
	query := h.server.GetDB().Preload("Variables")

	if input.Category != nil && *input.Category != "" {
		query = query.Where("category = ?", *input.Category)
	}

	err := query.Order("category, name").Find(&templates).Error
	if err != nil {
		return nil, types.TemplateListOutput{}, err
	}

	// Convert to types.Template
	var resultTemplates []types.Template
	for _, tmpl := range templates {
		var variables []types.TemplateVariable
		for _, v := range tmpl.Variables {
			variables = append(variables, types.TemplateVariable{
				Name:         v.Name,
				Type:         v.Type,
				Required:     v.Required,
				DefaultValue: v.DefaultValue,
				Description:  v.Description,
			})
		}

		resultTemplates = append(resultTemplates, types.Template{
			ID:          tmpl.ID,
			Name:        tmpl.Name,
			Description: tmpl.Description,
			Category:    tmpl.Category,
			Content:     tmpl.Content,
			Variables:   variables,
			CreatedAt:   tmpl.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   tmpl.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	if resultTemplates == nil {
		resultTemplates = []types.Template{}
	}

	return nil, types.TemplateListOutput{Templates: resultTemplates}, nil
}

func (h *TemplatesHandler) TemplateRegister(ctx context.Context, req *mcp.CallToolRequest, input types.TemplateRegisterInput) (*mcp.CallToolResult, types.TemplateRegisterOutput, error) {
	if strings.TrimSpace(input.ID) == "" || strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Category) == "" || strings.TrimSpace(input.Content) == "" {
		return nil, types.TemplateRegisterOutput{}, fmt.Errorf("id, name, category, and content are required")
	}

	description := ""
	if input.Description != nil {
		description = *input.Description
	}

	// Convert variables
	var variables []models.TemplateVariable
	for _, v := range input.Variables {
		variables = append(variables, models.TemplateVariable{
			TemplateID:   input.ID,
			Name:         v.Name,
			Type:         v.Type,
			Required:     v.Required,
			DefaultValue: v.DefaultValue,
			Description:  v.Description,
		})
	}

	template := models.MarkdownTemplate{
		ID:          input.ID,
		Name:        input.Name,
		Description: description,
		Category:    input.Category,
		Content:     input.Content,
		Variables:   variables,
	}

	err := h.server.GetDB().Create(&template).Error
	if err != nil {
		return nil, types.TemplateRegisterOutput{}, err
	}

	return nil, types.TemplateRegisterOutput{ID: input.ID, Success: true}, nil
}

func (h *TemplatesHandler) TemplateGet(ctx context.Context, req *mcp.CallToolRequest, input types.TemplateGetInput) (*mcp.CallToolResult, types.TemplateGetOutput, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, types.TemplateGetOutput{}, fmt.Errorf("template ID is required")
	}

	var tmpl models.MarkdownTemplate
	err := h.server.GetDB().Preload("Variables").Where("id = ?", input.ID).First(&tmpl).Error
	if err != nil {
		return nil, types.TemplateGetOutput{}, fmt.Errorf("template not found: %s", input.ID)
	}

	// Convert variables
	var variables []types.TemplateVariable
	for _, v := range tmpl.Variables {
		variables = append(variables, types.TemplateVariable{
			Name:         v.Name,
			Type:         v.Type,
			Required:     v.Required,
			DefaultValue: v.DefaultValue,
			Description:  v.Description,
		})
	}

	resultTemplate := types.Template{
		ID:          tmpl.ID,
		Name:        tmpl.Name,
		Description: tmpl.Description,
		Category:    tmpl.Category,
		Content:     tmpl.Content,
		Variables:   variables,
		CreatedAt:   tmpl.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   tmpl.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return nil, types.TemplateGetOutput{Template: resultTemplate}, nil
}

func (h *TemplatesHandler) TemplateUpdate(ctx context.Context, req *mcp.CallToolRequest, input types.TemplateUpdateInput) (*mcp.CallToolResult, types.TemplateUpdateOutput, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, types.TemplateUpdateOutput{}, fmt.Errorf("template ID is required")
	}

	// Build updates map
	updates := make(map[string]interface{})
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.Category != nil {
		updates["category"] = *input.Category
	}
	if input.Content != nil {
		updates["content"] = *input.Content
	}

	// Update template if there are changes
	if len(updates) > 0 {
		result := h.server.GetDB().Model(&models.MarkdownTemplate{}).Where("id = ?", input.ID).Updates(updates)
		if result.Error != nil {
			return nil, types.TemplateUpdateOutput{}, result.Error
		}
	}

	// Update variables if provided
	if len(input.Variables) > 0 {
		// Delete existing variables
		h.server.GetDB().Where("template_id = ?", input.ID).Delete(&models.TemplateVariable{})

		// Insert new variables
		var variables []models.TemplateVariable
		for _, v := range input.Variables {
			variables = append(variables, models.TemplateVariable{
				TemplateID:   input.ID,
				Name:         v.Name,
				Type:         v.Type,
				Required:     v.Required,
				DefaultValue: v.DefaultValue,
				Description:  v.Description,
			})
		}
		h.server.GetDB().Create(&variables)
	}

	return nil, types.TemplateUpdateOutput{Updated: true}, nil
}

func (h *TemplatesHandler) TemplateDelete(ctx context.Context, req *mcp.CallToolRequest, input types.TemplateDeleteInput) (*mcp.CallToolResult, types.TemplateDeleteOutput, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, types.TemplateDeleteOutput{}, fmt.Errorf("template ID is required")
	}

	result := h.server.GetDB().Where("id = ?", input.ID).Delete(&models.MarkdownTemplate{})
	if result.Error != nil {
		return nil, types.TemplateDeleteOutput{}, result.Error
	}

	return nil, types.TemplateDeleteOutput{Deleted: result.RowsAffected > 0}, nil
}

func (h *TemplatesHandler) TemplateApply(ctx context.Context, req *mcp.CallToolRequest, input types.TemplateApplyInput) (*mcp.CallToolResult, types.TemplateApplyOutput, error) {
	if strings.TrimSpace(input.TemplateID) == "" {
		return nil, types.TemplateApplyOutput{}, fmt.Errorf("template ID is required")
	}

	// Get template
	var tmpl models.MarkdownTemplate
	err := h.server.GetDB().Where("id = ?", input.TemplateID).First(&tmpl).Error
	if err != nil {
		return nil, types.TemplateApplyOutput{}, fmt.Errorf("template not found: %s", input.TemplateID)
	}

	// Parse and execute template
	t, err := template.New(tmpl.ID).Parse(tmpl.Content)
	if err != nil {
		return nil, types.TemplateApplyOutput{}, fmt.Errorf("template parse error: %v", err)
	}

	var result strings.Builder
	err = t.Execute(&result, input.Variables)
	if err != nil {
		return nil, types.TemplateApplyOutput{}, fmt.Errorf("template execution error: %v", err)
	}

	content := result.String()

	// Apply line wrapping to prevent long lines
	content = h.wrapTemplateContent(content)

	outputPath := ""

	// Write to file if output path specified or if template should auto-write
	if input.OutputPath != nil && *input.OutputPath != "" {
		fullPath := filepath.Join(h.server.GetRepoRoot(), *input.OutputPath)
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			return nil, types.TemplateApplyOutput{}, fmt.Errorf("failed to write file: %v", err)
		}
		outputPath = fullPath
	} else if input.OutputPath != nil && *input.OutputPath == "" {
		// Auto-generate filename in docs directory
		docsDir := h.server.GetDocsOutputPath()
		_ = os.MkdirAll(docsDir, 0755) // Ensure directory exists

		// Generate filename from template ID
		filename := tmpl.ID + ".md"
		fullPath := filepath.Join(docsDir, filename)
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			return nil, types.TemplateApplyOutput{}, fmt.Errorf("failed to write file: %v", err)
		}
		outputPath = fullPath
	}

	// Lint the file after output (only if a file was written)
	if outputPath != "" {
		markdownHandler := markdown.NewMarkdownHandler(h.server)
		fix := true
		configPath := filepath.Join(h.server.GetRepoRoot(), ".markdownlint.json")
		_, _, err = markdownHandler.MarkdownLint(ctx, req, types.MarkdownLintInput{
			Path:   &outputPath,
			Fix:    &fix,
			Config: &configPath,
		})
		if err != nil {
			return nil, types.TemplateApplyOutput{}, fmt.Errorf("failed to lint file: %v", err)
		}
	}

	return nil, types.TemplateApplyOutput{Content: content, Path: outputPath}, nil
}

// wrapTemplateContent applies line wrapping to template-generated content
func (h *TemplatesHandler) wrapTemplateContent(content string) string {
	lines := strings.Split(content, "\n")
	var wrappedLines []string

	for _, line := range lines {
		// Skip wrapping for code blocks, headers, and list items
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "  ") {
			wrappedLines = append(wrappedLines, line)
			continue
		}

		// Wrap long lines (over 80 characters)
		if len(line) > 80 {
			wrappedText := markdown.WrapText(line, 80)
			wrappedLines = append(wrappedLines, wrappedText...)
		} else {
			wrappedLines = append(wrappedLines, line)
		}
	}

	return strings.Join(wrappedLines, "\n")
}
