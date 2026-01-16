# Todo App

A modern, privacy-focused, cross-platform desktop Todo application built with Go (Wails) and Vue 3.

[中文文档](README_zh-CN.md)

## ✨ Features

- **Cross-Platform**: Runs on Windows, Linux, and macOS.
- **Modern UI/UX**: Built with Vue 3, Tailwind CSS, and Phosphor Icons.
- **Privacy-Focused**: All data is stored locally using SQLite.
- **Robust Backend**: Powered by Go 1.24+ and Wails v2.
- **Developer Friendly**: Unified workflow via Makefile.

## 🛠 Tech Stack

- **Backend**: Go 1.24+, Wails v2.11+, SQLite (`modernc.org/sqlite`)
- **Frontend**: Vue 3.5+, Pinia, Tailwind CSS 3.3+, Vite 5+, TypeScript
- **Communication**: HTTP REST API (localhost:8081)

## 🚀 Getting Started

### Prerequisites

- **Go**: v1.24 or higher
- **Node.js**: v18 or higher
- **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **Make**: For running unified commands (Windows users can use WSL or ensure Make is installed)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/todo-app.git
   cd todo-app
   ```

2. Install dependencies:
   ```bash
   make setup
   ```

### Development

Start the application in development mode with hot-reload:

```bash
make dev
```

The frontend will run on a dedicated Vite server, and the backend will compile and run in desktop mode.

## 📦 Building

To build the application for production:

```bash
make build
```

The binary will be located in `build/bin/`.

### Cross-Platform Build

- **Windows**: `make build-windows`
- **Linux**: `make build-linux`
- **macOS**: `make build-darwin`

## 🧪 Testing & Quality

- **Run all checks**: `make check` (Lint + Test + Build)
- **Unit Tests**: `make test`
- **E2E Tests**: `make test-e2e` (Requires running app)
- **Linting**: `make lint`
- **Formatting**: `make format`

## 📂 Project Structure

```
├── backend/      # Go backend code (Server, DB, Service)
├── frontend/     # Vue frontend code (Stores, Components)
├── docs/         # Documentation
├── build/        # Build artifacts
├── main.go       # Application entry point
└── Makefile      # Unified development commands
```

## 📚 Documentation

For more detailed documentation, please refer to the `docs/` directory:

- [Architecture](docs/ARCHITECTURE.md)
- [Code Patterns](docs/CODE_PATTERNS.md)
- [API Documentation](docs/API_DOCUMENTATION.md)

## 📄 License

MIT License
