#!/bin/bash

# Build script for Go applications
# Creates optimized binaries for multiple platforms
# Usage: ./build.sh [base-filename]
# Example: ./build.sh my-app

# Set base filename from argument or use default
BASE_NAME="${1:-bubblegum-physics-sim}"

echo "Building ${BASE_NAME} for multiple platforms..."

# Create build directory
mkdir -p builds

# Build for Linux (64-bit)
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/${BASE_NAME}-linux .

# Build for Windows (64-bit)
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o builds/${BASE_NAME}.exe .

# Build for macOS Intel (64-bit)
echo "Building for macOS Intel (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o builds/${BASE_NAME}-macos-intel .

# Build for macOS Apple Silicon (arm64)
echo "Building for macOS Apple Silicon (arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o builds/${BASE_NAME}-macos-arm .

# Build for current platform (local testing)
echo "Building for current platform..."
go build -ldflags="-s -w" -o builds/${BASE_NAME}-local .

echo "Build complete! Binaries are in the 'builds' directory:"
ls -la builds/ 