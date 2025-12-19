#!/bin/bash
set -euo pipefail

# pack.sh - Create azd extension package for release
# Usage: ./pack.sh [version]

VERSION="${1:-0.1.0}"
OUT_DIR="out"
PKG_NAME="azd-exec-${VERSION}.zip"

echo "Packaging azd exec extension version ${VERSION}"

rm -rf "${OUT_DIR}"
mkdir -p "${OUT_DIR}"/bin

# Copy binary and metadata
if [ -f bin/exec ]; then
  cp bin/exec "${OUT_DIR}/bin/"
fi
if [ -f bin/exec.exe ]; then
  cp bin/exec.exe "${OUT_DIR}/bin/"
fi

cp extension.yaml "${OUT_DIR}/"
cp README.md "${OUT_DIR}/"

pushd "${OUT_DIR}" >/dev/null
zip -r "../${PKG_NAME}" .
popd >/dev/null

echo "Created package: ${PKG_NAME}"
