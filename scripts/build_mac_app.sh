#!/bin/bash

APP_NAME="SV Printer"
APP_DIR="${APP_NAME}.app"
CONTENTS_DIR="${APP_DIR}/Contents"
MACOS_DIR="${CONTENTS_DIR}/MacOS"
RESOURCES_DIR="${CONTENTS_DIR}/Resources"

# Version: explicit arg/env wins, otherwise nearest git tag (v-prefix stripped).
VERSION="${1:-${VERSION:-}}"
if [ -z "$VERSION" ]; then
    VERSION=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//')
fi
if [ -z "$VERSION" ]; then
    VERSION="0.0.0"
fi
echo "App version: ${VERSION}"

if [ -f "sv-printer" ]; then
    echo "Using existing Go binary..."
else
    echo "Building Go binary..."
    go build -o sv-printer ./cmd/sv-printer
fi

echo "Creating App Bundle structure..."
mkdir -p "$MACOS_DIR"
mkdir -p "$RESOURCES_DIR"

echo "Copying binary and icon..."
cp sv-printer "$MACOS_DIR/"
cp assets/logo.icns "${RESOURCES_DIR}/icon.icns"

echo "Creating Info.plist..."
cat << PLIST > "${CONTENTS_DIR}/Info.plist"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>sv-printer</string>
    <key>CFBundleIconFile</key>
    <string>icon.icns</string>
    <key>CFBundleIdentifier</key>
    <string>com.svtech.svprinter</string>
    <key>CFBundleName</key>
    <string>SV Printer</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>${VERSION}</string>
    <key>LSUIElement</key>
    <true/>
</dict>
</plist>
PLIST

echo "Done! ${APP_DIR} created successfully."
