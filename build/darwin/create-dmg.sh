#!/bin/bash
set -e
APP_NAME="TodoApp"
VERSION=$(grep -o '"version"[[:space:]]*:[[:space:]]*"[^"]*"' frontend/package.json 2>/dev/null | head -1 | sed 's/.*"\([^"]*\)".*/\1/' || echo "1.0.0")
BUILD_DIR="build/bin"
DMG_DIR="build/dmg"
APP_PATH="${BUILD_DIR}/${APP_NAME}.app"
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
