# MCP Server Integration Tests

This directory contains integration tests for the MCP Server Go project.

## Running Tests

### From the project root:
```bash
make test-integration
```

### From the test directory:
```bash
cd test && go run main.go
```

## What the Tests Do

The integration test demonstrates how to interact with the MCP server tools:

1. **Initializes** the MCP server
2. **Lists available tools** to verify the server is working
3. **Lists current goals** from the database
4. **Adds a test goal** to demonstrate the goals_add tool
5. **Lists goals again** to show the new goal was added

## Available Tools Tested

- `goals_list` - List active goals
- `goals_add` - Add new goals
- `goals_update` - Update existing goals
- `adrs_list` - List Architecture Decision Records
- `adrs_get` - Get specific ADR content
- `ci_run_tests` - Run tests
- `ci_last_failure` - Get last test failure info
- `repo_search` - Search repository for text patterns
- `state_log_change` - Log changes to changelog
- `markdown_lint` - Lint markdown files
- `template_list` - List markdown templates
- `template_register` - Register new templates
- `template_get` - Get template details
- `template_update` - Update templates
- `template_delete` - Delete templates
- `template_apply` - Apply templates

## How It Works

The test uses Go's `os/exec` package to:
1. Build the MCP server binary
2. Start the server as a subprocess
3. Communicate via JSON-RPC over stdin/stdout
4. Parse responses and display results

This demonstrates how external clients can interact with the MCP server programmatically.
