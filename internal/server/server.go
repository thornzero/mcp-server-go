package server

import (
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
	// Initialize database
	dbPath := filepath.Join(repoRoot, ".agent", "state.db")
	_ = os.MkdirAll(filepath.Dir(dbPath), 0o755)

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
	)
	if err != nil {
		return nil, err
	}

	server := &Server{db: db, repoRoot: repoRoot}
	server.scanADRs()

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

func (s *Server) scanADRs() {
	adrDir := filepath.Join(s.repoRoot, "ADR")
	entries, err := os.ReadDir(adrDir)
	if err != nil {
		return // ADR directory doesn't exist, skip
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		p := filepath.Join("ADR", e.Name())
		id := strings.ToUpper(strings.TrimSuffix(strings.ReplaceAll(e.Name(), "-", ""), ".md"))
		title := strings.TrimSuffix(strings.TrimPrefix(e.Name(), "0000-"), ".md")

		// Use GORM's FirstOrCreate to insert or update
		adr := models.ADR{
			ID:    id,
			Title: title,
			Path:  p,
		}
		s.db.FirstOrCreate(&adr, models.ADR{ID: id})
	}
}
