# Mewling Goat Tavern MCP Server (Official SDK)

This is a Model Context Protocol (MCP) server built using the [official Go SDK](https://github.com/modelcontextprotocol/go-sdk) for the Mewling Goat Tavern project.

## Features

### Tools (9 available)

- **goals_list**: List active goals from the project
- **goals_add**: Add a new goal to the project  
- **goals_update**: Update an existing goal
- **adrs_list**: List Architecture Decision Records (ADRs)
- **adrs_get**: Get the content of a specific ADR
- **ci_run_tests**: Run tests for the project
- **ci_last_failure**: Get information about the last test failure
- **repo_search**: Search the repository for text patterns
- **state_log_change**: Log a change to the project changelog

## Implementation

This server is built using the official MCP Go SDK v0.7.0, which provides:

- ✅ **Automatic JSON schema generation** from Go structs
- ✅ **Type-safe tool implementations** with proper input/output validation
- ✅ **Full MCP protocol compliance** with latest version (2025-06-18)
- ✅ **Built-in transport handling** for stdio communication

## Configuration

The server is configured in Cursor's MCP settings (`~/.cursor/mcp.json`):

```json
{
  "mcpServers": {
    "mgt-sdk": {
      "command": "/home/thornzero/Repositories/mewling-goat-tavern/mcp/mcp-server-sdk",
      "transport": "stdio"
    }
  }
}
```

## Building

To build the server:

```bash
cd /home/thornzero/Repositories/mewling-goat-tavern/mcp
go mod tidy
go build -o mcp-server-sdk main.go
```

## Testing

Test the server manually:

```bash
# Test full handshake
cd /home/thornzero/Repositories/mewling-goat-tavern/mcp
./test-sdk-server.sh
```

Expected output shows:

- Proper MCP initialization with protocol version 2025-06-18
- Complete tools list with JSON schemas
- All 9 tools properly registered

## Dependencies

- Go 1.21+
- [github.com/modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) v0.7.0
- modernc.org/sqlite for database operations
- ripgrep (optional, for enhanced search functionality)

## Database

The server uses SQLite for persistence, storing data in `.agent/state.db`. The schema includes:

- `goals`: Project goals with priority, status, and notes
- `adrs`: Architecture Decision Records metadata
- `ci_runs`: Test run history

## Advantages of SDK Implementation

Compared to the manual JSON-RPC implementation, the official SDK provides:

1. **Automatic Schema Generation**: Input/output schemas are generated from Go structs
2. **Type Safety**: Compile-time validation of tool implementations
3. **Protocol Compliance**: Always uses the latest MCP protocol version
4. **Maintainability**: Updates with SDK releases for new MCP features
5. **Reliability**: Used by official MCP ecosystem
