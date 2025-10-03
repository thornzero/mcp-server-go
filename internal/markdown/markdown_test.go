package markdown

import (
	"strings"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder()
	if builder == nil {
		t.Errorf("NewBuilder() returned nil")
	}
	if builder.doc == nil {
		t.Errorf("NewBuilder() created builder with nil document")
	}
}

func TestBuilder_AddHeader(t *testing.T) {
	builder := NewBuilder()
	builder.AddHeader(1, "Test Header")

	result := builder.String()
	if !strings.Contains(result, "# Test Header") {
		t.Errorf("AddHeader() did not add header correctly. Got: %s", result)
	}
}

func TestBuilder_AddParagraph(t *testing.T) {
	builder := NewBuilder()
	builder.AddParagraph("Test paragraph")

	result := builder.String()
	if !strings.Contains(result, "Test paragraph") {
		t.Errorf("AddParagraph() did not add paragraph correctly. Got: %s", result)
	}
}

func TestBuilder_AddList(t *testing.T) {
	builder := NewBuilder()
	items := []string{"Item 1", "Item 2", "Item 3"}
	builder.AddList(items)

	result := builder.String()
	for _, item := range items {
		if !strings.Contains(result, "- "+item) {
			t.Errorf("AddList() did not add item '%s' correctly. Got: %s", item, result)
		}
	}
}

func TestBuilder_AddDefinitionList(t *testing.T) {
	builder := NewBuilder()
	items := map[string]string{
		"Term 1": "Definition 1",
		"Term 2": "Definition 2",
	}
	builder.AddDefinitionList(items)

	result := builder.String()

	// Check that terms are present without prefix
	if !strings.Contains(result, "Term 1") {
		t.Errorf("AddDefinitionList() did not add term 'Term 1'. Got: %s", result)
	}
	if !strings.Contains(result, "Term 2") {
		t.Errorf("AddDefinitionList() did not add term 'Term 2'. Got: %s", result)
	}

	// Check that definitions are present with : prefix
	if !strings.Contains(result, ": Definition 1") {
		t.Errorf("AddDefinitionList() did not add definition ': Definition 1'. Got: %s", result)
	}
	if !strings.Contains(result, ": Definition 2") {
		t.Errorf("AddDefinitionList() did not add definition ': Definition 2'. Got: %s", result)
	}
}

func TestBuilder_AddCodeBlock(t *testing.T) {
	builder := NewBuilder()
	builder.AddCodeBlock("go", "package main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}")

	result := builder.String()
	expected := "```go\npackage main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}\n```"
	if !strings.Contains(result, expected) {
		t.Errorf("AddCodeBlock() did not add code block correctly. Got: %s", result)
	}
}

func TestBuilder_AddInlineCode(t *testing.T) {
	builder := NewBuilder()
	builder.AddInlineCode("fmt.Println")

	result := builder.String()
	if !strings.Contains(result, "`fmt.Println`") {
		t.Errorf("AddInlineCode() did not add inline code correctly. Got: %s", result)
	}
}

func TestBuilder_AddBold(t *testing.T) {
	builder := NewBuilder()
	builder.AddBold("Bold Text")

	result := builder.String()
	if !strings.Contains(result, "**Bold Text**") {
		t.Errorf("AddBold() did not add bold text correctly. Got: %s", result)
	}
}

func TestBuilder_AddHorizontalRule(t *testing.T) {
	builder := NewBuilder()
	builder.AddHorizontalRule()

	result := builder.String()
	if !strings.Contains(result, "---") {
		t.Errorf("AddHorizontalRule() did not add horizontal rule correctly. Got: %s", result)
	}
}

func TestBuilder_AddLink(t *testing.T) {
	builder := NewBuilder()
	builder.AddLink("Google", "https://google.com")

	result := builder.String()
	if !strings.Contains(result, "[Google](https://google.com)") {
		t.Errorf("AddLink() did not add link correctly. Got: %s", result)
	}
}

func TestBuilder_AddSection(t *testing.T) {
	builder := NewBuilder()
	builder.AddSection(2, "Test Section", func(b *Builder) {
		b.AddParagraph("Section content")
	})

	result := builder.String()
	if !strings.Contains(result, "## Test Section") {
		t.Errorf("AddSection() did not add section header correctly. Got: %s", result)
	}
	if !strings.Contains(result, "Section content") {
		t.Errorf("AddSection() did not add section content correctly. Got: %s", result)
	}
}

func TestBuilder_String(t *testing.T) {
	builder := NewBuilder()
	builder.AddHeader(1, "Test")
	builder.AddParagraph("This is a test.")

	result := builder.String()
	if result == "" {
		t.Errorf("String() returned empty string")
	}
	if !strings.Contains(result, "# Test") {
		t.Errorf("String() did not contain expected header. Got: %s", result)
	}
	if !strings.Contains(result, "This is a test.") {
		t.Errorf("String() did not contain expected paragraph. Got: %s", result)
	}
}

func TestBuilder_WriteToFile(t *testing.T) {
	builder := NewBuilder()
	builder.AddHeader(1, "Test Document")
	builder.AddParagraph("This is a test document.")

	// Create a temporary file
	tempFile := "/tmp/test_markdown.md"
	err := builder.WriteToFile(tempFile)
	if err != nil {
		t.Errorf("WriteToFile() error: %v", err)
	}

	// Clean up
	defer func() {
		// Note: In a real test, you'd want to clean up the file
		// but we'll leave it for inspection
	}()
}
