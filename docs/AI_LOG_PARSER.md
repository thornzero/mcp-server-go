# AI Log Parser (MCP Tool)

A Go-based log parser integrated into the MCP server, specifically designed for analyzing Cursor/VS Code logs with AI model optimization in mind.

## Features

- **MCP Integration**: Available as `mcp_mcp-server-go_log_parse()` tool
- **Multiple Output Formats**: JSON, summary, detailed, and AI-friendly formats
- **Comprehensive Error Analysis**: Categorizes errors by type and severity
- **Critical Issue Detection**: Identifies disk space, syntax, and permission errors
- **AI-Optimized Output**: Structured data perfect for AI model consumption
- **Performance**: Handles large log files efficiently with Go's performance

## Usage

The log parser is now integrated into the MCP server as a tool. Use it directly through the MCP interface:

```bash
# Build the MCP server (includes log parser)
make build

# The log parser is now available as an MCP tool
# Use: mcp_mcp-server-go_log_parse()
```

## MCP Tool Usage

### Basic Usage
```javascript
// Parse default log file with AI-friendly output
mcp_mcp-server-go_log_parse()

// Parse specific log file
mcp_mcp-server-go_log_parse({file_path: "/path/to/log.json"})

// Get JSON output for programmatic analysis
mcp_mcp-server-go_log_parse({format: "json"})

// Get summary output
mcp_mcp-server-go_log_parse({format: "summary"})
```

### Parameters
- `file_path` (optional): Path to log file (auto-detects if not provided)
- `format` (optional): Output format - "json", "summary", "detailed", "ai-friendly" (default: "ai-friendly")

## Output Formats

### JSON Format
Structured JSON output perfect for programmatic analysis:
```json
{
  "file": "vscode-app-1759085846096.log",
  "analysis_time": "2025-09-28T18:04:49-04:00",
  "error_counts": {
    "network_issues": 330,
    "file_errors": 4,
    "memory_leaks": 351,
    "mcp_errors": 4
  },
  "critical_issues": {
    "disk_space": 0,
    "syntax_errors": 0,
    "permission_denied": 0
  },
  "recommendations": [
    "🧠 MEDIUM: Restart Cursor to clear memory leaks (351 warnings)",
    "📁 MEDIUM: Clean up missing files (4 errors)"
  ]
}
```

### AI-Friendly Format
Optimized for AI model consumption with clear sections:
- File Statistics
- Error Summary
- Critical Issues
- Detailed Error Analysis
- Recent Errors
- Missing Files
- AI Recommendations
- Context for AI Analysis

### Summary Format
Quick overview with key metrics and status indicators.

### Detailed Format
Comprehensive analysis with full error patterns and recommendations.

## Error Categories

The parser identifies and categorizes these error types:

- **Network Issues**: Internet connectivity failures (usually harmless)
- **File Errors**: Missing files and file system errors
- **Memory Leaks**: Potential listener leaks (restart recommended)
- **MCP Errors**: MCP server-related issues (investigation needed)
- **JavaScript Errors**: Type, reference, and syntax errors
- **Permission Errors**: Access denied and permission issues
- **Connection Errors**: Connection refused and timeout errors

## Critical Issues Detection

Automatically flags critical problems requiring immediate attention:
- **Disk Space**: `ENOSPC` errors indicating full disk
- **Syntax Errors**: JavaScript syntax problems
- **Permission Denied**: File access permission issues

## AI Integration

The AI-friendly format provides:
- Structured error categorization
- Severity levels (low, medium, high, critical)
- Priority-based recommendations
- Context for AI analysis
- Actionable next steps

## Makefile Integration

```bash
# Build the MCP server (includes log parser)
make build

# The log parser is now available as an MCP tool
# No separate executable needed!
```

## Examples

### Quick Status Check
```javascript
mcp_mcp-server-go_log_parse({format: "summary"})
```

### Programmatic Analysis
```javascript
mcp_mcp-server-go_log_parse({format: "json"})
```

### AI-Assisted Debugging
```javascript
mcp_mcp-server-go_log_parse({format: "ai-friendly"})
```

## Benefits of MCP Integration

- **Unified Interface**: All tools available through the same MCP server
- **AI-Optimized**: Designed specifically for AI model consumption
- **No Separate Executables**: Everything integrated into one server
- **Consistent API**: Same parameter patterns as other MCP tools
- **Better Error Handling**: Integrated error reporting through MCP
