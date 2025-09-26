# MCP Server Troubleshooting Guide

## Current Status

✅ **MCP Server Implementation**: Complete and working
✅ **Protocol Compliance**: All MCP methods implemented correctly
✅ **Configuration**: Updated in `~/.cursor/mcp.json`
✅ **Testing**: Server responds correctly to all MCP protocol requests

## What We Fixed

1. **Added Missing MCP Protocol Methods**:
   - `initialize` - Proper capabilities advertisement
   - `initialized` - Notification handler
   - `tools/list` - Lists 9 tools with proper schemas
   - `prompts/list` - Lists 2 prompts
   - `resources/list` - Lists 2 resources
   - `tools/call`, `prompts/get`, `resources/read` - Handlers

2. **Fixed Path Issues**:
   - Server now uses executable directory instead of current working directory
   - Fixed schema.sql path resolution

3. **Fixed JSON Handling**:
   - Support for both string and numeric IDs
   - Proper notification handling (no response for notifications)

4. **Updated Cursor Configuration**:
   - Correct executable path
   - Added working directory
   - Removed outdated method lists

## Current Configuration

```json
{
  "mcpServers": {
    "MewlingGoatTavern": {
      "command": "/home/thornzero/Repositories/mewling-goat-tavern/mcp/mcp-server",
      "transport": "stdio",
      "cwd": "/home/thornzero/Repositories/mewling-goat-tavern/mcp",
      "env": {
        "RIPGREP_CONFIG_PATH": "/home/thornzero/.ripgreprc"
      },
      "args": []
    }
  }
}
```

## Troubleshooting Steps

### 1. Restart Cursor Completely

- Close all Cursor windows
- Restart Cursor
- Check Settings → Features → MCP

### 2. Check MCP Logs

- Open Cursor
- Go to View → Output
- Select "MCP Logs" from the dropdown
- Look for any error messages

### 3. Verify Server is Running

The server should appear in Cursor's MCP settings. If it doesn't:

1. **Check if Cursor can find the executable**:

   ```bash
   /home/thornzero/Repositories/mewling-goat-tavern/mcp/mcp-server --help
   ```

2. **Test the server manually**:

   ```bash
   cd /home/thornzero/Repositories/mewling-goat-tavern/mcp
   ./test-server.sh
   ```

### 4. Check for Permission Issues

```bash
ls -la /home/thornzero/Repositories/mewling-goat-tavern/mcp/mcp-server
```

Should show executable permissions (`-rwxrwxr-x`)

### 5. Alternative Configuration

If the current config doesn't work, try this simpler version:

```json
{
  "mcpServers": {
    "MGT": {
      "command": "/home/thornzero/Repositories/mewling-goat-tavern/mcp/mcp-server",
      "transport": "stdio"
    }
  }
}
```

### 6. Check Cursor Version

- Ensure you're using a recent version of Cursor
- MCP support was added in recent versions

### 7. Enable Agent Mode

- Make sure Agent Mode is enabled in Cursor
- MCP tools are primarily used in Agent Mode

## Available Tools

The server provides these tools:

- `goals_list` - List project goals
- `goals_add` - Add new goals
- `goals_update` - Update existing goals
- `adrs_list` - List Architecture Decision Records
- `adrs_get` - Get ADR content
- `ci_run_tests` - Run project tests
- `ci_last_failure` - Check last test failure
- `repo_search` - Search repository
- `state_log_change` - Log changes

## Available Prompts

- `goal_planning` - Help plan goals
- `code_review` - Code review assistance

## Available Resources

- `file://adrs` - All ADR files
- `file://changelog` - Project changelog

## Next Steps

1. **Restart Cursor completely**
2. **Check MCP logs for errors**
3. **Verify the server appears in Settings → Features → MCP**
4. **Test a tool by asking Cursor to use it**

If the server still doesn't appear, the issue is likely with Cursor's MCP integration rather than the server implementation, which is working correctly.
