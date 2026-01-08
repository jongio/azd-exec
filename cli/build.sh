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
        
        # Create a safe version of extension ID (replace dots with dashes)
        EXTENSION_ID_SAFE="${EXTENSION_ID:-jongio.azd.exec}"
        EXTENSION_ID_SAFE="${EXTENSION_ID_SAFE//./-}"
        
        # List of OS and architecture combinations
        PLATFORMS=(
            "linux/amd64"
            "linux/arm64"
            "darwin/amd64"
            "darwin/arm64"
            "windows/amd64"
        )
        
        # Loop through platforms and build
        for PLATFORM in "${PLATFORMS[@]}"; do
            OS=$(echo "$PLATFORM" | cut -d'/' -f1)
            ARCH=$(echo "$PLATFORM" | cut -d'/' -f2)
            
            OUTPUT_NAME="bin/${EXTENSION_ID_SAFE}-${OS}-${ARCH}"
            
            if [ "$OS" = "windows" ]; then
                OUTPUT_NAME+='.exe'
            fi
            
            echo "  Building for $OS/$ARCH..."
            
            # Delete the output file if it already exists
            [ -f "$OUTPUT_NAME" ] && rm -f "$OUTPUT_NAME"
            
            # Build for this platform
            GOOS=$OS GOARCH=$ARCH go build -ldflags "${LDFLAGS}" -o "$OUTPUT_NAME" ./src/cmd/script
            
            if [ $? -ne 0 ]; then
                echo "ERROR: Build failed for $OS/$ARCH"
                exit 1
            fi
        done

        echo "All builds complete!"
        echo "Verifying binaries:"
        find bin -type f -name "${EXTENSION_ID_SAFE}-*" | while read -r file; do
            echo "  - ${file} ($(stat -f%z "$file" 2>/dev/null || stat -c%s "$file" 2>/dev/null || echo "?") bytes)"
        done
fi

echo "Build complete!"
