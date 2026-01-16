#!/bin/bash
set -e
APP_NAME="TodoApp"
VERSION=$(grep -o '"version"[[:space:]]*:[[:space:]]*"[^"]*"' frontend/package.json 2>/dev/null | head -1 | sed 's/.*"\([^"]*\)".*/\1/' || echo "1.0.0")
ARCH=${ARCH:-amd64}
APP_PUBLISHER="TodoApp Developer"
APP_DESCRIPTION="A Modern Todo Application"
BUILD_DIR="build/bin"
APPDIR="build/appimage/${APP_NAME}.AppDir"
APPIMAGE_NAME="${APP_NAME}-${VERSION}-linux-${ARCH}.AppImage"

if [ ! -f "${BUILD_DIR}/${APP_NAME}" ]; then
    echo "Error: Binary not found at ${BUILD_DIR}/${APP_NAME}"
    exit 1
fi

rm -rf "build/appimage"
mkdir -p "${APPDIR}/usr/bin"
mkdir -p "${APPDIR}/usr/share/applications"
mkdir -p "${APPDIR}/usr/share/icons/hicolor/256x256/apps"

cp "${BUILD_DIR}/${APP_NAME}" "${APPDIR}/usr/bin/"
chmod +x "${APPDIR}/usr/bin/${APP_NAME}"

cat > "${APPDIR}/usr/share/applications/${APP_NAME}.desktop" << EOF
[Desktop Entry]
Type=Application
Name=${APP_NAME}
Comment=${APP_DESCRIPTION}
Exec=${APP_NAME}
Icon=${APP_NAME}
Categories=Utility;
Terminal=false
StartupWMClass=${APP_NAME}
EOF

cat > "${APPDIR}/AppRun" << 'EOF'
#!/bin/bash
SELF=$(readlink -f "$0")
HERE=${SELF%/*}
export PATH="${HERE}/usr/bin:${PATH}"
export LD_LIBRARY_PATH="${HERE}/usr/lib:${LD_LIBRARY_PATH}"
exec "${HERE}/usr/bin/TodoApp" "$@"
EOF
chmod +x "${APPDIR}/AppRun"

if [ -f "build/appicon.png" ]; then
    cp "build/appicon.png" "${APPDIR}/usr/share/icons/hicolor/256x256/apps/${APP_NAME}.png"
    cp "build/appicon.png" "${APPDIR}/${APP_NAME}.png"
else
    echo "Warning: No icon found at build/appicon.png"
fi

cp "${APPDIR}/usr/share/applications/${APP_NAME}.desktop" "${APPDIR}/"

APPIMAGE_ARCH="${ARCH}"
if [ "${ARCH}" = "amd64" ]; then
    APPIMAGE_ARCH="x86_64"
elif [ "${ARCH}" = "arm64" ]; then
    APPIMAGE_ARCH="aarch64"
fi

APPIMAGETOOL="build/appimagetool-${APPIMAGE_ARCH}.AppImage"
if [ ! -f "${APPIMAGETOOL}" ]; then
    wget -q "https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-${APPIMAGE_ARCH}.AppImage" -O "${APPIMAGETOOL}"
    chmod +x "${APPIMAGETOOL}"
fi

ARCH="${APPIMAGE_ARCH}" "${APPIMAGETOOL}" --no-appstream --verbose "${APPDIR}" "${BUILD_DIR}/${APPIMAGE_NAME}"
rm -rf "build/appimage"
echo "AppImage created: ${BUILD_DIR}/${APPIMAGE_NAME}"
