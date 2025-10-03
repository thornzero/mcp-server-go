// /home/thornzero/Repositories/project-manager/internal/cursorrules/cursorrules.go
package cursorrules

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thornzero/project-manager/internal/models"
	"github.com/thornzero/project-manager/internal/server"
	"github.com/thornzero/project-manager/internal/templates"
	"github.com/thornzero/project-manager/internal/types"
)

type CursorRulesHandler struct {
	server *server.Server
}

type ruleDetails struct {
	description        string
	globs              string
	language           string
	bestPractices      string
	securityGuidelines string
}

func NewCursorRulesHandler(srv *server.Server) *CursorRulesHandler {
	return &CursorRulesHandler{server: srv}
}

func (h *CursorRulesHandler) CursorRulesList(ctx context.Context, req *mcp.CallToolRequest, input types.CursorRulesListInput) (*mcp.CallToolResult, types.CursorRulesListOutput, error) {
	var rules []models.CursorRule
	query := h.server.GetDB()

	// Apply filters
	if input.Category != "" {
		query = query.Where("category = ?", input.Category)
	}
	if input.Tags != "" {
		query = query.Where("tags LIKE ?", "%"+input.Tags+"%")
	}
	if input.Source != "" {
		query = query.Where("source = ?", input.Source)
	}
	if input.Active != nil {
		query = query.Where("is_active = ?", *input.Active)
	}

	// Apply limit
	if input.Limit > 0 {
		query = query.Limit(input.Limit)
	}

	// Order by creation date (newest first)
	err := query.Order("created_at DESC").Find(&rules).Error
	if err != nil {
		return nil, types.CursorRulesListOutput{}, err
	}

	// Convert to output format
	var resultRules []types.CursorRule
	for _, rule := range rules {
		resultRules = append(resultRules, types.CursorRule{
			ID:          rule.ID,
			Name:        rule.Name,
			Category:    rule.Category,
			Description: rule.Description,
			Content:     rule.Content,
			Tags:        rule.Tags,
			Source:      rule.Source,
			IsActive:    rule.IsActive,
			CreatedAt:   rule.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   rule.UpdatedAt.Format(time.RFC3339),
		})
	}

	return nil, types.CursorRulesListOutput{Rules: resultRules}, nil
}

func (h *CursorRulesHandler) CursorRulesAdd(ctx context.Context, req *mcp.CallToolRequest, input types.CursorRulesAddInput) (*mcp.CallToolResult, types.CursorRulesAddOutput, error) {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Category) == "" || strings.TrimSpace(input.Content) == "" {
		return nil, types.CursorRulesAddOutput{}, fmt.Errorf("name, category, and content are required")
	}

	// Validate MDC format
	if err := h.validateMDCFormat(input.Content); err != nil {
		return nil, types.CursorRulesAddOutput{}, fmt.Errorf("invalid MDC format: %v", err)
	}

	// Set defaults
	source := input.Source
	if source == "" {
		source = "local"
	}

	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	rule := models.CursorRule{
		Name:        input.Name,
		Category:    input.Category,
		Description: input.Description,
		Content:     input.Content,
		Tags:        input.Tags,
		Source:      source,
		IsActive:    isActive,
	}

	err := h.server.GetDB().Create(&rule).Error
	if err != nil {
		return nil, types.CursorRulesAddOutput{}, err
	}

	return nil, types.CursorRulesAddOutput{ID: rule.ID, Success: true}, nil
}

func (h *CursorRulesHandler) CursorRulesUpdate(ctx context.Context, req *mcp.CallToolRequest, input types.CursorRulesUpdateInput) (*mcp.CallToolResult, types.CursorRulesUpdateOutput, error) {
	var rule models.CursorRule
	err := h.server.GetDB().First(&rule, input.ID).Error
	if err != nil {
		return nil, types.CursorRulesUpdateOutput{}, fmt.Errorf("rule not found")
	}

	// Validate MDC format if content is being updated
	if input.Content != "" {
		if err := h.validateMDCFormat(input.Content); err != nil {
			return nil, types.CursorRulesUpdateOutput{}, fmt.Errorf("invalid MDC format: %v", err)
		}
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
	if input.Content != "" {
		updates["content"] = input.Content
	}
	if input.Tags != "" {
		updates["tags"] = input.Tags
	}
	if input.Source != "" {
		updates["source"] = input.Source
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	if len(updates) > 0 {
		err = h.server.GetDB().Model(&rule).Updates(updates).Error
		if err != nil {
			return nil, types.CursorRulesUpdateOutput{}, err
		}
	}

	return nil, types.CursorRulesUpdateOutput{Success: true}, nil
}

func (h *CursorRulesHandler) CursorRulesDelete(ctx context.Context, req *mcp.CallToolRequest, input types.CursorRulesDeleteInput) (*mcp.CallToolResult, types.CursorRulesDeleteOutput, error) {
	err := h.server.GetDB().Delete(&models.CursorRule{}, input.ID).Error
	if err != nil {
		return nil, types.CursorRulesDeleteOutput{}, err
	}

	return nil, types.CursorRulesDeleteOutput{Success: true}, nil
}

func (h *CursorRulesHandler) CursorRulesSuggest(ctx context.Context, req *mcp.CallToolRequest, input types.CursorRulesSuggestInput) (*mcp.CallToolResult, types.CursorRulesSuggestOutput, error) {
	// Get properly formatted example suggestions following current guidelines
	suggestions := h.getExampleRules()

	// Filter suggestions based on input criteria
	var filteredSuggestions []types.CommunityRule
	for _, suggestion := range suggestions {
		include := true

		if input.Language != "" && !strings.Contains(strings.ToLower(suggestion.Tags), strings.ToLower(input.Language)) {
			include = false
		}
		if input.Category != "" && !strings.EqualFold(suggestion.Category, input.Category) {
			include = false
		}
		if input.Tags != "" {
			inputTags := strings.Split(input.Tags, ",")
			found := false
			for _, tag := range inputTags {
				if strings.Contains(strings.ToLower(suggestion.Tags), strings.TrimSpace(strings.ToLower(tag))) {
					found = true
					break
				}
			}
			if !found {
				include = false
			}
		}

		if include {
			filteredSuggestions = append(filteredSuggestions, suggestion)
		}
	}

	// Apply limit
	if input.Limit > 0 && len(filteredSuggestions) > input.Limit {
		filteredSuggestions = filteredSuggestions[:input.Limit]
	}

	return nil, types.CursorRulesSuggestOutput{Suggestions: filteredSuggestions}, nil
}

func (h *CursorRulesHandler) CursorRulesInstall(ctx context.Context, req *mcp.CallToolRequest, input types.CursorRulesInstallInput) (*mcp.CallToolResult, types.CursorRulesInstallOutput, error) {
	// Create the .cursor/rules directory if it doesn't exist
	rulesDir := filepath.Join(h.server.GetRepoRoot(), ".cursor", "rules")
	err := os.MkdirAll(rulesDir, 0755)
	if err != nil {
		return nil, types.CursorRulesInstallOutput{}, fmt.Errorf("failed to create rules directory: %v", err)
	}

	// Fetch the rule content from GitHub if URL is provided
	var ruleContent string
	if input.URL != "" {
		content, err := h.fetchRuleFromGitHub(input.URL)
		if err != nil {
			return nil, types.CursorRulesInstallOutput{}, fmt.Errorf("failed to fetch rule: %v", err)
		}
		ruleContent = content
	} else {
		// For now, create a placeholder rule with proper MDC format
		ruleContent = h.generatePlaceholderRule(input.RuleName)
	}

	// Validate the MDC format
	if err := h.validateMDCFormat(ruleContent); err != nil {
		return nil, types.CursorRulesInstallOutput{}, fmt.Errorf("downloaded rule has invalid format: %v", err)
	}

	// Save the rule to the .cursor/rules directory
	filename := strings.ToLower(strings.ReplaceAll(input.RuleName, " ", "-")) + ".mdc"
	rulePath := filepath.Join(rulesDir, filename)
	err = os.WriteFile(rulePath, []byte(ruleContent), 0644)
	if err != nil {
		return nil, types.CursorRulesInstallOutput{}, fmt.Errorf("failed to write rule file: %v", err)
	}

	return nil, types.CursorRulesInstallOutput{
		Success: true,
		Message: fmt.Sprintf("Rule '%s' installed successfully to %s", input.RuleName, rulePath),
	}, nil
}

// Helper function to fetch rule content from GitHub (placeholder)
func (h *CursorRulesHandler) fetchRuleFromGitHub(ruleURL string) (string, error) {
	resp, err := http.Get(ruleURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch rule: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// validateMDCFormat validates that the content follows proper MDC format
func (h *CursorRulesHandler) validateMDCFormat(content string) error {
	lines := strings.Split(content, "\n")

	// Check for frontmatter
	if len(lines) < 3 {
		return fmt.Errorf("content must have frontmatter with at least description")
	}

	if !strings.HasPrefix(lines[0], "---") {
		return fmt.Errorf("content must start with frontmatter delimiter '---'")
	}

	// Find end of frontmatter
	frontmatterEnd := -1
	for i, line := range lines {
		if i > 0 && strings.HasPrefix(line, "---") {
			frontmatterEnd = i
			break
		}
	}

	if frontmatterEnd == -1 {
		return fmt.Errorf("content must have closing frontmatter delimiter '---'")
	}

	// Check for required frontmatter fields
	frontmatter := strings.Join(lines[1:frontmatterEnd], "\n")
	hasDescription := strings.Contains(frontmatter, "description:")
	hasGlobs := strings.Contains(frontmatter, "globs:")
	hasAlwaysApply := strings.Contains(frontmatter, "alwaysApply:")

	if !hasDescription {
		return fmt.Errorf("frontmatter must include 'description' field")
	}

	if !hasGlobs {
		return fmt.Errorf("frontmatter must include 'globs' field")
	}

	if !hasAlwaysApply {
		return fmt.Errorf("frontmatter must include 'alwaysApply' field")
	}

	// Check that there's content after frontmatter
	if len(lines) <= frontmatterEnd+1 {
		return fmt.Errorf("content must have rules after frontmatter")
	}

	return nil
}

// getExampleRules returns properly formatted example rules following current guidelines
func (h *CursorRulesHandler) getExampleRules() []types.CommunityRule {
	return []types.CommunityRule{
		{
			Name:        "Go Best Practices",
			Category:    "language",
			Description: "Best practices for Go development including error handling, naming conventions, and performance",
			Tags:        "go,golang,best-practices",
			Source:      "awesome-cursor-rules-mdc",
			URL:         "https://github.com/sanjeed5/awesome-cursor-rules-mdc/blob/main/rules-mdc/go.mdc",
		},
		{
			Name:        "React Development",
			Category:    "framework",
			Description: "React development rules including hooks, component patterns, and performance optimization",
			Tags:        "react,javascript,frontend",
			Source:      "awesome-cursor-rules-mdc",
			URL:         "https://github.com/sanjeed5/awesome-cursor-rules-mdc/blob/main/rules-mdc/react.mdc",
		},
		{
			Name:        "TypeScript Guidelines",
			Category:    "language",
			Description: "TypeScript development guidelines including type safety and modern features",
			Tags:        "typescript,javascript,types",
			Source:      "awesome-cursor-rules-mdc",
			URL:         "https://github.com/sanjeed5/awesome-cursor-rules-mdc/blob/main/rules-mdc/typescript.mdc",
		},
		{
			Name:        "Python Standards",
			Category:    "language",
			Description: "Python development standards including PEP 8 compliance and best practices",
			Tags:        "python,pep8,standards",
			Source:      "awesome-cursor-rules-mdc",
			URL:         "https://github.com/sanjeed5/awesome-cursor-rules-mdc/blob/main/rules-mdc/python.mdc",
		},
		{
			Name:        "Security Best Practices",
			Category:    "security",
			Description: "General security best practices for application development",
			Tags:        "security,best-practices,owasp",
			Source:      "awesome-cursor-rules-mdc",
			URL:         "https://github.com/sanjeed5/awesome-cursor-rules-mdc/blob/main/rules-mdc/security.mdc",
		},
	}
}

// getRuleDetails determines rule configuration based on rule name
func (h *CursorRulesHandler) getRuleDetails(ruleName string) ruleDetails {
	lowerName := strings.ToLower(ruleName)

	switch {
	case strings.Contains(lowerName, "go"):
		return ruleDetails{
			description:        "Go development guidelines and best practices",
			globs:              `["*.go", "*.md"]`,
			language:           "Go",
			bestPractices:      "Use meaningful variable names, handle errors explicitly, use context for cancellation, follow Go idioms",
			securityGuidelines: "Validate all inputs, use secure HTTP clients, avoid SQL injection with parameterized queries",
		}
	case strings.Contains(lowerName, "react"):
		return ruleDetails{
			description:        "React development rules and patterns",
			globs:              `["*.js", "*.jsx", "*.ts", "*.tsx", "*.md"]`,
			language:           "React",
			bestPractices:      "Use functional components with hooks, implement proper state management, optimize re-renders",
			securityGuidelines: "Sanitize user inputs, use CSP headers, validate props, avoid XSS vulnerabilities",
		}
	case strings.Contains(lowerName, "typescript"):
		return ruleDetails{
			description:        "TypeScript development standards",
			globs:              `["*.ts", "*.tsx", "*.md"]`,
			language:           "TypeScript",
			bestPractices:      "Use strict type checking, leverage advanced types, implement proper interfaces",
			securityGuidelines: "Use strict type checking to prevent runtime errors, validate external data",
		}
	case strings.Contains(lowerName, "python"):
		return ruleDetails{
			description:        "Python development standards and PEP compliance",
			globs:              `["*.py", "*.md"]`,
			language:           "Python",
			bestPractices:      "Follow PEP 8, use type hints, implement proper exception handling",
			securityGuidelines: "Validate inputs, use secure coding practices, avoid eval() and exec()",
		}
	case strings.Contains(lowerName, "security"):
		return ruleDetails{
			description:        "Security best practices and guidelines",
			globs:              `["*.go", "*.js", "*.ts", "*.py", "*.md"]`,
			language:           "General",
			bestPractices:      "Follow OWASP guidelines, implement proper authentication and authorization",
			securityGuidelines: "Validate all inputs, use secure communication, implement proper error handling",
		}
	default:
		return ruleDetails{
			description:        "Development guidelines and best practices",
			globs:              `["*.md"]`,
			language:           "General",
			bestPractices:      "Follow established patterns and conventions",
			securityGuidelines: "Follow security best practices and validate all inputs",
		}
	}
}

// generatePlaceholderRule creates a properly formatted MDC rule using the template system
func (h *CursorRulesHandler) generatePlaceholderRule(ruleName string) string {
	// Determine rule details based on name
	ruleDetails := h.getRuleDetails(ruleName)

	// Create a templates handler to use the template system
	templatesHandler := templates.NewTemplatesHandler(h.server)

	// Prepare template variables
	variables := map[string]interface{}{
		"rule_name":           ruleName,
		"description":         ruleDetails.description,
		"globs":               ruleDetails.globs,
		"language":            ruleDetails.language,
		"best_practices":      ruleDetails.bestPractices,
		"security_guidelines": ruleDetails.securityGuidelines,
	}

	// Apply the cursor rule template
	_, result, err := templatesHandler.TemplateApply(context.Background(), nil, types.TemplateApplyInput{
		TemplateID: "cursor-rule-template",
		Variables:  variables,
	})

	if err != nil {
		// Fallback to a basic template if template application fails
		return fmt.Sprintf(`---
description: %s
globs: %s
alwaysApply: true
---

# %s

## Overview
This rule provides guidelines for %s development.

## Key Guidelines
- Follow consistent naming conventions
- Write clear, self-documenting code
- Handle errors appropriately
- Write tests for new functionality`, ruleDetails.description, ruleDetails.globs, ruleName, ruleDetails.language)
	}

	return result.Content
}
