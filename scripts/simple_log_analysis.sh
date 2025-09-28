#!/bin/bash

# Simple Cursor Log Analysis
# Usage: ./simple_log_analysis.sh [log_file]

LOG_FILE="${1:-/home/thornzero/Repositories/mcp-server-go/docs/logs/vscode-app-1759085846096.log}"

echo "📊 Cursor Log Analysis: $(basename "$LOG_FILE")"
echo "==============================================="

# Basic stats
echo "📈 Basic Statistics:"
echo "  Total lines: $(wc -l < "$LOG_FILE")"
echo "  File size: $(du -h "$LOG_FILE" | cut -f1)"
echo ""

# Error counts
echo "🚨 Error Summary:"
echo "  Internet issues: $(grep -c "ERR_INTERNET_DISCONNECTED" "$LOG_FILE" 2>/dev/null || echo "0")"
echo "  Missing files: $(grep -c "ENOENT\|EntryNotFound" "$LOG_FILE" 2>/dev/null || echo "0")"
echo "  Memory leaks: $(grep -c "listener LEAK" "$LOG_FILE" 2>/dev/null || echo "0")"
echo "  Composer errors: $(grep -c "composer.*Error" "$LOG_FILE" 2>/dev/null || echo "0")"
echo ""

# Show recent errors (last 10)
echo "🔍 Recent Errors (last 10):"
grep -i "error\|warn\|fail" "$LOG_FILE" | tail -10 | sed 's/^/  /'
echo ""

# Show unique missing files
echo "📁 Missing Files:"
grep "ENOENT.*no such file" "$LOG_FILE" | sed 's/.*stat '\''//' | sed 's/'\''.*//' | sort | uniq | sed 's/^/  /'
echo ""

echo "✅ Action Items:"
echo "  1. Clean up missing files"
echo "  2. Restart Cursor to clear memory leaks"
echo "  3. Check internet connection"
echo "  4. Most errors are harmless for development"
