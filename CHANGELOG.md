# Changelog

All notable changes to this project will be documented in this file.

## [0.1.0] - 2026-01-16

### Added

#### ✨ Core Features
- Initial release of the Todo App.
- Cross-platform desktop application support (Windows, Linux, macOS) powered by Wails v2.
- Local data privacy with SQLite storage (`modernc.org/sqlite`).
- Dedicated HTTP REST API server running on `localhost:8081` for frontend-backend communication.
- Internationalization (i18n) support with English (default) and Chinese locales.

#### 🎨 Frontend
- Modern UI built with Vue 3.5+, Tailwind CSS 3.3+, and Phosphor Icons.
- State management using Pinia.
- TypeScript support for type safety.

#### 🛠 Infrastructure & Tooling
- Unified development workflow via `Makefile` (`make dev`, `make build`, `make test`, etc.).
- Automated testing setup:
  - Backend: Go standard library testing.
  - Frontend: Cypress for E2E testing.
- Comprehensive documentation in `docs/`:
  - Architecture overview (`ARCHITECTURE.md`)
  - API documentation (`API_DOCUMENTATION.md`)
  - Code patterns and guidelines (`CODE_PATTERNS.md`)
  - Build and testing guides.

#### 🏗 Tech Stack
- **Backend**: Go 1.24+, Wails v2.11+, SQLite.
- **Frontend**: Vue 3, Vite 5, Pinia, TypeScript.
