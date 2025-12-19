#!/usr/bin/env bash
set -euo pipefail

# Build script for azd-exec extension

VERSION="${VERSION:-0.1.0}"
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS="-X main.version=${VERSION} -X main.buildDate=${BUILD_DATE} -X main.gitCommit=${GIT_COMMIT}"

echo "Building azd-exec ${VERSION} (${GIT_COMMIT})"

# Create output directory
mkdir -p bin

# Build for current platform
echo "Building for current platform..."
if ! command -v go >/dev/null 2>&1; then
    echo "Go not found in PATH; skipping build. Install Go to build from source." >&2
    exit 0
fi

go build -ldflags "${LDFLAGS}" -o bin/exec ./src/cmd/script

# Build for multiple platforms if requested
if [ "${BUILD_ALL:-false}" = "true" ]; then
        echo "Building for multiple platforms..."
        # Linux AMD64
        GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o bin/linux-amd64/exec ./src/cmd/script

        # Linux ARM64
        GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o bin/linux-arm64/exec ./src/cmd/script

        # macOS AMD64
        GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o bin/darwin-amd64/exec ./src/cmd/script

        # macOS ARM64 (Apple Silicon)
        GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o bin/darwin-arm64/exec ./src/cmd/script

        # Windows AMD64
        GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o bin/windows-amd64/exec.exe ./src/cmd/script

        echo "All builds complete!"
fi

echo "Build complete!"
