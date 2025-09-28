# MCP Server Table-Driven Tests

This directory contains comprehensive table-driven tests for the MCP Server Go project.

## Test Structure

The tests are organized using Go's table-driven testing pattern with the following structure:

### Test Categories

1. **`TestMCPProtocol`** - Tests MCP protocol compliance
   - Server initialization
   - Tool listing
   - Basic JSON-RPC communication

2. **`TestGoalsTools`** - Tests goal management functionality
   - List goals
   - Add goals (valid and invalid data)
   - Goal validation

3. **`TestADRTools`** - Tests Architecture Decision Records
   - List ADRs
   - Get specific ADR content

4. **`TestCITools`** - Tests CI/testing functionality
   - Run tests
   - Get last test failure

5. **`TestMCPIntegration`** - Full integration test
   - Runs all tool tests in sequence
   - End-to-end workflow testing

### Test Case Structure

Each test category uses a table-driven approach:

```go
type ToolTestCase struct {
    Name        string
    ToolName    string
    Arguments   map[string]interface{}
    ExpectError bool
    Validate    func(t *testing.T, resp *MCPResponse)
}
```

### Key Features

- **Modular Design**: Each tool category has its own test function
- **Reusable Setup**: Common server setup/teardown functions
- **Flexible Validation**: Custom validation functions for each test case
- **Error Testing**: Both success and failure scenarios
- **Clear Naming**: Descriptive test names for easy identification

## Running Tests

### All Integration Tests
```bash
make test-integration
```

### Specific Test Categories
```bash
make test-goals      # Goals management tests
make test-protocol   # MCP protocol tests
make test-ci         # CI/testing tests
```

### Individual Test Functions
```bash
cd test
go test -v -run TestGoalsTools
go test -v -run TestMCPProtocol
go test -v -run TestCITools
```

### Run Specific Test Cases
```bash
cd test
go test -v -run "TestGoalsTools/Add_Goal"
go test -v -run "TestMCPProtocol/Initialize_Server"
```

## Adding New Tests

### 1. Add Test Cases to Existing Categories

```go
func getGoalsToolTests() []ToolTestCase {
    return []ToolTestCase{
        // ... existing tests ...
        {
            Name:      "Update Goal",
            ToolName:  "goals_update",
            Arguments: map[string]interface{}{
                "id":     1,
                "status": "done",
            },
            ExpectError: false,
            Validate: func(t *testing.T, resp *MCPResponse) {
                // Custom validation logic
            },
        },
    }
}
```

### 2. Create New Test Categories

```go
func getNewToolTests() []ToolTestCase {
    return []ToolTestCase{
        {
            Name:      "Test New Tool",
            ToolName:  "new_tool",
            Arguments: map[string]interface{}{},
            ExpectError: false,
            Validate: func(t *testing.T, resp *MCPResponse) {
                // Validation logic
            },
        },
    }
}

func TestNewTools(t *testing.T) {
    cmd, stdin, stdout := setupTestServer(t)
    defer teardownTestServer(cmd)
    
    // Initialize server...
    
    tests := getNewToolTests()
    for _, tt := range tests {
        t.Run(tt.Name, func(t *testing.T) {
            // Test logic
        })
    }
}
```

### 3. Add Makefile Target

```makefile
.PHONY: test-new
test-new:
	@echo "Running new tool tests..."
	cd test && go test -v -run TestNewTools
```

## Test Data Management

- **Server Setup**: Each test function sets up its own server instance
- **Database State**: Tests use the actual database (consider cleanup for production)
- **Isolation**: Each test case is independent
- **Cleanup**: Server processes are properly terminated

## Best Practices

1. **Descriptive Names**: Use clear, descriptive test case names
2. **Error Testing**: Include both success and failure scenarios
3. **Validation**: Write specific validation functions for each test case
4. **Isolation**: Ensure tests don't depend on each other
5. **Cleanup**: Always clean up resources (server processes, etc.)
6. **Logging**: Use `t.Logf()` for informative output during tests

## Debugging Tests

### Verbose Output
```bash
go test -v
```

### Run Single Test
```bash
go test -v -run "TestGoalsTools/Add_Goal"
```

### Debug Mode
```bash
go test -v -run "TestGoalsTools" -args -test.v
```

This table-driven approach makes the tests more maintainable, readable, and easier to extend with new test cases.
