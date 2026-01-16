# Build Requirements & Setup Guide

This document outlines the requirements and steps to build and run the Todo App.

## System Requirements

### General
- **Go**: v1.24 or higher
- **Node.js**: v20 or higher
- **NPM**: v10 or higher
- **Wails CLI**: v2.11 or higher (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Linux Specifics (Debian/Ubuntu)
The application relies on WebKitGTK. You must install the development headers:

```bash
sudo apt update
sudo apt install -y pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev
```

> **Note**: If your system only has `libwebkit2gtk-4.0-dev`, you may run into build issues. We recommend `4.1`. If you must use `4.0`, you might need to adjust build tags, but our default Makefile setup assumes `4.1` or modern environments.

## Installation

1. **Clone the repository**:
   ```bash
   git clone <repo-url>
   cd todo
   ```

2. **Install dependencies**:
   ```bash
   make install-deps
   ```
   This command runs `go mod tidy` for the backend and `npm install` for the frontend.

## Running Development Server

To start the application in development mode (with hot-reload):

```bash
make dev
```

On Linux, this automatically applies the `-tags webkit2_41` flag to ensure compatibility.

## Building for Production

To create a production build:

```bash
make build
```

The executable will be placed in `build/bin/`.

## Troubleshooting

### "pkg-config: executable file not found" or "Package webkit2gtk-4.0 not found"
This usually happens on Linux when the required libraries are missing.
- **Fix**: Run the `sudo apt install ...` command listed in "Linux Specifics".
- **Fix**: Ensure you are using `make dev` which handles the build tags correctly.
