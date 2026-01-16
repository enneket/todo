# GitHub Copilot Instructions

> **Quick Links**: [Architecture](../docs/ARCHITECTURE.md) | [Code Patterns](../docs/CODE_PATTERNS.md) | [Testing](../docs/TESTING.md) | [Version Management](../docs/VERSION_MANAGEMENT.md)

This project uses **Wails + Vue + Go**. Please follow these guidelines when generating code:

1. **Communication Protocol**:
   - ALWAYS use `axios` to call `http://localhost:8081/api/...` for backend interaction.
   - DO NOT suggest using Wails bindings (`@/wailsjs/...`) for business logic.

2. **Frontend Style**:
   - Use **Vue 3 Composition API** (`<script setup lang="ts">`).
   - Use **Tailwind CSS** for styling.
   - Use **Phosphor Icons** with `Ph` prefix (e.g., `PhCheckCircle`).

3. **Backend Style**:
   - Follow standard Go project layout.
   - Use `modernc.org/sqlite` for database driver.
   - Handlers should be in `backend/server`, Logic in `backend/service`.

4. **Development**:
   - Use `make dev` to start the server.
   - If `pkg-config` error occurs regarding `webkit2gtk-4.0`, ensure `-tags webkit2_41` is used (handled by `make dev`).
