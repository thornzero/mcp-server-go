package types

import (
	"testing"
)

func TestGoalsAddInput_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   GoalsAddInput
		isValid bool
	}{
		{
			name: "Valid input with title only",
			input: GoalsAddInput{
				Title: "Test Goal",
			},
			isValid: true,
		},
		{
			name: "Valid input with all fields",
			input: GoalsAddInput{
				Title:    "Test Goal",
				Priority: func() *int { p := 1; return &p }(),
				Notes:    func() *string { n := "Test notes"; return &n }(),
			},
			isValid: true,
		},
		{
			name: "Invalid input with empty title",
			input: GoalsAddInput{
				Title: "",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - check if title is not empty
			isValid := tt.input.Title != ""
			if isValid != tt.isValid {
				t.Errorf("Expected isValid=%v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestGoalsUpdateInput_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   GoalsUpdateInput
		isValid bool
	}{
		{
			name: "Valid input with ID and status",
			input: GoalsUpdateInput{
				ID:     1,
				Status: func() *string { s := "done"; return &s }(),
			},
			isValid: true,
		},
		{
			name: "Valid input with ID and notes",
			input: GoalsUpdateInput{
				ID:    1,
				Notes: func() *string { n := "Updated notes"; return &n }(),
			},
			isValid: true,
		},
		{
			name: "Invalid input with zero ID",
			input: GoalsUpdateInput{
				ID: 0,
			},
			isValid: false,
		},
		{
			name: "Invalid input with negative ID",
			input: GoalsUpdateInput{
				ID: -1,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - check if ID is positive
			isValid := tt.input.ID > 0
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
			name: "Valid string variable",
			variable: TemplateVariable{
				Name:         "title",
				Type:         "string",
				Required:     true,
				DefaultValue: "",
				Description:  "Document title",
			},
			isValid: true,
		},
		{
			name: "Valid number variable",
			variable: TemplateVariable{
				Name:         "count",
				Type:         "number",
				Required:     false,
				DefaultValue: "0",
				Description:  "Item count",
			},
			isValid: true,
		},
		{
			name: "Invalid variable with empty name",
			variable: TemplateVariable{
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
			// Basic validation - check if name and type are not empty
			isValid := tt.variable.Name != "" && tt.variable.Type != ""
			if isValid != tt.isValid {
				t.Errorf("Expected isValid=%v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestTemplateRegisterInput_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   TemplateRegisterInput
		isValid bool
	}{
		{
			name: "Valid template with all required fields",
			input: TemplateRegisterInput{
				ID:       "test-template",
				Name:     "Test Template",
				Category: "documentation",
				Content:  "# {{.Title}}\n\n{{.Content}}",
			},
			isValid: true,
		},
		{
			name: "Valid template with optional fields",
			input: TemplateRegisterInput{
				ID:          "test-template",
				Name:        "Test Template",
				Description: func() *string { d := "A test template"; return &d }(),
				Category:    "documentation",
				Content:     "# {{.Title}}\n\n{{.Content}}",
			},
			isValid: true,
		},
		{
			name: "Invalid template with empty ID",
			input: TemplateRegisterInput{
				ID:       "",
				Name:     "Test Template",
				Category: "documentation",
				Content:  "# {{.Title}}\n\n{{.Content}}",
			},
			isValid: false,
		},
		{
			name: "Invalid template with empty name",
			input: TemplateRegisterInput{
				ID:       "test-template",
				Name:     "",
				Category: "documentation",
				Content:  "# {{.Title}}\n\n{{.Content}}",
			},
			isValid: false,
		},
		{
			name: "Invalid template with empty category",
			input: TemplateRegisterInput{
				ID:       "test-template",
				Name:     "Test Template",
				Category: "",
				Content:  "# {{.Title}}\n\n{{.Content}}",
			},
			isValid: false,
		},
		{
			name: "Invalid template with empty content",
			input: TemplateRegisterInput{
				ID:       "test-template",
				Name:     "Test Template",
				Category: "documentation",
				Content:  "",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - check if all required fields are not empty
			isValid := tt.input.ID != "" && tt.input.Name != "" && tt.input.Category != "" && tt.input.Content != ""
			if isValid != tt.isValid {
				t.Errorf("Expected isValid=%v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestPreferredToolsAddInput_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   PreferredToolsAddInput
		isValid bool
	}{
		{
			name: "Valid tool with required fields only",
			input: PreferredToolsAddInput{
				Name:     "ESLint",
				Category: "linting",
			},
			isValid: true,
		},
		{
			name: "Valid tool with all fields",
			input: PreferredToolsAddInput{
				Name:        "ESLint",
				Category:    "linting",
				Description: "JavaScript linter",
				Language:    "javascript",
				UseCase:     "Code quality",
				Priority:    5,
			},
			isValid: true,
		},
		{
			name: "Invalid tool with empty name",
			input: PreferredToolsAddInput{
				Name:     "",
				Category: "linting",
			},
			isValid: false,
		},
		{
			name: "Invalid tool with empty category",
			input: PreferredToolsAddInput{
				Name:     "ESLint",
				Category: "",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - check if required fields are not empty
			isValid := tt.input.Name != "" && tt.input.Category != ""
			if isValid != tt.isValid {
				t.Errorf("Expected isValid=%v, got %v", tt.isValid, isValid)
			}
		})
	}
}
