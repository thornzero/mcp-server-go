#!/bin/bash

# Setup script for Go documentation tools
echo "🔧 Setting up Go documentation tools..."

# Add Go bin to PATH if not already there
export PATH="$(go env GOPATH)/bin:$PATH"

# Install godoc if not available
if ! command -v godoc &> /dev/null; then
    echo "📦 Installing godoc..."
    go install golang.org/x/tools/cmd/godoc@latest
fi

# Install pkgsite if not available
if ! command -v pkgsite &> /dev/null; then
    echo "📦 Installing pkgsite..."
    go install golang.org/x/pkgsite/cmd/pkgsite@latest
fi

echo "✅ Documentation tools ready!"
echo ""
echo "📖 Available commands:"
echo "  make docs        - Start godoc server (http://localhost:6060)"
echo "  make docs-modern - Start pkgsite server (http://localhost:8080)"
echo "  make docs-build  - Build static documentation"
echo ""
echo "🎯 Quick test:"
echo "  godoc -http=:6060 &"
echo "  # Then visit http://localhost:6060/pkg/github.com/thornzero/mcp-server-go/"
