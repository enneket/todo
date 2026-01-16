# AI Agent Guidelines for Todo App

> **Quick Links**: [Architecture](docs/ARCHITECTURE.md) | [Code Patterns](docs/CODE_PATTERNS.md) | [Testing](docs/TESTING.md) | [Build Requirements](docs/BUILD_REQUIREMENTS.md) | [API Docs](docs/API_DOCUMENTATION.md)

## Project Overview

**Todo App** is a modern, privacy-focused, cross-platform desktop Todo application.

### Tech Stack

- **Backend**: Go 1.24+ with Wails v2.11+ framework, SQLite with `modernc.org/sqlite`
- **Frontend**: Vue 3.5+ Composition API, Pinia, Tailwind CSS 3.3+, Vite 5+, TypeScript
- **Communication**: HTTP REST API (running on localhost:8081), NOT Wails bindings
- **Icons**: Phosphor Icons | **I18n**: vue-i18n (English/Chinese)
- **Testing**: Cypress (E2E), Go Test (Backend)

## Development Workflow

### Unified Commands (Makefile)

Use `make` to run common tasks. Avoid running raw commands manually.

- `make dev`: Start the application in development mode (with webkit2_41 support)
- `make build`: Build the application for production
- `make lint`: Run linters for both backend and frontend
- `make format`: Format code for both backend and frontend
- `make test`: Run unit tests
- `make test-e2e`: Run Cypress E2E tests (requires app running)
- `make clean`: Clean build artifacts

### Directory Structure

- `backend/`: Go backend code
  - `server/`: HTTP server implementation
  - `db/`: Database models and initialization
  - `service/`: Business logic
- `frontend/`: Vue frontend code
  - `src/stores/`: Pinia stores
  - `src/locales/`: i18n JSON files
  - `cypress/`: E2E tests
- `docs/`: Project documentation for AI and Developers

## Critical Constraints

1. **HTTP Communication**: The Frontend MUST communicate with the Backend via HTTP REST API (localhost:8081). Do NOT use Wails runtime bindings (`wails.runtime.*`) for data transmission.
2. **WebKit Version**: On Linux, always use `-tags webkit2_41` when building or running dev to ensure compatibility with modern WebKitGTK.
