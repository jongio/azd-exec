#!/usr/bin/env bash
set -euo pipefail

# Build script for azd-exec extension

VERSION="${VERSION:-0.1.0}"
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

VERSION_IMPORT_PATH="github.com/jongio/azd-exec/cli/src/internal/version"
LDFLAGS="-X ${VERSION_IMPORT_PATH}.Version=${VERSION} -X ${VERSION_IMPORT_PATH}.BuildDate=${BUILD_DATE} -X ${VERSION_IMPORT_PATH}.GitCommit=${GIT_COMMIT}"

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
        
        # Create platform-specific directories
        mkdir -p bin/linux-amd64 bin/linux-arm64 bin/darwin-amd64 bin/darwin-arm64 bin/windows-amd64
        
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
        echo "Verifying binaries:"
        find bin -type f -name "exec*" | while read -r file; do
            echo "  - ${file} ($(stat -f%z "$file" 2>/dev/null || stat -c%s "$file" 2>/dev/null || echo "?") bytes)"
        done
fi

echo "Build complete!"
