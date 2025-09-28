#!/bin/bash
# /home/thornzero/Repositories/mcp-server-go/scripts/install-mcp-tools.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
MCP_SERVER_PATH="/home/thornzero/Repositories/mcp-server-go"
SETUP_TOOL="$MCP_SERVER_PATH/build/setup-mcp-tools"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS] [PROJECT_PATH]"
    echo ""
    echo "Install MCP tools for project management and AI assistance"
    echo ""
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  -f, --force    Force installation even if rules exist"
    echo "  -u, --update   Update existing MCP tools installation"
    echo "  -c, --config   Show Cursor configuration needed"
    echo ""
    echo "Arguments:"
    echo "  PROJECT_PATH   Path to project directory (default: current directory)"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Install in current directory"
    echo "  $0 /path/to/my/project               # Install in specific project"
    echo "  $0 --update /path/to/my/project      # Update existing installation"
    echo "  $0 --config                          # Show Cursor configuration"
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check if MCP server exists
    if [ ! -f "$MCP_SERVER_PATH/build/mcp-server" ]; then
        print_error "MCP server not found at $MCP_SERVER_PATH/build/mcp-server"
        print_status "Building MCP server..."
        cd "$MCP_SERVER_PATH"
        make build
        if [ ! -f "$MCP_SERVER_PATH/build/mcp-server" ]; then
            print_error "Failed to build MCP server"
            exit 1
        fi
    fi
    
    # Check if setup tool exists
    if [ ! -f "$SETUP_TOOL" ]; then
        print_status "Building setup tool..."
        cd "$MCP_SERVER_PATH"
        make build-setup
        if [ ! -f "$SETUP_TOOL" ]; then
            print_error "Failed to build setup tool"
            exit 1
        fi
    fi
    
    print_success "Prerequisites check complete"
}

# Function to show Cursor configuration
show_cursor_config() {
    echo ""
    echo "=========================================="
    echo "Cursor MCP Configuration Required"
    echo "=========================================="
    echo ""
    echo "Add this to your Cursor settings (settings.json):"
    echo ""
    echo '{'
    echo '  "mcp.servers": {'
    echo '    "mcp-server-go": {'
    echo "      \"command\": \"$MCP_SERVER_PATH/build/mcp-server\","
    echo '      "args": []'
    echo '    }'
    echo '  }'
    echo '}'
    echo ""
    echo "Or add to your .cursor/settings.json file:"
    echo ""
    echo '{'
    echo '  "mcp.servers": {'
    echo '    "mcp-server-go": {'
    echo "      \"command\": \"$MCP_SERVER_PATH/build/mcp-server\","
    echo '      "args": []'
    echo '    }'
    echo '  }'
    echo '}'
    echo ""
    echo "After adding this configuration:"
    echo "1. Restart Cursor completely"
    echo "2. Test with: mcp_mcp-server-go_goals_list()"
    echo "3. Initialize with: mcp_mcp-server-go_goals_add({title: 'Test Goal'})"
    echo ""
}

# Function to install MCP tools
install_mcp_tools() {
    local project_path="$1"
    local force="$2"
    
    print_status "Installing MCP tools in: $project_path"
    
    # Check if project path exists
    if [ ! -d "$project_path" ]; then
        print_error "Project path does not exist: $project_path"
        exit 1
    fi
    
    # Check if rules already exist
    local rules_dir="$project_path/.cursor/rules"
    if [ -d "$rules_dir" ] && [ "$(ls -A "$rules_dir" 2>/dev/null)" ]; then
        if [ "$force" != "true" ]; then
            print_warning "Cursor rules already exist in $rules_dir"
            echo "Use --force to overwrite existing rules"
            exit 1
        else
            print_warning "Overwriting existing rules..."
        fi
    fi
    
    # Run setup tool
    print_status "Running MCP setup tool..."
    if "$SETUP_TOOL" "$project_path"; then
        print_success "MCP tools installed successfully!"
    else
        print_error "Failed to install MCP tools"
        exit 1
    fi
    
    # Show next steps
    echo ""
    echo "=========================================="
    echo "Installation Complete!"
    echo "=========================================="
    echo ""
    echo "Next steps:"
    echo "1. Configure Cursor MCP server (see --config option)"
    echo "2. Restart Cursor completely"
    echo "3. Test the tools:"
    echo "   mcp_mcp-server-go_goals_list()"
    echo "4. Initialize with test data:"
    echo "   mcp_mcp-server-go_goals_add({title: 'Project Goal'})"
    echo ""
    echo "Available tools:"
    echo "- Goals management"
    echo "- Cursor rules management"
    echo "- Documentation generation"
    echo "- Repository search"
    echo "- CI integration"
    echo ""
    echo "For help: Check .cursor/rules/mcp-tools-usage.mdc"
    echo ""
}

# Function to update MCP tools
update_mcp_tools() {
    local project_path="$1"
    
    print_status "Updating MCP tools in: $project_path"
    
    # Force update
    install_mcp_tools "$project_path" "true"
    
    print_success "MCP tools updated successfully!"
}

# Main script logic
main() {
    local project_path="."
    local force="false"
    local update="false"
    local show_config="false"
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_usage
                exit 0
                ;;
            -f|--force)
                force="true"
                shift
                ;;
            -u|--update)
                update="true"
                shift
                ;;
            -c|--config)
                show_config="true"
                shift
                ;;
            -*)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
            *)
                project_path="$1"
                shift
                ;;
        esac
    done
    
    # Convert relative path to absolute
    project_path=$(realpath "$project_path")
    
    # Show configuration if requested
    if [ "$show_config" = "true" ]; then
        show_cursor_config
        exit 0
    fi
    
    # Check prerequisites
    check_prerequisites
    
    # Install or update MCP tools
    if [ "$update" = "true" ]; then
        update_mcp_tools "$project_path"
    else
        install_mcp_tools "$project_path" "$force"
    fi
}

# Run main function with all arguments
main "$@"
