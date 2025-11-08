#!/bin/bash
set -e

# Build script for azd-script extension

VERSION="${VERSION:-0.1.0}"
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS="-X main.version=${VERSION} -X main.buildDate=${BUILD_DATE} -X main.gitCommit=${GIT_COMMIT}"

echo "Building azd-script ${VERSION} (${GIT_COMMIT})"

# Create output directory
mkdir -p bin

# Build for current platform
echo "Building for current platform..."
go build -ldflags "${LDFLAGS}" -o bin/script ./src/cmd/script

# Build for multiple platforms if requested
if [ "$BUILD_ALL" = "true" ]; then
    echo "Building for multiple platforms..."
    
    # Linux AMD64
    GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o bin/linux-amd64/script ./src/cmd/script
    
    # Linux ARM64
    GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o bin/linux-arm64/script ./src/cmd/script
    
    # macOS AMD64
    GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o bin/darwin-amd64/script ./src/cmd/script
    
    # macOS ARM64 (Apple Silicon)
    GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o bin/darwin-arm64/script ./src/cmd/script
    
    # Windows AMD64
    GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o bin/windows-amd64/script.exe ./src/cmd/script
    
    echo "All builds complete!"
fi

echo "Build complete!"
