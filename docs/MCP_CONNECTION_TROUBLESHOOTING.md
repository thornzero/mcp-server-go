# MCP Tools Connection Issues - Troubleshooting Guide

## 🔍 **Issue Identified: MCP Server Connection Problems**

Based on the investigation, the MCP server is designed to run as a **stdio-based process**, not a standalone server. This means it communicates via standard input/output and must be started by an MCP client (like Cursor).

## 🚨 **Common Symptoms**

- `{"error":"Not connected"}` when calling MCP tools
- MCP server starts and exits immediately when run directly
- Tools work intermittently or not at all
- Multiple MCP server processes running

## 🛠️ **Root Cause Analysis**

### **MCP Server Architecture**
The MCP server uses `mcp.StdioTransport{}` which means:
- It communicates via stdin/stdout
- It's not a standalone HTTP server
- It must be started by an MCP client (Cursor)
- Running it directly causes immediate exit

### **Connection Issues**
- Cursor may not be properly configured to use the MCP server
- Multiple server instances may be conflicting
- Database state may be corrupted
- Server may not be finding the correct project root

## ✅ **Solutions**

### **Solution 1: Verify Cursor MCP Configuration**

1. **Check Cursor Settings:**
   ```json
   // In Cursor settings, ensure MCP server is configured
   {
     "mcp.servers": {
       "mcp-server-go": {
         "command": "/home/thornzero/Repositories/mcp-server-go/build/mcp-server",
         "args": []
       }
     }
   }
   ```

2. **Restart Cursor** after configuration changes

### **Solution 2: Clean Up Server Processes**

```bash
# Kill any existing MCP server processes
pkill -f mcp-server

# Verify no processes are running
ps aux | grep mcp-server
```

### **Solution 3: Verify Database State**

```bash
# Check if database exists and has content
ls -la /home/thornzero/Repositories/mcp-server-go/.agent/state.db

# Check database size (should be > 0)
du -h /home/thornzero/Repositories/mcp-server-go/.agent/state.db
```

### **Solution 4: Rebuild MCP Server**

```bash
cd /home/thornzero/Repositories/mcp-server-go
go build -o build/mcp-server ./cmd/mcp-server-go
```

### **Solution 5: Test MCP Server Manually**

```bash
# Test if server responds to MCP protocol
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}' | ./build/mcp-server
```

## 🔧 **Prevention Strategies**

### **1. Proper MCP Server Management**
- Never run MCP server directly as a standalone process
- Let Cursor manage the server lifecycle
- Use proper MCP client configuration

### **2. Database Maintenance**
- Regularly check database integrity
- Backup `.agent/state.db` before major changes
- Initialize with test data if database is empty

### **3. Process Monitoring**
- Monitor for multiple MCP server instances
- Clean up orphaned processes regularly
- Use proper process management

## 📋 **Diagnostic Checklist**

### **Before Reporting Issues:**

- [ ] Cursor is properly configured with MCP server
- [ ] No multiple MCP server processes running
- [ ] Database file exists and has content
- [ ] MCP server binary is up-to-date
- [ ] Cursor has been restarted after configuration changes

### **Quick Tests:**

```bash
# 1. Check processes
ps aux | grep mcp-server

# 2. Check database
ls -la .agent/state.db

# 3. Test server binary
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}' | ./build/mcp-server

# 4. Check Cursor logs for MCP errors
```

## 🚀 **Working Configuration**

### **Successful Setup:**
1. **MCP Server Binary**: `/home/thornzero/Repositories/mcp-server-go/build/mcp-server`
2. **Database**: `/home/thornzero/Repositories/mcp-server-go/.agent/state.db`
3. **Project Root**: `/home/thornzero/Repositories/mcp-server-go`
4. **Transport**: stdio (not HTTP)

### **Expected Behavior:**
- MCP tools work consistently
- Database persists between sessions
- No manual server management required
- Cursor handles server lifecycle

## ⚠️ **Common Mistakes**

### **Don't Do:**
- Run MCP server as standalone process
- Start multiple server instances
- Modify database directly
- Ignore connection errors

### **Do:**
- Configure MCP server in Cursor settings
- Let Cursor manage server lifecycle
- Use MCP tools for database operations
- Monitor for connection issues

## 🔄 **Recovery Procedures**

### **If Tools Stop Working:**

1. **Check Connection:**
   ```javascript
   // Try a simple tool call
   mcp_mcp-server-go_goals_list()
   ```

2. **Restart Cursor** if connection fails

3. **Clean Up Processes:**
   ```bash
   pkill -f mcp-server
   ```

4. **Verify Configuration** in Cursor settings

5. **Test Basic Functionality:**
   ```javascript
   mcp_mcp-server-go_goals_add({ title: "Test Goal" })
   mcp_mcp-server-go_goals_list()
   ```

## 📊 **Success Indicators**

✅ **Working Correctly:**
- MCP tools respond without errors
- Database persists between sessions
- No manual server management needed
- Consistent tool behavior

❌ **Not Working:**
- "Not connected" errors
- Tools fail intermittently
- Need to manually start server
- Database issues

---

**Last Updated**: 2025-09-28  
**Status**: Production Ready  
**Next Review**: When connection issues occur
