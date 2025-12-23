#!/usr/bin/env bash
# Test script for azd x release flow
# This simulates what the GitHub Actions workflow will do

set -e

VERSION="${1:-0.2.0-test}"
SKIP_BUILD="${SKIP_BUILD:-false}"
SKIP_PACK="${SKIP_PACK:-false}"

echo "🧪 Testing azd x release flow"
echo "Version: $VERSION"
echo ""

# Check if azd is installed
echo "Checking azd installation..."
if ! command -v azd >/dev/null 2>&1; then
    echo "❌ azd not installed. Install from https://aka.ms/azd"
    exit 1
fi
AZD_VERSION=$(azd version 2>&1 || echo "unknown")
echo "✅ azd installed: $AZD_VERSION"

# Check if microsoft.azd.extensions is installed
echo "Checking azd extensions..."
if azd extension list --output json 2>/dev/null | grep -q "microsoft.azd.extensions"; then
    echo "✅ microsoft.azd.extensions installed"
else
    echo "Installing microsoft.azd.extensions..."
    azd extension install microsoft.azd.extensions || true
fi

# Navigate to cli directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

echo ""
echo "=== Step 1: Build Website ==="
if [ -d "../web" ]; then
    cd "../web"
    if [ -f "package.json" ]; then
        echo "Running pnpm install..."
        pnpm install
        echo "Running pnpm build..."
        pnpm run build
        echo "✅ Website built"
    else
        echo "⚠️  No package.json found, skipping website build"
    fi
    cd "../cli"
else
    echo "⚠️  No web directory found"
fi

echo ""
echo "=== Step 2: Build Extension Binaries ==="
if [ "$SKIP_BUILD" != "true" ]; then
    export EXTENSION_ID="jongio.azd.exec"
    export EXTENSION_VERSION="$VERSION"
    
    echo "Building for all platforms..."
    echo "  EXTENSION_ID=$EXTENSION_ID"
    echo "  EXTENSION_VERSION=$EXTENSION_VERSION"
    
    azd x build --all
    
    if [ -d "bin" ]; then
        echo "✅ Built binaries:"
        find bin -type f -name "exec*" | while read -r file; do
            echo "   - ${file#./}"
        done
    else
        echo "❌ No binaries found in bin/"
        exit 1
    fi
else
    echo "⏭️  Skipped (SKIP_BUILD=true)"
fi

echo ""
echo "=== Step 3: Package Extension ==="
if [ "$SKIP_PACK" != "true" ]; then
    echo "Packaging extension..."
    
    # Temporarily update extension.yaml version for testing
    ORIGINAL_VERSION=$(grep '^version:' ../extension.yaml | awk '{print $2}')
    echo "  Original version in extension.yaml: $ORIGINAL_VERSION"
    echo "  Using test version: $VERSION"
    
    # Backup and update
    cp ../extension.yaml ../extension.yaml.bak
    sed -i.tmp "s/^version:.*/version: $VERSION/" ../extension.yaml
    rm -f ../extension.yaml.tmp
    
    # Run pack and restore regardless of success/failure
    if azd x pack; then
        PACK_SUCCESS=true
    else
        PACK_SUCCESS=false
    fi
    
    # Restore original version
    mv ../extension.yaml.bak ../extension.yaml
    echo "  Restored extension.yaml to version $ORIGINAL_VERSION"
    
    if [ "$PACK_SUCCESS" != "true" ]; then
        exit 1
    fi
    
    REGISTRY_PATH="$HOME/.azd/registry/jongio.azd.exec/$VERSION"
    if [ -d "$REGISTRY_PATH" ]; then
        echo "✅ Packaged archives:"
        find "$REGISTRY_PATH" -type f \( -name "*.zip" -o -name "*.tar.gz" \) | while read -r file; do
            echo "   - $(basename "$file")"
        done
        
        echo ""
        echo "Registry location:"
        echo "  $REGISTRY_PATH"
    else
        echo "❌ No packages found at $REGISTRY_PATH"
        echo ""
        echo "Expected location:"
        echo "  $REGISTRY_PATH"
        exit 1
    fi
else
    echo "⏭️  Skipped (SKIP_PACK=true)"
fi

echo ""
echo "=== Test Summary ==="
echo "✅ All steps completed successfully!"
echo ""
echo "Artifacts location:"
echo "  $REGISTRY_PATH"
echo ""
echo "To clean up this test:"
echo "  rm -rf '$REGISTRY_PATH'"
