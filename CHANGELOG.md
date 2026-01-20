# Changelog

All notable changes to this project will be documented in this file.

## [0.1.3] - 2026-01-20

### Added

#### ✨ Features
- **Task Notifications**: Added native desktop notifications for task reminders.
- **Repeating Tasks**: Added support for recurring tasks (Daily, Weekly, Monthly, Weekdays).
- **Reminders**: Users can now set specific reminder times for tasks, separate from due dates.

#### 🧪 Testing
- **E2E Testing**: Added comprehensive E2E tests for new features (reminders, repeat rules) using Cypress.
- **Accessibility**: Added basic a11y checks in E2E suite.

### Changed

#### 🎨 UI/UX
- **Unified Input Height**: Standardized the height of all input fields in the task creation modal for better visual consistency.
- **Visual Indicators**: Added bell and repeat icons to task cards to indicate active reminders and recurrence rules.

#### 📚 Documentation
- **Updated Docs**: Updated README and Chinese documentation with new features and Linux dependencies (`libnotify-dev`).

## [0.1.2] - 2026-01-19

### Added

#### ✨ Features
- **Tags System**: Added full support for task tags (add, edit, display, and search).
- **Task Description**: Added a detailed description field for tasks.
- **Enhanced Search**: Search now filters by title, description, and tags.
- **Search Bar**: Added a dedicated search bar in the application header.

#### 🎨 UI/UX
- **Modal Task Creation**: Replaced inline input with a "Add" button and a comprehensive modal dialog for creating tasks.
- **Tag Visualization**: Tags are now displayed as pills in the task list.
- **Improved List View**: Descriptions are now shown in the task list (truncated).

### Changed
- **Database**: Added `description` and `tags` columns to the schema.
- **API**: Updated Create/Update endpoints to handle tags and descriptions.

### Fixed
- **Testing**: Added headless server support (`backend/cmd/server`) for reliable CI/CD testing.
- **Code Quality**: Fixed various frontend linting issues and type definitions.

## [0.1.1] - 2026-01-16

### Added

#### ✨ Features
- **Task Priority**: Added support for High, Medium, and Low priorities with visual indicators.
- **Due Dates**: Added ability to set due dates for tasks.
- **Smart Defaults**: New tasks automatically default to a due date of tomorrow.
- **Task Editing**: Users can now edit task details (title, priority, due date) after creation.

#### 🎨 UI/UX
- **Custom Select Component**: Introduced `BaseSelect` for consistent rendering across platforms (fixes Linux WebKit font issues).
- **Better Typography**: Added `Noto Sans SC` webfont to ensure correct Chinese character rendering.

### Changed
- **UI Alignment**: Unified height and alignment for all input controls in the add task bar.
- **Database**: Updated schema to include `priority` and `due_date` columns (auto-migration included).
- **API**: Enhanced Create/Update endpoints to support new fields.

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
