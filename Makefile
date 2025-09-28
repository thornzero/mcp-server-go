# MCP Server Go Makefile

# Variables
BINARY_NAME=mcp-server
MAIN_PACKAGE=./cmd/mcp-server-go
BUILD_DIR=build
GO_VERSION=1.25.1

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	go clean
	@echo "Clean complete"

# Run tests
.PHONY: test
test:
	@echo "Running all tests..."
	go test -v ./...

# Run all integration tests
.PHONY: test-all
test-all: test-integration test-goals test-protocol test-ci test-markdown test-templates
	@echo "All integration tests completed!"

# Run integration tests
.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	cd test && go test -v

# Run specific test categories
.PHONY: test-goals
test-goals:
	@echo "Running goals tests..."
	cd test && go test -v -run TestGoalsTools

.PHONY: test-protocol
test-protocol:
	@echo "Running MCP protocol tests..."
	cd test && go test -v -run TestMCPProtocol

.PHONY: test-ci
test-ci:
	@echo "Running CI tests..."
	cd test && go test -v -run TestCITools

.PHONY: test-markdown
test-markdown:
	@echo "Running markdown tests..."
	cd test && go test -v -run TestMarkdownTools

.PHONY: test-templates
test-templates:
	@echo "Running template tests..."
	cd test && go test -v -run TestTemplateTools

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linting
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	go mod tidy
	go mod verify

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download

# Run the application
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Install the binary to GOPATH/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to GOPATH/bin..."
	go install ./cmd/mcp-server-go

# Cross-compile for different platforms
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)

.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)

.PHONY: build-all-platforms
build-all-platforms: build-linux build-windows build-darwin

# Development helpers
.PHONY: dev-setup
dev-setup: deps tidy
	@echo "Setting up development environment..."
	@echo "Installing golangci-lint..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development setup complete"

# Check if required tools are installed
.PHONY: check-tools
check-tools:
	@echo "Checking required tools..."
	@command -v go >/dev/null 2>&1 || { echo "Go is not installed"; exit 1; }
	@go version | grep -q "go$(GO_VERSION)" || { echo "Go version $(GO_VERSION) required"; exit 1; }
	@echo "All required tools are installed"

# Generate documentation
.PHONY: docs
docs:
	@echo "Generating documentation..."
	godoc -http=:6060 &
	@echo "Documentation server started at http://localhost:6060"
	@echo "Press Ctrl+C to stop"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build              - Build the application"
	@echo "  build-all-platforms - Cross-compile for Linux, Windows, and macOS"
	@echo "  clean              - Clean build artifacts"
	@echo "  test               - Run tests"
	@echo "  test-all           - Run all integration tests"
	@echo "  test-integration   - Run integration tests with MCP tools"
	@echo "  test-goals         - Run goals-specific tests"
	@echo "  test-protocol      - Run MCP protocol tests"
	@echo "  test-ci            - Run CI-specific tests"
	@echo "  test-markdown      - Run markdown-specific tests"
	@echo "  test-templates     - Run template-specific tests"
	@echo "  test-coverage      - Run tests with coverage report"
	@echo "  lint               - Run linter"
	@echo "  fmt                - Format code"
	@echo "  tidy               - Tidy dependencies"
	@echo "  deps               - Install dependencies"
	@echo "  run                - Build and run the application"
	@echo "  install            - Install binary to GOPATH/bin"
	@echo "  dev-setup          - Set up development environment"
	@echo "  check-tools        - Check if required tools are installed"
	@echo "  docs               - Start documentation server"
	@echo "  help               - Show this help message"
