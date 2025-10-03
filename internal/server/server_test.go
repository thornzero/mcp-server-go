package server

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewServer(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	tests := []struct {
		name      string
		repoRoot  string
		wantError bool
	}{
		{
			name:      "Valid repository root",
			repoRoot:  tempDir,
			wantError: false,
		},
		{
			name:      "Non-existent directory",
			repoRoot:  "/non/existent/path",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewServer(tt.repoRoot)

			if tt.wantError {
				if err == nil {
					t.Errorf("NewServer() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("NewServer() unexpected error: %v", err)
				return
			}

			if server == nil {
				t.Errorf("NewServer() returned nil server")
				return
			}

			// Verify .agent directory was created
			agentDir := filepath.Join(tt.repoRoot, ".agent")
			if _, err := os.Stat(agentDir); os.IsNotExist(err) {
				t.Errorf("NewServer() did not create .agent directory")
			}

			// Verify database file was created
			dbPath := filepath.Join(agentDir, "state.db")
			if _, err := os.Stat(dbPath); os.IsNotExist(err) {
				t.Errorf("NewServer() did not create state.db file")
			}

			// Clean up
			server.Close()
		})
	}
}

func TestServer_GetRepoRoot(t *testing.T) {
	tempDir := t.TempDir()
	server, err := NewServer(tempDir)
	if err != nil {
		t.Fatalf("NewServer() error: %v", err)
	}
	defer server.Close()

	if server.GetRepoRoot() != tempDir {
		t.Errorf("GetRepoRoot() = %v, want %v", server.GetRepoRoot(), tempDir)
	}
}

func TestServer_GetDocsOutputPath(t *testing.T) {
	tempDir := t.TempDir()
	server, err := NewServer(tempDir)
	if err != nil {
		t.Fatalf("NewServer() error: %v", err)
	}
	defer server.Close()

	tests := []struct {
		name     string
		envPath  string
		expected string
	}{
		{
			name:     "Default docs path",
			envPath:  "",
			expected: filepath.Join(tempDir, "docs"),
		},
		{
			name:     "Custom relative path",
			envPath:  "custom-docs",
			expected: filepath.Join(tempDir, "custom-docs"),
		},
		{
			name:     "Custom absolute path",
			envPath:  "/absolute/path",
			expected: "/absolute/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.envPath != "" {
				os.Setenv("MCP_DOCS_OUTPUT_PATH", tt.envPath)
				defer os.Unsetenv("MCP_DOCS_OUTPUT_PATH")
			}

			result := server.GetDocsOutputPath()
			if result != tt.expected {
				t.Errorf("GetDocsOutputPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestServer_Close(t *testing.T) {
	tempDir := t.TempDir()
	server, err := NewServer(tempDir)
	if err != nil {
		t.Fatalf("NewServer() error: %v", err)
	}

	// Close should not return an error
	if err := server.Close(); err != nil {
		t.Errorf("Close() error: %v", err)
	}
}
