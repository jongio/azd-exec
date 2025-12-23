#!/usr/bin/env bash
set -euo pipefail

# pack.sh - Manual packaging script (for local testing)
# For releases, use: azd x pack
# Usage: ./pack.sh [version]

VERSION="${1:-0.1.0}"

echo "🎁 Packaging azd exec extension version ${VERSION}"
echo ""
echo "Note: This is for local testing. The release workflow uses 'azd x pack' automatically."
echo ""

# Check if azd is available
if ! command -v azd &> /dev/null; then
    echo "❌ azd not found. Install from https://aka.ms/azd"
    echo "   For local testing, you can still create a basic zip manually."
    exit 1
fi

# Check if azd x extensions are enabled
if ! azd extension list &> /dev/null; then
    echo "⚠️  azd extensions not enabled. Enabling..."
    azd config set alpha.extension.enabled on
    echo "📦 Installing microsoft.azd.extensions..."
    azd extension install microsoft.azd.extensions
fi

# Use azd x pack
echo "📦 Running azd x pack..."
azd x pack

echo ""
echo "✅ Package created in: ~/.azd/registry/jongio.azd.exec/${VERSION}/"
echo ""
echo "To test locally:"
echo "  azd extension source add -n local -t file -l \"\$HOME/.azd/registry/jongio.azd.exec/${VERSION}/registry.json\""
echo "  azd extension install jongio.azd.exec --version ${VERSION}"
