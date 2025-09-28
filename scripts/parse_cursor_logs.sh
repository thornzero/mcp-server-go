#!/bin/bash

# Cursor Log Parser
# Usage: ./parse_cursor_logs.sh [log_file]

LOG_FILE="${1:-/home/thornzero/Repositories/mcp-server-go/docs/logs/vscode-app-1759085846096.log}"

if [[ ! -f "$LOG_FILE" ]]; then
    echo "Error: Log file not found: $LOG_FILE"
    exit 1
fi

echo "🔍 Analyzing Cursor Console Log: $(basename "$LOG_FILE")"
echo "📊 Total lines: $(wc -l < "$LOG_FILE")"
echo ""

# Function to count and show errors
analyze_errors() {
    local pattern="$1"
    local description="$2"
    local count=$(grep -c "$pattern" "$LOG_FILE" 2>/dev/null || echo "0")
    
    if [[ $count -gt 0 ]]; then
        echo "❌ $description: $count occurrences"
        if [[ $count -le 10 ]]; then
            echo "   Recent occurrences:"
            grep "$pattern" "$LOG_FILE" | tail -5 | sed 's/^/   /'
        fi
        echo ""
    fi
}

# Function to show unique errors
show_unique_errors() {
    local pattern="$1"
    local description="$2"
    echo "🔍 $description:"
    grep "$pattern" "$LOG_FILE" | sort | uniq -c | sort -nr | head -10 | sed 's/^/   /'
    echo ""
}

echo "🚨 ERROR ANALYSIS"
echo "=================="

# Network/Connectivity Issues
analyze_errors "ERR_INTERNET_DISCONNECTED" "Internet connectivity failures"

# File System Errors
analyze_errors "ENOENT.*no such file" "Missing file errors"
analyze_errors "EntryNotFound.*FileSystemError" "File system errors"

# Memory/Performance Issues
analyze_errors "potential listener LEAK" "Memory leak warnings"

# MCP/Server Specific Issues
analyze_errors "mcp.*error\|MCP.*error" "MCP-related errors"
analyze_errors "server.*error\|Server.*error" "Server-related errors"

# Composer Issues
analyze_errors "composer.*Error" "Composer context errors"

echo "📋 SUMMARY"
echo "=========="

# Count different error types
echo "Error Summary:"
echo "  🌐 Network issues: $(grep -c "ERR_INTERNET_DISCONNECTED" "$LOG_FILE" 2>/dev/null || echo "0")"
echo "  📁 File errors: $(grep -c "ENOENT\|EntryNotFound" "$LOG_FILE" 2>/dev/null || echo "0")"
echo "  🧠 Memory leaks: $(grep -c "listener LEAK" "$LOG_FILE" 2>/dev/null || echo "0")"
echo "  🎭 Composer errors: $(grep -c "composer.*Error" "$LOG_FILE" 2>/dev/null || echo "0")"
echo "  🔧 MCP errors: $(grep -c "mcp.*error\|MCP.*error" "$LOG_FILE" 2>/dev/null || echo "0")"

echo ""
echo "✅ RECOMMENDATIONS"
echo "=================="

# Check for specific issues and provide recommendations
if grep -q "ERR_INTERNET_DISCONNECTED" "$LOG_FILE" 2>/dev/null; then
    echo "  🌐 Internet connectivity issues detected"
    echo "     → Check your internet connection"
    echo "     → These are usually harmless for local development"
    echo ""
fi

if grep -q "ENOENT.*no such file" "$LOG_FILE" 2>/dev/null; then
    echo "  📁 Missing file errors detected"
    echo "     → Clean up temporary test files"
    echo "     → These errors should stop appearing after cleanup"
    echo ""
fi

if grep -q "listener LEAK" "$LOG_FILE" 2>/dev/null; then
    echo "  🧠 Memory leak warnings detected"
    echo "     → These are usually Cursor/VS Code internal issues"
    echo "     → Consider restarting Cursor if performance degrades"
    echo ""
fi

if grep -q "composer.*Error" "$LOG_FILE" 2>/dev/null; then
    echo "  🎭 Composer context errors detected"
    echo "     → These are related to Cursor's AI context gathering"
    echo "     → Usually harmless, but indicate missing files"
    echo ""
fi

echo "🎯 QUICK FIXES"
echo "=============="
echo "  • Clean up temporary files: rm -f test-*.md test_*.go debug_*.go"
echo "  • Restart Cursor to clear memory leaks"
echo "  • Check internet connection for connectivity issues"
echo "  • Most errors are harmless for local development"

echo ""
echo "📝 To analyze a different log file:"
echo "   ./parse_cursor_logs.sh /path/to/other/log/file"
