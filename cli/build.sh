#!/usr/bin/env bash
set -euo pipefail

# Build script called by azd x build
# This script is invoked by `azd x build` with specific environment variables:
# - EXTENSION_ID: The extension identifier (e.g., jongio.azd.exec)
# - EXTENSION_VERSION: The extension version (e.g., 0.1.0)
# - GOOS: Target operating system
# - GOARCH: Target architecture
# - OUTPUT_PATH: Where to write the binary

# Get environment variables set by azd x build
EXTENSION_ID="${EXTENSION_ID:-}"
EXTENSION_VERSION="${EXTENSION_VERSION:-0.1.0}"
TARGET_OS="${GOOS:-}"
TARGET_ARCH="${GOARCH:-}"
OUTPUT_PATH="${OUTPUT_PATH:-}"

if [ -z "$EXTENSION_ID" ]; then
    echo "ERROR: EXTENSION_ID environment variable not set" >&2
    exit 1
fi

if [ -z "$EXTENSION_VERSION" ]; then
    echo "ERROR: EXTENSION_VERSION environment variable not set" >&2
    exit 1
fi

# Build metadata
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags with version info
VERSION_IMPORT_PATH="github.com/jongio/azd-exec/cli/src/internal/version"
LDFLAGS="-X ${VERSION_IMPORT_PATH}.Version=${EXTENSION_VERSION} -X ${VERSION_IMPORT_PATH}.BuildDate=${BUILD_DATE} -X ${VERSION_IMPORT_PATH}.GitCommit=${GIT_COMMIT}"

echo "Building $EXTENSION_ID v$EXTENSION_VERSION"

# If OUTPUT_PATH is set, this is a targeted build for azd x build
if [ -n "$OUTPUT_PATH" ]; then
    # Default OS/Arch if not explicitly set by azd x build
    if [ -z "$TARGET_OS" ]; then
        TARGET_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
        case "$TARGET_OS" in
            darwin) TARGET_OS="darwin" ;;
            linux) TARGET_OS="linux" ;;
            *) TARGET_OS="linux" ;;
        esac
    fi
    if [ -z "$TARGET_ARCH" ]; then
        TARGET_ARCH=$(uname -m)
        case "$TARGET_ARCH" in
            x86_64) TARGET_ARCH="amd64" ;;
            aarch64|arm64) TARGET_ARCH="arm64" ;;
            *) TARGET_ARCH="amd64" ;;
        esac
    fi
    
    echo "  OS/Arch: $TARGET_OS/$TARGET_ARCH"
    echo "  Output: $OUTPUT_PATH"
    
    # IMPORTANT: azd x pack expects binaries with platform-specific names
    # Build DIRECTLY to the platform-specific name that pack expects
    EXTENSION_ID_SAFE="${EXTENSION_ID//./-}"
    BINARY_EXT=""
    if [ "$TARGET_OS" = "windows" ]; then
        BINARY_EXT=".exe"
    fi
    PLATFORM_SPECIFIC_NAME="${EXTENSION_ID_SAFE}-${TARGET_OS}-${TARGET_ARCH}${BINARY_EXT}"
    
    BIN_DIR=$(dirname "$OUTPUT_PATH")
    PLATFORM_SPECIFIC_PATH="${BIN_DIR}/${PLATFORM_SPECIFIC_NAME}"
    
    echo "  Building to: $PLATFORM_SPECIFIC_PATH"
    
    GOOS=$TARGET_OS GOARCH=$TARGET_ARCH go build -ldflags "$LDFLAGS" -o "$PLATFORM_SPECIFIC_PATH" ./src/cmd/script
    
    if [ $? -ne 0 ]; then
        echo "ERROR: Build failed" >&2
        exit $?
    fi
    
    echo "✅ Build successful!"
else
    # Fallback: build for current platform to bin/
    echo "  Building for current platform..."
    
    mkdir -p bin
    
    BINARY_NAME="exec"
    if [ "${GOOS:-}" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
    fi
    
    OUTPUT_PATH="bin/${BINARY_NAME}"
    
    go build -ldflags "$LDFLAGS" -o "$OUTPUT_PATH" ./src/cmd/script
    
    if [ $? -ne 0 ]; then
        echo "ERROR: Build failed" >&2
        exit $?
    fi
    
    echo "✅ Build successful: $OUTPUT_PATH"
fi
