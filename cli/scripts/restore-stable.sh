#!/bin/bash
# Restore the stable version of azd exec extension

set -e

# Colors
CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
GRAY='\033[0;90m'
WHITE='\033[0;37m'
NC='\033[0m' # No Color

REPO="jongio/azd-exec"
EXTENSION_ID="jongio.azd.exec"
STABLE_REGISTRY_URL="https://jongio.github.io/azd-extensions/registry.json"

echo -e "${CYAN}🔄 Restoring stable azd exec extension${NC}"
echo ""

# Step 1: Uninstall current extension
echo -e "${GRAY}🗑️  Uninstalling current extension...${NC}"
azd extension uninstall $EXTENSION_ID 2>/dev/null || true
# Ignore errors - extension might not be installed

# Step 2: Remove all PR registry sources
echo -e "${GRAY}🧹 Removing PR registry sources...${NC}"
SOURCES=$(azd extension source list --output json 2>/dev/null || echo "[]")
if [ "$SOURCES" != "[]" ]; then
    echo "$SOURCES" | grep -o '"name":"pr-[0-9]*"' | sed 's/"name":"\(.*\)"/\1/' | while read -r source; do
        if [ ! -z "$source" ]; then
            echo -e "${GRAY}   Removing: $source${NC}"
            azd extension source remove "$source" 2>/dev/null || true
        fi
    done
fi

# Step 3: Clean up pr-registry.json files
echo -e "${GRAY}🧹 Cleaning up pr-registry.json files...${NC}"
rm -f ./pr-registry.json
rm -f ~/pr-registry.json
rm -f $HOME/pr-registry.json

# Step 4: Add stable registry source
echo -e "${GRAY}🔗 Adding stable registry source...${NC}"
# Remove if exists
azd extension source remove "jongio" 2>/dev/null || true
azd extension source add -n "jongio" -t url -l "$STABLE_REGISTRY_URL"
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Failed to add stable registry source${NC}"
    exit 1
fi

# Step 5: Install stable version
echo -e "${GRAY}📦 Installing latest stable version...${NC}"
azd extension install $EXTENSION_ID
if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Failed to install stable version${NC}"
    exit 1
fi

# Step 6: Verify installation
echo ""
echo -e "${GREEN}✅ Restoration complete!${NC}"
echo ""
echo -e "${GRAY}🔍 Verifying installation...${NC}"
INSTALLED_VERSION=$(azd exec version 2>&1)
if [ $? -eq 0 ]; then
    echo -e "${GRAY}   $INSTALLED_VERSION${NC}"
    echo ""
    echo -e "${GREEN}✨ Success! Stable version restored.${NC}"
else
    echo -e "${YELLOW}⚠️  Could not verify version${NC}"
fi

echo ""
echo -e "${CYAN}Try these commands:${NC}"
echo -e "${WHITE}  azd exec version${NC}"
echo -e "${WHITE}  azd exec ./my-script.sh${NC}"
