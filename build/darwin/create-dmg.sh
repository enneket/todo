#!/bin/bash
set -e
APP_NAME="TodoApp"

# Allow version to be passed as first argument
if [ -n "$1" ]; then
    VERSION="$1"
else
    VERSION=$(grep -o '"version"[[:space:]]*:[[:space:]]*"[^"]*"' frontend/package.json 2>/dev/null | head -1 | sed 's/.*"\([^"]*\)".*/\1/' || echo "1.0.0")
fi

BUILD_DIR="build/bin"

# Dynamically find the .app bundle to handle potential naming variations
APP_PATH=$(find "${BUILD_DIR}" -maxdepth 1 -name "*.app" | head -n 1)

if [ -z "${APP_PATH}" ]; then
    echo "Error: App bundle not found in ${BUILD_DIR}"
    echo "Contents of ${BUILD_DIR}:"
    ls -la "${BUILD_DIR}"
    exit 1
fi

# Use the found app name for DMG creation
APP_FILENAME=$(basename "${APP_PATH}")

DMG_DIR="build/dmg"
DMG_NAME="${APP_NAME}-${VERSION}-darwin-universal.dmg"

if [ ! -d "${APP_PATH}" ]; then
    echo "Error: App bundle not found at ${APP_PATH}"
    exit 1
fi

rm -rf "${DMG_DIR}"
mkdir -p "${DMG_DIR}"
cp -R "${APP_PATH}" "${DMG_DIR}/"
ln -s /Applications "${DMG_DIR}/Applications"
rm -f "${BUILD_DIR}/${DMG_NAME}"
hdiutil create -volname "${APP_NAME}" -srcfolder "${DMG_DIR}" -ov -format UDZO "${BUILD_DIR}/${DMG_NAME}"
rm -rf "${DMG_DIR}"
echo "DMG created: ${BUILD_DIR}/${DMG_NAME}"
