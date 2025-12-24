#!/usr/bin/env bash
# Local build and install script for net-tools
# This script builds the binary from source using idiomatic Go build process

set -e

echo "Building net-tools from source..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Build the binary from cmd/net-tools (idiomatic Go structure)
echo "Building binary from cmd/net-tools..."
go build -o bin/net-tools ./cmd/net-tools/

# Make it executable
chmod +x bin/net-tools

# Install to local bin directory
echo "Installing to ~/.local/bin..."
mkdir -p ~/.local/bin
cp bin/net-tools ~/.local/bin/net-tools

echo "✅ net-tools successfully built and installed!"
echo "Binary location: ~/.local/bin/net-tools"
echo "Make sure ~/.local/bin is in your PATH"

# Test the installation
if command -v net-tools &> /dev/null; then
    echo "✅ net-tools is now available in your PATH"
    net-tools --help
else
    echo "⚠️  net-tools is installed but not in your PATH"
    echo "Run: export PATH=~/.local/bin:\$PATH"
fi