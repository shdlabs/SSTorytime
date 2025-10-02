#!/bin/bash

# Text2N4L Standalone Server Build Script

echo "Building Text2N4L Standalone Server..."

# Ensure we're in the right directory
cd "$(dirname "$0")"

# Initialize go modules if needed
if [ ! -f "go.mod" ]; then
    echo "Initializing Go modules..."
    go mod init text2n4l-server
fi

# Tidy dependencies
echo "Downloading dependencies..."
go mod tidy

# Build the server
echo "Building server..."
go build -o text2n4l-server .

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "To run the server:"
    echo "  ./text2n4l-server"
    echo ""
    echo "Then open http://localhost:3001 in your browser"
else
    echo "❌ Build failed!"
    exit 1
fi