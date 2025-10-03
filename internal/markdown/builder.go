// /home/thornzero/Repositories/project-manager/internal/markdown/builder.go
package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

// AddWrappedList adds a list of items with text wrapping at the specified width
func (b *Builder) AddWrappedList(items []string, width int) *Builder {
	list := &ast.List{
		ListFlags: 0, // 0 = bullet list
	}

	for _, item := range items {
		wrappedLines := WrapText(item, width)
		listItem := &ast.ListItem{}

		for i, line := range wrappedLines {
			paragraph := &ast.Paragraph{}
			if i == 0 {
				// First line gets the bullet point
				paragraph.Literal = []byte(line)
			} else {
				// Subsequent lines are indented with 2 spaces (standard markdown indentation)
				paragraph.Literal = []byte("  " + line)
			}
			ast.AppendChild(listItem, paragraph)
		}

		ast.AppendChild(list, listItem)
	}

	ast.AppendChild(b.doc, list)
	return b
}

func (b *Builder) AddDefinitionList(items map[string]string) *Builder {
	list := &ast.List{
		ListFlags: ast.ListTypeDefinition,
	}

	for key, value := range items {
		// Create term (key) list item
		termItem := &ast.ListItem{
			ListFlags: ast.ListTypeTerm,
		}
		termParagraph := &ast.Paragraph{}
		termParagraph.Literal = []byte(key)
		ast.AppendChild(termItem, termParagraph)
		ast.AppendChild(list, termItem)

		// Create definition (value) list item
		defItem := &ast.ListItem{
			ListFlags: ast.ListTypeDefinition,
		}
		defParagraph := &ast.Paragraph{}
		defParagraph.Literal = []byte(value)
		ast.AppendChild(defItem, defParagraph)
		ast.AppendChild(list, defItem)
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
					// Check if this is a definition list item
					if n.ListFlags&ast.ListTypeDefinition != 0 {
						// This is a definition (value) - indent it
						result.WriteString(fmt.Sprintf(": %s\n", text))
					} else if n.ListFlags&ast.ListTypeTerm != 0 {
						// This is a term (key) - no prefix needed, just the term
						result.WriteString(fmt.Sprintf("%s\n", text))
					} else {
						// Regular list item
						result.WriteString(fmt.Sprintf("- %s\n", text))
					}
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




// WrapText wraps text to fit within the specified width using greedy word-wrap
func WrapText(s string, width int) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	var currentLine strings.Builder
	currentLength := 0

	for _, word := range words {
		wordLength := len(word)

		if currentLength == 0 {
			// First word in line
			currentLine.WriteString(word)
			currentLength = wordLength
		} else {
			// Check if adding this word would exceed the width
			if currentLength+1+wordLength <= width {
				// Add word to current line
				currentLine.WriteString(" " + word)
				currentLength += 1 + wordLength
			} else {
				// Finish current line and start new one
				lines = append(lines, currentLine.String())
				currentLine.Reset()
				currentLine.WriteString(word)
				currentLength = wordLength
			}
		}
	}

	// Add the last line if it has content
	if currentLength > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}
