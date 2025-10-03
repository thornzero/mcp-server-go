package main

import (
	"context"
	"fmt"
	"log"

	"github.com/thornzero/project-manager/internal/server"
	"github.com/thornzero/project-manager/internal/setup"
	"github.com/thornzero/project-manager/internal/types"
)

func main() {
	// Initialize server
	srv, err := server.NewServer("/home/thornzero/Repositories/project-manager")
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Close()

	// Create setup handler
	setupHandler := setup.NewSetupHandler(srv)

	// Test setup_mcp_tools
	input := types.SetupMCPToolsInput{
		ProjectPath: "/tmp/test-mcp-setup",
	}

	result, output, err := setupHandler.SetupMCPTools(context.Background(), nil, input)
	if err != nil {
		log.Fatal("Setup failed:", err)
	}

	_ = result // Ignore unused result
	fmt.Printf("Setup result: %+v\n", output)
	fmt.Printf("Files created: %v\n", output.FilesCreated)
}
