# MCP Server Go - Modular Architecture

This project has been refactored into a modular architecture for better maintainability and navigation.

## Project Structure

```
mcp-server-go/
├── cmd/
│   └── mcp-server/          # Main application entry point
│       └── main.go
├── internal/               # Internal packages (not importable by external projects)
│   ├── server/            # Core server logic and database management
│   ├── goals/             # Goal management functionality
│   ├── adrs/              # Architecture Decision Records management
│   ├── ci/                # CI/testing functionality
│   ├── search/            # Repository search functionality
│   ├── state/             # State management and change logging
│   ├── markdown/          # Markdown linting functionality
│   └── templates/         # Template system functionality
├── pkg/
│   └── types/             # Shared type definitions
└── main.go                # Legacy main file (for backward compatibility)
```

## Module Organization

### Core Modules

- **`internal/server`**: Contains the main Server struct, database initialization, and ADR scanning logic
- **`pkg/types`**: All shared type definitions and input/output structs

### Feature Modules

Each feature module follows a consistent pattern:

- **Handler struct**: Contains a reference to the server instance
- **Constructor**: `New[Module]Handler(server)` function
- **Methods**: Individual tool implementations

#### Available Modules

1. **Goals** (`internal/goals`): Goal management (list, add, update)
2. **ADRs** (`internal/adrs`): Architecture Decision Records (list, get)
3. **CI** (`internal/ci`): Continuous Integration (run tests, last failure)
4. **Search** (`internal/search`): Repository search functionality
5. **State** (`internal/state`): Change logging and state management
6. **Markdown** (`internal/markdown`): Markdown linting tools
7. **Templates** (`internal/templates`): Template system (list, register, get, update, delete, apply)

## Benefits of This Structure

1. **Separation of Concerns**: Each module handles a specific domain
2. **Easier Navigation**: Related functionality is grouped together
3. **Better Testing**: Modules can be tested independently
4. **Maintainability**: Changes to one feature don't affect others
5. **Scalability**: Easy to add new features or modify existing ones

## Building

### Build from cmd directory (recommended):
```bash
go build -o mcp-server ./cmd/mcp-server
```

### Build from root (legacy):
```bash
go build -o mcp-server main.go
```

## Adding New Features

To add a new feature:

1. Create a new directory under `internal/`
2. Define handler struct with server reference
3. Implement constructor and methods
4. Add types to `pkg/types` if needed
5. Register tools in `main.go` or `cmd/mcp-server/main.go`

## Migration Notes

- The original `main.go` is preserved for backward compatibility
- All functionality remains the same, just better organized
- Database schema and behavior unchanged
- All MCP tools continue to work as before
