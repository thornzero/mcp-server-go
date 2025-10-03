package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thornzero/mcp-server-go/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Server struct {
	db       *gorm.DB
	repoRoot string
}

func NewServer(repoRoot string) (*Server, error) {
	// Create .agent directory if it doesn't exist
	agentDir := filepath.Join(repoRoot, ".agent")
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .agent directory: %v", err)
	}

	// Initialize database in .agent directory
	dbPath := filepath.Join(agentDir, "state.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&models.Goal{},
		&models.ADR{},
		&models.CIRun{},
		&models.MarkdownTemplate{},
		&models.TemplateVariable{},
		&models.PreferredTool{},
		&models.CursorRule{},
		&models.ChangelogEntry{},
	)
	if err != nil {
		return nil, err
	}

	server := &Server{db: db, repoRoot: repoRoot}
	server.migrateChangelogToDB()

	return server, nil
}

func (s *Server) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s *Server) GetDB() *gorm.DB {
	return s.db
}

func (s *Server) GetRepoRoot() string {
	return s.repoRoot
}

// GetDocsOutputPath returns the configured docs output path
// Priority: 1. Environment variable MCP_DOCS_OUTPUT_PATH, 2. Default "docs"
func (s *Server) GetDocsOutputPath() string {
	if envPath := os.Getenv("MCP_DOCS_OUTPUT_PATH"); envPath != "" {
		// If absolute path, use as-is; if relative, join with repo root
		if filepath.IsAbs(envPath) {
			return envPath
		}
		return filepath.Join(s.repoRoot, envPath)
	}
	// Default fallback
	return filepath.Join(s.repoRoot, "docs")
}

func (s *Server) migrateChangelogToDB() {
	changelogPath := filepath.Join(s.repoRoot, "CHANGELOG_AGENT.md")
	content, err := os.ReadFile(changelogPath)
	if err != nil {
		return // CHANGELOG_AGENT.md doesn't exist, skip migration
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "- ") {
			continue
		}

		// Parse the line format: "- timestamp — summary"
		parts := strings.SplitN(line, " — ", 2)
		if len(parts) != 2 {
			continue
		}

		summary := strings.TrimSpace(parts[1])
		if summary == "" {
			continue
		}

		// Check if this entry already exists in the database
		var existing models.ChangelogEntry
		err := s.db.Where("summary = ?", summary).First(&existing).Error
		if err == nil {
			continue // Entry already exists, skip
		}

		// Create new changelog entry
		entry := models.ChangelogEntry{
			Summary: summary,
		}
		s.db.Create(&entry)
	}
}
