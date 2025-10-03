# MCP Tools Troubleshooting Guide

## 🎯 **Purpose**

This guide ensures MCP server tools work reliably across any project, providing consistent project management capabilities for AI agents.

## 🔍 **Common Issues & Solutions**

### **Issue 1: Parameter Type Validation Errors**

**Symptoms:**
```
Error calling tool: Parameter 'active' must be of type null,boolean, got string
Error calling tool: Parameter 'notes' must be of type null,string, got string
```

**Root Cause:** MCP SDK parameter type handling issues with optional parameters (pointers to primitives).

**Solutions:**

1. **Avoid Optional Parameters Initially:**
   ```javascript
   // ✅ WORKING - Required parameters only
   mcp_project-manager_goals_add({
     title: "Goal Title"
   })
   
   // ❌ PROBLEMATIC - Optional parameters
   mcp_project-manager_goals_add({
     title: "Goal Title",
     notes: "Some notes"  // This causes type validation errors
   })
   ```

2. **Use Required Parameters First:**
   - Start with tools that have only required parameters
   - Add optional parameters after confirming basic functionality

### **Issue 2: JSON Schema Validation Errors**

**Symptoms:**
```
MCP error 0: validating tool output: validating root: validating /properties/rules: type: <invalid reflect.Value> has type "null", want "array"
```

**Root Cause:** MCP server returns `null` instead of empty arrays when no data exists.

**Solutions:**

1. **Initialize Database with Test Data:**
   ```javascript
   // Add a test goal
   mcp_project-manager_goals_add({
     title: "Test Goal for MCP Tools"
   })
   
   // Add a test cursor rule
   mcp_project-manager_cursor_rules_add({
     name: "Test Rule",
     category: "testing", 
     content: `---
description: Test rule for MCP tools validation
globs: ["*.md"]
alwaysApply: true
---

# Test Rule

This is a test rule to validate MCP tools functionality.`
   })
   ```

2. **Verify Database State:**
   ```bash
   # Check if database exists
   ls -la .agent/state.db
   
   # Check database size (should be > 0)
   du -h .agent/state.db
   ```

### **Issue 3: Empty Database State**

**Symptoms:**
- Tools return empty arrays `[]` or validation errors
- No data visible in list operations

**Solutions:**

1. **Bootstrap Process:**
   ```javascript
   // Step 1: Add initial goal
   mcp_project-manager_goals_add({
     title: "Project Setup Complete"
   })
   
   // Step 2: Add initial cursor rule
   mcp_project-manager_cursor_rules_add({
     name: "Project Guidelines",
     category: "general",
     content: `---
description: General project guidelines
globs: ["**/*"]
alwaysApply: true
---

# Project Guidelines

Follow these general guidelines for all development work.`
   })
   
   // Step 3: Verify tools work
   mcp_project-manager_goals_list()
   mcp_project-manager_cursor_rules_list()
   ```

## 🛠️ **Reliable Usage Patterns**

### **Pattern 1: Safe Tool Initialization**

```javascript
// Always start with these tools (no optional parameters)
const safeTools = [
  'mcp_project-manager_goals_add',
  'mcp_project-manager_cursor_rules_add', 
  'mcp_project-manager_goals_list',
  'mcp_project-manager_cursor_rules_list'
]

// Use required parameters only initially
mcp_project-manager_goals_add({
  title: "Project Goal"
})
```

### **Pattern 2: Progressive Enhancement**

```javascript
// Step 1: Basic functionality
mcp_project-manager_goals_add({ title: "Goal 1" })

// Step 2: Add more complex data
mcp_project-manager_cursor_rules_add({
  name: "Rule 1",
  category: "language",
  content: "Rule content here"
})

// Step 3: Test list operations
mcp_project-manager_goals_list()
mcp_project-manager_cursor_rules_list()
```

### **Pattern 3: Error Recovery**

```javascript
// If tools fail, try this sequence:
try {
  // 1. Check if database exists
  // 2. Add minimal test data
  mcp_project-manager_goals_add({ title: "Test" })
  
  // 3. Verify basic functionality
  mcp_project-manager_goals_list()
  
  // 4. Proceed with normal operations
} catch (error) {
  // Log error and provide fallback
  console.log("MCP tools error:", error)
  // Use manual documentation updates as fallback
}
```

## 🔧 **Configuration Requirements**

### **MCP Server Setup:**

1. **Server Running:** Ensure MCP server is running
   ```bash
   ps aux | grep mcp-server
   ```

2. **Database Location:** Check `.agent/state.db` exists
   ```bash
   ls -la .agent/state.db
   ```

3. **Permissions:** Ensure write access to `.agent/` directory
   ```bash
   ls -la .agent/
   ```

### **Project Structure:**

```
project-root/
├── .agent/
│   └── state.db          # MCP server database
├── .cursor/
│   └── rules/            # Cursor rules directory
└── docs/                 # Documentation output
```

## 📋 **Best Practices**

### **1. Always Initialize First**
- Add test data before using list operations
- Verify tools work with simple operations first

### **2. Use Required Parameters Only**
- Avoid optional parameters until basic functionality is confirmed
- Add complexity gradually

### **3. Test After Each Change**
- Verify tools work after adding data
- Test list operations to ensure data persistence

### **4. Have Fallback Plans**
- Manual documentation updates as backup
- File-based goal tracking as alternative

### **5. Monitor Database State**
- Check database file size and modification time
- Ensure `.agent/` directory has proper permissions

## 🚨 **Emergency Procedures**

### **If All Tools Fail:**

1. **Check Server Status:**
   ```bash
   ps aux | grep mcp-server
   ```

2. **Restart MCP Server:**
   ```bash
   cd /path/to/project-manager
   ./build/mcp-server
   ```

3. **Reinitialize Database:**
   ```bash
   rm -rf .agent/state.db
   # Restart server to recreate database
   ```

4. **Use Manual Fallback:**
   - Update documentation files directly
   - Use file-based goal tracking
   - Maintain project state in version control

## 📊 **Success Indicators**

✅ **Tools Working Correctly:**
- `mcp_project-manager_goals_list()` returns array of goals
- `mcp_project-manager_cursor_rules_list()` returns array of rules
- No parameter type validation errors
- Database file exists and has content

❌ **Tools Not Working:**
- Parameter type validation errors
- JSON schema validation errors
- Empty arrays or null responses
- Database file missing or empty

## 🔄 **Maintenance**

### **Regular Checks:**
- Verify MCP server is running
- Check database file integrity
- Test basic tool functionality
- Monitor for parameter type issues

### **Updates:**
- Keep MCP server updated
- Test after any server updates
- Verify compatibility with new Cursor versions

---

**Last Updated:** 2025-09-28  
**Version:** 1.0  
**Status:** Production Ready
