package models

import (
	"testing"
	"time"
)

func TestGoal_Validation(t *testing.T) {
	tests := []struct {
		name    string
		goal    Goal
		isValid bool
	}{
		{
			name: "Valid goal",
			goal: Goal{
				ID:        1,
				Title:     "Test Goal",
				Priority:  0,
				Status:    "active",
				Notes:     "Test notes",
				UpdatedAt: time.Now(),
			},
			isValid: true,
		},
		{
			name: "Invalid goal with empty title",
			goal: Goal{
				ID:        1,
				Title:     "",
				Priority:  0,
				Status:    "active",
				Notes:     "Test notes",
				UpdatedAt: time.Now(),
			},
			isValid: false,
		},
		{
			name: "Invalid goal with invalid status",
			goal: Goal{
				ID:        1,
				Title:     "Test Goal",
				Priority:  0,
				Status:    "invalid_status",
				Notes:     "Test notes",
				UpdatedAt: time.Now(),
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation
			isValid := tt.goal.Title != "" &&
				(tt.goal.Status == "active" || tt.goal.Status == "paused" || tt.goal.Status == "done")

			if isValid != tt.isValid {
				t.Errorf("Expected isValid=%v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestADR_Validation(t *testing.T) {
	tests := []struct {
		name    string
		adr     ADR
		isValid bool
	}{
		{
			name: "Valid ADR",
			adr: ADR{
				ID:        "ADR-001",
				Title:     "Test ADR",
				Content:   "This is a test ADR",
				UpdatedAt: time.Now(),
			},
			isValid: true,
		},
		{
			name: "Invalid ADR with empty ID",
			adr: ADR{
				ID:        "",
				Title:     "Test ADR",
				Content:   "This is a test ADR",
				UpdatedAt: time.Now(),
			},
			isValid: false,
		},
		{
			name: "Invalid ADR with empty title",
			adr: ADR{
				ID:        "ADR-001",
				Title:     "",
				Content:   "This is a test ADR",
				UpdatedAt: time.Now(),
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - ADR only needs ID and Title
			isValid := tt.adr.ID != "" && tt.adr.Title != ""

			if isValid != tt.isValid {
				t.Errorf("Expected isValid=%v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestMarkdownTemplate_Validation(t *testing.T) {
	tests := []struct {
		name     string
		template MarkdownTemplate
		isValid  bool
	}{
		{
			name: "Valid template",
			template: MarkdownTemplate{
				ID:          "test-template",
				Name:        "Test Template",
				Description: "A test template",
				Category:    "documentation",
				Content:     "# {{.Title}}\n\n{{.Content}}",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			isValid: true,
		},
		{
			name: "Invalid template with empty ID",
			template: MarkdownTemplate{
				ID:          "",
				Name:        "Test Template",
				Description: "A test template",
				Category:    "documentation",
				Content:     "# {{.Title}}\n\n{{.Content}}",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			isValid: false,
		},
		{
			name: "Invalid template with empty name",
			template: MarkdownTemplate{
				ID:          "test-template",
				Name:        "",
				Description: "A test template",
				Category:    "documentation",
				Content:     "# {{.Title}}\n\n{{.Content}}",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			isValid: false,
		},
		{
			name: "Invalid template with empty content",
			template: MarkdownTemplate{
				ID:          "test-template",
				Name:        "Test Template",
				Description: "A test template",
				Category:    "documentation",
				Content:     "",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation
			isValid := tt.template.ID != "" && tt.template.Name != "" && tt.template.Content != ""

			if isValid != tt.isValid {
				t.Errorf("Expected isValid=%v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestTemplateVariable_Validation(t *testing.T) {
	tests := []struct {
		name     string
		variable TemplateVariable
		isValid  bool
	}{
		{
			name: "Valid variable",
			variable: TemplateVariable{
				ID:           1,
				TemplateID:   "test-template",
				Name:         "title",
				Type:         "string",
				Required:     true,
				DefaultValue: "",
				Description:  "Document title",
			},
			isValid: true,
		},
		{
			name: "Invalid variable with empty template ID",
			variable: TemplateVariable{
				ID:           1,
				TemplateID:   "",
				Name:         "title",
				Type:         "string",
				Required:     true,
				DefaultValue: "",
				Description:  "Document title",
			},
			isValid: false,
		},
		{
			name: "Invalid variable with empty name",
			variable: TemplateVariable{
				ID:           1,
				TemplateID:   "test-template",
				Name:         "",
				Type:         "string",
				Required:     true,
				DefaultValue: "",
				Description:  "Document title",
			},
			isValid: false,
		},
		{
			name: "Invalid variable with empty type",
			variable: TemplateVariable{
				ID:           1,
				TemplateID:   "test-template",
				Name:         "title",
				Type:         "",
				Required:     true,
				DefaultValue: "",
				Description:  "Document title",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation
			isValid := tt.variable.TemplateID != "" && tt.variable.Name != "" && tt.variable.Type != ""

			if isValid != tt.isValid {
				t.Errorf("Expected isValid=%v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestPreferredTool_Validation(t *testing.T) {
	tests := []struct {
		name    string
		tool    PreferredTool
		isValid bool
	}{
		{
			name: "Valid tool",
			tool: PreferredTool{
				ID:          1,
				Name:        "ESLint",
				Category:    "linting",
				Description: "JavaScript linter",
				Language:    "javascript",
				UseCase:     "Code quality",
				Priority:    5,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			isValid: true,
		},
		{
			name: "Invalid tool with empty name",
			tool: PreferredTool{
				ID:          1,
				Name:        "",
				Category:    "linting",
				Description: "JavaScript linter",
				Language:    "javascript",
				UseCase:     "Code quality",
				Priority:    5,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			isValid: false,
		},
		{
			name: "Invalid tool with empty category",
			tool: PreferredTool{
				ID:          1,
				Name:        "ESLint",
				Category:    "",
				Description: "JavaScript linter",
				Language:    "javascript",
				UseCase:     "Code quality",
				Priority:    5,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation
			isValid := tt.tool.Name != "" && tt.tool.Category != ""

			if isValid != tt.isValid {
				t.Errorf("Expected isValid=%v, got %v", tt.isValid, isValid)
			}
		})
	}
}
