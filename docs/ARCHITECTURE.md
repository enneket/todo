# Architecture Documentation

## System Architecture

The application follows a **Decoupled Client-Server** architecture, packaged as a single Desktop application using Wails.

```mermaid
graph TD
    User[User] --> GUI[Wails Frontend (Vue 3)]
    GUI -- HTTP REST API (localhost:8081) --> Server[Go HTTP Server]
    Server --> Service[Business Logic Service]
    Service --> ProjectService[Project Logic]
    Service --> TodoService[Todo Logic]
    Service --> SubtaskService[Subtask Logic]
    Service --> DB[(SQLite Database)]
```

### Key Design Decisions

1. **HTTP over Wails Bindings**:
   - **Decision**: We use a standard Go HTTP server (`net/http`) instead of Wails' native JS bindings (`wails.runtime.*`).
   - **Reasoning**:
     - **Decoupling**: The frontend is a standard SPA that can be developed/tested in a browser without the Wails runtime.
     - **Standardization**: Uses standard REST patterns familiar to web developers.
     - **Flexibility**: The backend can easily serve other clients (mobile, web) in the future.

2. **SQLite Database**:
   - **Decision**: Use `modernc.org/sqlite` (pure Go implementation).
   - **Reasoning**: Removes the need for CGO, making cross-compilation easier and reducing runtime dependency issues (like `libc` versions).

3. **Data Relations**:
   - **Projects**: Todos can optionally belong to a Project.
   - **Subtasks**: Todos can contain multiple Subtasks (simple checklist items).
   - **Foreign Keys**: Enforced at DB level (`ON DELETE CASCADE` for Subtasks, `SET NULL` for Projects).

### Directory Structure

```
todo/
├── backend/            # Go Backend Code
│   ├── db/             # Database initialization and models
│   ├── server/         # HTTP Handlers and Routing
│   └── service/        # Business Logic
├── frontend/           # Vue 3 Frontend Code
│   ├── src/
│   │   ├── components/ # UI Components
│   │   ├── stores/     # Pinia State Stores
│   │   └── api/        # Axios API Client
│   └── cypress/        # E2E Tests
├── docs/               # Project Documentation
├── build/              # Build Artifacts
└── main.go             # Application Entry Point
```

### Components

1. **Wails Container**:
   - Manages Window lifecycle (Resize, Minimize, Close).
   - Hosts the Webview.
   - *Note*: Acts primarily as a "browser" shell.

2. **Frontend (Vue 3)**:
   - **UI**: Tailwind CSS + Phosphor Icons.
   - **State**: Pinia (Store).
   - **Network**: Axios (HTTP Client).
   - **I18n**: vue-i18n.

3. **Backend (Go)**:
   - **HTTP Server**: `net/http` standard library.
   - **Router**: Standard `http.ServeMux`.
   - **Database**: `database/sql` with `modernc.org/sqlite`.
