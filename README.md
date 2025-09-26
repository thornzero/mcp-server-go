# MCP Server Go

A comprehensive Model Context Protocol (MCP) server implementation in Go that provides project management, documentation, and markdown tools for any project.

## Features

### 🎯 **Project Management Tools**
- **Goals Management**: Create, update, list, and track project goals
- **Architecture Decision Records (ADRs)**: Manage ADR documents with automatic discovery
- **Change Logging**: Track project changes with timestamps and file references

### 🔍 **Development Tools**
- **Repository Search**: Search codebase using ripgrep/grep with line numbers
- **CI Integration**: Run tests and check last failure status
- **Markdown Linting**: Validate and auto-fix markdown formatting issues

### 📝 **Template System**
- **Template Management**: Register, update, delete markdown templates
- **Variable Support**: Typed variables (string, date, list, number, boolean)
- **Content Generation**: Apply templates to generate standardized documents

## Installation

### Prerequisites
- Go 1.25+ 
- SQLite (included via modernc.org/sqlite)
- markdownlint-cli (optional, for markdown linting)

### Build from Source
```bash
git clone https://github.com/thornzero/mcp-server-go.git
cd mcp-server-go
go build -o mcp-server main.go
```

### Install markdownlint-cli (optional)
```bash
npm install -g markdownlint-cli
```

## Usage

### Cursor Integration

Add to your `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "mcp-server-go": {
      "command": "/path/to/mcp-server-go/mcp-server",
      "transport": "stdio"
    }
  }
}
```

### Available Tools

#### Project Management
- `goals_list` - List active project goals
- `goals_add` - Add new project goals
- `goals_update` - Update existing goals
- `adrs_list` - List Architecture Decision Records
- `adrs_get` - Get ADR content by ID
- `state_log_change` - Log project changes

#### Development
- `repo_search` - Search repository for text patterns
- `ci_run_tests` - Run project tests
- `ci_last_failure` - Get last test failure information
- `markdown_lint` - Lint markdown files for formatting issues

#### Templates
- `template_list` - List available markdown templates
- `template_register` - Register new templates
- `template_get` - Get template details
- `template_update` - Update existing templates
- `template_delete` - Delete templates
- `template_apply` - Apply templates to generate content

## Configuration

### Markdown Linting

Create a `.markdownlint.json` file in your project root:

```json
{
  "MD013": {
    "line_length": 120,
    "code_blocks": false,
    "tables": false
  },
  "MD022": true,
  "MD032": true,
  "MD029": {
    "style": "ordered"
  },
  "MD036": false,
  "MD041": false
}
```

### Database

The server automatically creates a SQLite database at `.agent/state.db` with the following tables:
- `goals` - Project goals and tasks
- `adrs` - Architecture Decision Records
- `ci_runs` - CI test run history
- `markdown_templates` - Template definitions
- `template_variables` - Template variable definitions

## Template System

### Creating Templates

```bash
# Register a README template
mcp-server template_register \
  --id "README" \
  --name "Project README Template" \
  --category "documentation" \
  --content "# {{.ProjectName}}\n\n{{.Description}}" \
  --variables '[
    {"name": "ProjectName", "type": "string", "required": true},
    {"name": "Description", "type": "string", "required": true}
  ]'
```

### Applying Templates

```bash
# Generate README from template
mcp-server template_apply \
  --template_id "README" \
  --variables '{"ProjectName": "My Project", "Description": "A great project"}' \
  --output_path "README.md"
```

## Development

### Project Structure
```
mcp-server-go/
├── main.go              # Main server implementation
├── schema.sql           # Database schema
├── go.mod              # Go module definition
├── README.md           # This file
└── docs/               # Documentation
    ├── MARKDOWN_TOOLS_RESEARCH.md
    └── TROUBLESHOOTING.md
```

### Adding New Tools

1. Define input/output structs
2. Implement handler function
3. Register tool in main.go
4. Add to database schema if needed

### Testing

```bash
# Test server directly
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {}}' | ./mcp-server

# Test specific tool
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "goals_list", "arguments": {}}}' | ./mcp-server
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Acknowledgments

- [Model Context Protocol](https://modelcontextprotocol.io/) for the MCP specification
- [Go SDK for MCP](https://github.com/modelcontextprotocol/go-sdk) for the Go implementation
- [markdownlint](https://github.com/DavidAnson/markdownlint) for markdown validation rules