// /home/thornzero/Repositories/mcp-server-go/internal/markdown/builder.go
package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
)

// Builder represents a markdown document builder using gomarkdown AST
type Builder struct {
	doc *ast.Document
}

// NewBuilder creates a new markdown builder
func NewBuilder() *Builder {
	return &Builder{
		doc: &ast.Document{},
	}
}

// AddHeader adds a header to the document
func (b *Builder) AddHeader(level int, textContent string) *Builder {
	heading := &ast.Heading{
		Level: level,
	}
	heading.Literal = []byte(textContent)
	ast.AppendChild(b.doc, heading)
	return b
}

// AddParagraph adds a paragraph
func (b *Builder) AddParagraph(textContent string) *Builder {
	paragraph := &ast.Paragraph{}
	paragraph.Literal = []byte(textContent)
	ast.AppendChild(b.doc, paragraph)
	return b
}

// AddList adds a list
func (b *Builder) AddList(items []string) *Builder {
	list := &ast.List{
		ListFlags: 0, // 0 = bullet list
	}

	for _, item := range items {
		listItem := &ast.ListItem{}
		paragraph := &ast.Paragraph{}
		paragraph.Literal = []byte(item)
		ast.AppendChild(listItem, paragraph)
		ast.AppendChild(list, listItem)
	}

	ast.AppendChild(b.doc, list)
	return b
}

// AddCodeBlock adds a code block
func (b *Builder) AddCodeBlock(language, code string) *Builder {
	codeBlock := &ast.CodeBlock{}
	codeBlock.Info = []byte(language)
	codeBlock.Literal = []byte(code)
	ast.AppendChild(b.doc, codeBlock)
	return b
}

// AddInlineCode adds inline code
func (b *Builder) AddInlineCode(code string) *Builder {
	codeSpan := &ast.Code{}
	codeSpan.Literal = []byte(code)
	ast.AppendChild(b.doc, codeSpan)
	return b
}

// AddBold adds bold text
func (b *Builder) AddBold(textContent string) *Builder {
	strong := &ast.Strong{}
	strong.Literal = []byte(textContent)
	ast.AppendChild(b.doc, strong)
	return b
}

// AddLineBreak adds a line break
func (b *Builder) AddLineBreak() *Builder {
	paragraph := &ast.Paragraph{}
	paragraph.Literal = []byte("\n")
	ast.AppendChild(b.doc, paragraph)
	return b
}

// AddHorizontalRule adds a horizontal rule
func (b *Builder) AddHorizontalRule() *Builder {
	hr := &ast.HorizontalRule{}
	ast.AppendChild(b.doc, hr)
	return b
}

// AddLink adds a link
func (b *Builder) AddLink(textContent, url string) *Builder {
	link := &ast.Link{}
	link.Destination = []byte(url)
	link.Literal = []byte(textContent)
	ast.AppendChild(b.doc, link)
	return b
}

// AddTable adds a table (simplified - tables require extension)
func (b *Builder) AddTable(headers []string, rows [][]string) *Builder {
	// For now, create a simple text representation
	// Tables require the table extension which is complex to implement
	tableText := "| " + strings.Join(headers, " | ") + " |\n"
	separator := make([]string, len(headers))
	for i := range headers {
		separator[i] = "---"
	}
	tableText += "| " + strings.Join(separator, " | ") + " |\n"
	for _, row := range rows {
		tableText += "| " + strings.Join(row, " | ") + " |\n"
	}

	paragraph := &ast.Paragraph{}
	paragraph.Literal = []byte(tableText)
	ast.AppendChild(b.doc, paragraph)
	return b
}

// AddChecklist adds a checklist
func (b *Builder) AddChecklist(items []string, checked []bool) *Builder {
	list := &ast.List{
		ListFlags: 0, // 0 = bullet list
	}

	for i, item := range items {
		listItem := &ast.ListItem{}

		// Add checkbox
		checkbox := "- [ ]"
		if i < len(checked) && checked[i] {
			checkbox = "- [x]"
		}

		paragraph := &ast.Paragraph{}
		paragraph.Literal = []byte(checkbox + " " + item)
		ast.AppendChild(listItem, paragraph)
		ast.AppendChild(list, listItem)
	}

	ast.AppendChild(b.doc, list)
	return b
}

// AddSection adds a section with header and content
func (b *Builder) AddSection(level int, title string, content func(*Builder)) *Builder {
	b.AddHeader(level, title)
	content(b)
	return b
}

// String returns the markdown content as a string
func (b *Builder) String() string {
	var result strings.Builder

	// Walk the AST and build markdown string
	ast.WalkFunc(b.doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.GoToNext
		}

		switch n := node.(type) {
		case *ast.Heading:
			level := n.Level
			prefix := strings.Repeat("#", level)
			text := string(n.Literal)
			result.WriteString(fmt.Sprintf("%s %s\n\n", prefix, text))
		case *ast.Paragraph:
			// Skip paragraphs that are inside list items
			if n.GetParent() != nil {
				if _, ok := n.GetParent().(*ast.ListItem); ok {
					return ast.GoToNext
				}
			}
			text := string(n.Literal)
			result.WriteString(fmt.Sprintf("%s\n\n", text))
		case *ast.List:
			// Lists are handled by their items
		case *ast.ListItem:
			// ListItem contains a Paragraph, so we need to get the text from the paragraph
			if len(n.GetChildren()) > 0 {
				if para, ok := n.GetChildren()[0].(*ast.Paragraph); ok {
					text := string(para.Literal)
					result.WriteString(fmt.Sprintf("- %s\n", text))
				}
			}
		case *ast.CodeBlock:
			lang := string(n.Info)
			code := string(n.Literal)
			result.WriteString(fmt.Sprintf("```%s\n%s\n```\n\n", lang, code))
		case *ast.Code:
			text := string(n.Literal)
			result.WriteString(fmt.Sprintf("`%s`", text))
		case *ast.Strong:
			text := string(n.Literal)
			result.WriteString(fmt.Sprintf("**%s**", text))
		case *ast.HorizontalRule:
			result.WriteString("---\n\n")
		case *ast.Link:
			text := string(n.Literal)
			url := string(n.Destination)
			result.WriteString(fmt.Sprintf("[%s](%s)", text, url))
		}

		return ast.GoToNext
	})

	return result.String()
}

// WriteToFile writes the markdown content to a file
func (b *Builder) WriteToFile(filePath string) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write content to file
	return os.WriteFile(filePath, []byte(b.String()), 0644)
}

// ValidateMarkdown validates the generated markdown using gomarkdown
func (b *Builder) ValidateMarkdown() error {
	markdownContent := []byte(b.String())

	// Try to parse the markdown
	doc := markdown.Parse(markdownContent, nil)
	if doc == nil {
		return fmt.Errorf("markdown validation failed: could not parse generated markdown")
	}

	return nil
}

// MCPToolsRuleBuilder builds the MCP tools rule markdown
func MCPToolsRuleBuilder(mcpServerPath string) *Builder {
	b := NewBuilder()

	// Front matter
	b.AddParagraph("---")
	b.AddParagraph("description: MCP Server Tools - Project Management and Development Assistance")
	b.AddParagraph("globs: [\"**/*\"]")
	b.AddParagraph("alwaysApply: true")
	b.AddParagraph("---")
	b.AddLineBreak()

	// Main header
	b.AddHeader(1, "MCP Server Tools")

	// Available tools section
	b.AddSection(2, "🎯 Available Tools", func(b *Builder) {
		b.AddParagraph("This project has access to MCP (Model Context Protocol) server tools for enhanced project management and AI assistance.")

		// Goals Management
		b.AddSection(3, "📋 Goals Management", func(b *Builder) {
			b.AddList([]string{
				"mcp_mcp-server-go_goals_list() - List active project goals",
				"mcp_mcp-server-go_goals_add({title: \"Goal Title\"}) - Add new goal",
				"mcp_mcp-server-go_goals_update({id: 1, status: \"done\"}) - Update goal status",
			})
		})

		// Cursor Rules Management
		b.AddSection(3, "📝 Cursor Rules Management", func(b *Builder) {
			b.AddList([]string{
				"mcp_mcp-server-go_cursor_rules_list() - List active cursor rules",
				"mcp_mcp-server-go_cursor_rules_add({name: \"Rule Name\", category: \"general\", content: \"Rule content\"}) - Add new rule",
				"mcp_mcp-server-go_cursor_rules_update({id: 1, content: \"Updated content\"}) - Update existing rule",
			})
		})

		// Documentation & ADRs
		b.AddSection(3, "📚 Documentation & ADRs", func(b *Builder) {
			b.AddList([]string{
				"mcp_mcp-server-go_adrs_list() - List Architecture Decision Records",
				"mcp_mcp-server-go_adrs_get({id: \"ADR-001\"}) - Get specific ADR content",
				"mcp_mcp-server-go_template_list() - List available documentation templates",
				"mcp_mcp-server-go_template_apply({template_id: \"template-name\", variables: {}}) - Generate documentation",
			})
		})

		// Repository Tools
		b.AddSection(3, "🔍 Repository Tools", func(b *Builder) {
			b.AddList([]string{
				"mcp_mcp-server-go_repo_search({q: \"search pattern\"}) - Search codebase",
				"mcp_mcp-server-go_markdown_lint({path: \"docs/\"}) - Lint markdown files",
				"mcp_mcp-server-go_state_log_change({summary: \"Change description\", files: [\"file1.go\"]}) - Log project changes",
			})
		})

		// CI & Testing
		b.AddSection(3, "🧪 CI & Testing", func(b *Builder) {
			b.AddList([]string{
				"mcp_mcp-server-go_ci_run_tests({scope: \"./cmd\"}) - Run tests",
				"mcp_mcp-server-go_ci_last_failure() - Get last test failure info",
			})
		})
	})

	// Getting Started section
	b.AddSection(2, "🚀 Getting Started", func(b *Builder) {
		b.AddSection(3, "1. Initialize Project Data", func(b *Builder) {
			b.AddParagraph("First, add initial project data:")
			b.AddCodeBlock("javascript", `// Add initial project goal
mcp_mcp-server-go_goals_add({
  title: "Project Setup Complete"
})

// Add project guidelines rule
mcp_mcp-server-go_cursor_rules_add({
  name: "Project Guidelines",
  category: "general",
  content: "Your project guidelines here..."
})

// Verify tools work
mcp_mcp-server-go_goals_list()
mcp_mcp-server-go_cursor_rules_list()`)
		})

		b.AddSection(3, "2. Use Tools Regularly", func(b *Builder) {
			b.AddList([]string{
				"Track project goals and milestones",
				"Maintain development rules and guidelines",
				"Document architecture decisions",
				"Log significant changes",
			})
		})
	})

	// Important Notes section
	b.AddSection(2, "⚠️ Important Notes", func(b *Builder) {
		b.AddList([]string{
			"**Always initialize with test data** before using list operations",
			"**Use required parameters only** initially (avoid optional parameters)",
			"**Restart Cursor** if tools show \"Not connected\" errors",
			"**Check troubleshooting guide** if tools fail",
		})
	})

	// Troubleshooting section
	b.AddSection(2, "🔧 Troubleshooting", func(b *Builder) {
		b.AddParagraph("If tools don't work:")
		b.AddList([]string{
			"Check if MCP server is configured in Cursor settings",
			"Restart Cursor completely",
			"Verify database exists: .agent/state.db",
			"Try adding test data first: mcp_mcp-server-go_goals_add({title: \"Test\"})",
		})
	})

	// Documentation section
	b.AddSection(2, "📚 Documentation", func(b *Builder) {
		b.AddList([]string{
			"**Usage Guide**: See mcp-tools-usage.mdc",
			"**Troubleshooting**: See mcp-tools-troubleshooting.mdc",
			fmt.Sprintf("**Full Documentation**: %s/docs/", mcpServerPath),
		})
	})

	// Footer
	b.AddHorizontalRule()
	b.AddParagraph(fmt.Sprintf("**Last Updated**: %s", time.Now().Format("2006-01-02")))
	b.AddParagraph("**Status**: Production Ready")

	return b
}

// MCPUsageGuideBuilder builds the MCP usage guide markdown
func MCPUsageGuideBuilder() *Builder {
	b := NewBuilder()

	// Front matter
	b.AddParagraph("---")
	b.AddParagraph("description: Comprehensive guide for using MCP server tools effectively")
	b.AddParagraph("globs: [\"**/*\"]")
	b.AddParagraph("alwaysApply: false")
	b.AddParagraph("---")
	b.AddLineBreak()

	// Main header
	b.AddHeader(1, "MCP Tools Usage Guide")

	// Purpose section
	b.AddSection(2, "🎯 Purpose", func(b *Builder) {
		b.AddParagraph("MCP (Model Context Protocol) tools provide AI agents with project management capabilities, allowing them to:")
		b.AddList([]string{
			"Track goals and milestones",
			"Maintain development rules",
			"Document architecture decisions",
			"Search and analyze code",
			"Generate documentation",
		})
	})

	// Core Workflows section
	b.AddSection(2, "📋 Core Workflows", func(b *Builder) {
		b.AddSection(3, "Project Initialization", func(b *Builder) {
			b.AddCodeBlock("javascript", `// 1. Add initial project goal
mcp_mcp-server-go_goals_add({
  title: "Project Setup Complete"
})

// 2. Add project guidelines
mcp_mcp-server-go_cursor_rules_add({
  name: "Project Guidelines",
  category: "general",
  content: "Your project guidelines here..."
})

// 3. Verify tools work
mcp_mcp-server-go_goals_list()
mcp_mcp-server-go_cursor_rules_list()`)
		})

		b.AddSection(3, "Goal Management", func(b *Builder) {
			b.AddCodeBlock("javascript", `// Add goals for different phases
mcp_mcp-server-go_goals_add({title: "Phase 1: Core Features"})
mcp_mcp-server-go_goals_add({title: "Phase 2: Performance"})
mcp_mcp-server-go_goals_add({title: "Phase 3: Polish"})

// Update goal status
mcp_mcp-server-go_goals_update({
  id: 1,
  status: "done"
})

// List all goals
mcp_mcp-server-go_goals_list()`)
		})

		b.AddSection(3, "Rule Management", func(b *Builder) {
			b.AddCodeBlock("javascript", `// Add coding standards
mcp_mcp-server-go_cursor_rules_add({
  name: "Code Quality Standards",
  category: "quality",
  content: "Code quality guidelines..."
})

// Add specific technology rules
mcp_mcp-server-go_cursor_rules_add({
  name: "React Best Practices",
  category: "framework",
  content: "React development guidelines..."
})`)
		})
	})

	// Best Practices section
	b.AddSection(2, "📊 Best Practices", func(b *Builder) {
		b.AddSection(3, "1. Always Initialize First", func(b *Builder) {
			b.AddList([]string{
				"Add test data before using list operations",
				"Verify tools work with simple operations",
			})
		})

		b.AddSection(3, "2. Use Required Parameters Only", func(b *Builder) {
			b.AddList([]string{
				"Avoid optional parameters initially",
				"Add complexity gradually",
			})
		})

		b.AddSection(3, "3. Regular Maintenance", func(b *Builder) {
			b.AddList([]string{
				"Update goals as project progresses",
				"Maintain current rules and guidelines",
				"Log significant changes",
			})
		})

		b.AddSection(3, "4. Error Handling", func(b *Builder) {
			b.AddList([]string{
				"Check for \"Not connected\" errors",
				"Restart Cursor if needed",
				"Use troubleshooting guide",
			})
		})
	})

	// Footer
	b.AddHorizontalRule()
	b.AddParagraph(fmt.Sprintf("**Last Updated**: %s", time.Now().Format("2006-01-02")))
	b.AddParagraph("**Status**: Production Ready")

	return b
}

// MCPTroubleshootingGuideBuilder builds the MCP troubleshooting guide markdown
func MCPTroubleshootingGuideBuilder(mcpServerPath string) *Builder {
	b := NewBuilder()

	// Front matter
	b.AddParagraph("---")
	b.AddParagraph("description: Troubleshooting guide for MCP server tools")
	b.AddParagraph("globs: [\"**/*\"]")
	b.AddParagraph("alwaysApply: false")
	b.AddParagraph("---")
	b.AddLineBreak()

	// Main header
	b.AddHeader(1, "MCP Tools Troubleshooting")

	// Common Issues section
	b.AddSection(2, "🚨 Common Issues", func(b *Builder) {
		b.AddSection(3, "\"Not connected\" Error", func(b *Builder) {
			b.AddCodeBlock("", `{"error":"Not connected"}`)
			b.AddParagraph("**Cause**: MCP server not properly configured in Cursor")
			b.AddParagraph("**Solution**:")
			b.AddList([]string{
				"Check Cursor MCP configuration",
				"Restart Cursor completely",
				"Verify server binary exists",
			})
		})

		b.AddSection(3, "Parameter Type Validation Errors", func(b *Builder) {
			b.AddCodeBlock("", `Error calling tool: Parameter 'active' must be of type null,boolean, got string`)
			b.AddParagraph("**Cause**: Incorrect parameter types")
			b.AddParagraph("**Solution**: Use required parameters only initially")
		})

		b.AddSection(3, "JSON Schema Validation Errors", func(b *Builder) {
			b.AddCodeBlock("", `MCP error 0: validating tool output: type: <invalid reflect.Value> has type "null", want "array"`)
			b.AddParagraph("**Cause**: Empty database state")
			b.AddParagraph("**Solution**: Initialize with test data first")
		})
	})

	// Quick Fixes section
	b.AddSection(2, "🔧 Quick Fixes", func(b *Builder) {
		b.AddSection(3, "1. Restart Everything", func(b *Builder) {
			b.AddCodeBlock("bash", `# Kill any existing processes
pkill -f mcp-server

# Restart Cursor
# (Close and reopen Cursor application)`)
		})

		b.AddSection(3, "2. Initialize Database", func(b *Builder) {
			b.AddCodeBlock("javascript", `// Add test data
mcp_mcp-server-go_goals_add({
  title: "Test Goal"
})

// Verify tools work
mcp_mcp-server-go_goals_list()`)
		})

		b.AddSection(3, "3. Check Configuration", func(b *Builder) {
			b.AddList([]string{
				"Verify MCP server path in Cursor settings",
				"Ensure server binary exists and is executable",
				"Check database file exists: .agent/state.db",
			})
		})
	})

	// Diagnostic Checklist section
	b.AddSection(2, "📋 Diagnostic Checklist", func(b *Builder) {
		b.AddChecklist([]string{
			"Cursor MCP server configured",
			"No multiple server processes running",
			"Database file exists and has content",
			"Server binary is up-to-date",
			"Cursor restarted after configuration",
		}, []bool{false, false, false, false, false})
	})

	// Support section
	b.AddSection(2, "📞 Support", func(b *Builder) {
		b.AddList([]string{
			fmt.Sprintf("**Documentation**: %s/docs/", mcpServerPath),
			fmt.Sprintf("**Server Path**: %s/build/mcp-server", mcpServerPath),
			"**Database**: .agent/state.db",
		})
	})

	// Footer
	b.AddHorizontalRule()
	b.AddParagraph(fmt.Sprintf("**Last Updated**: %s", time.Now().Format("2006-01-02")))
	b.AddParagraph("**Status**: Production Ready")

	return b
}
