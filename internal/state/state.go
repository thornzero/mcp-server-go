package state

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thornzero/project-manager/internal/markdown"
	"github.com/thornzero/project-manager/internal/models"
	"github.com/thornzero/project-manager/internal/server"
	"github.com/thornzero/project-manager/internal/types"
)

type StateHandler struct {
	server *server.Server
}

func NewStateHandler(s *server.Server) *StateHandler {
	return &StateHandler{server: s}
}

func (h *StateHandler) StateLogChange(ctx context.Context, req *mcp.CallToolRequest, input types.StateLogChangeInput) (*mcp.CallToolResult, types.StateLogChangeOutput, error) {
	if strings.TrimSpace(input.Summary) == "" {
		return nil, types.StateLogChangeOutput{}, fmt.Errorf("summary required")
	}

	// Create changelog entry in database
	entry := models.ChangelogEntry{
		Summary: input.Summary,
		Files:   strings.Join(input.Files, ", "),
	}

	err := h.server.GetDB().Create(&entry).Error
	if err != nil {
		return nil, types.StateLogChangeOutput{}, err
	}

	return nil, types.StateLogChangeOutput{OK: true}, nil
}

func (h *StateHandler) ChangelogGenerate(ctx context.Context, req *mcp.CallToolRequest, input types.ChangelogGenerateInput) (*mcp.CallToolResult, types.ChangelogGenerateOutput, error) {
	var entries []models.ChangelogEntry
	query := h.server.GetDB().Order("created_at DESC")

	if input.Limit > 0 {
		query = query.Limit(input.Limit)
	}

	err := query.Find(&entries).Error
	if err != nil {
		return nil, types.ChangelogGenerateOutput{}, err
	}

	if len(entries) == 0 {
		return nil, types.ChangelogGenerateOutput{}, fmt.Errorf("no changelog entries found")
	}

	var content string
	var path string

	if input.Format == "json" {
		// Generate JSON format with proper line breaks
		if entries == nil {
			return nil, types.ChangelogGenerateOutput{}, fmt.Errorf("no changelog entries found")
		}
		contentMap := map[string]any{"changelog": entries}
		contentBytes, err := json.Marshal(contentMap)
		if err != nil {
			return nil, types.ChangelogGenerateOutput{}, err
		}
		content = string(contentBytes[:])
		path = filepath.Join(h.server.GetRepoRoot(), "CHANGELOG.json")
	} else {
		md := markdown.NewBuilder()
		md.AddHeader(1, "Changelog")
		md.AddParagraph("All notable changes to this project are documented in this file.")
		for _, entry := range entries {
			md.AddHeader(2, entry.CreatedAt.Format("2006-01-02"))
			md.AddParagraph(entry.Summary)
			if entry.Files != "" {
				md.AddHeader(3, "Files")
				md.AddWrappedList(strings.Split(entry.Files, ", "), 80)
			}
		}
		content = md.String()
		path = filepath.Join(h.server.GetRepoRoot(), "CHANGELOG.md")
	}

	// Write to file
	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return nil, types.ChangelogGenerateOutput{}, err
	}

	return nil, types.ChangelogGenerateOutput{
		Content: content,
		Path:    path,
		Entries: len(entries),
	}, nil
}

// wrapText wraps text to fit within the specified width
func (h *StateHandler) wrapText(text string, width int) []string {
	words := strings.Fields(text)
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
