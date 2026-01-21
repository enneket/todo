# Testing Strategy

## Backend Testing

- **Unit Tests**: Located in `backend/service/service_test.go` and `backend/server/server_test.go`.
- **Framework**: Standard Go `testing` package with **In-Memory SQLite** (`:memory:`) and `httptest`.
- **Run Tests**: `go test -v ./backend/...` (or `make test`)
- **Run with Coverage**: `go test -coverprofile=coverage.out ./backend/... && go tool cover -func=coverage.out`

## Frontend Testing

- **Unit Tests**: Located in `frontend/src/stores/*.test.ts` and `frontend/src/components/*.test.ts`.
- **Framework**: **Vitest** + **HappyDOM** + **Axios Mocking** + **Vue Test Utils**.
- **Run Tests**: `npm run test` (in frontend dir) or `make test` (from root).
- **Run with Coverage**: `npm run test -- --coverage`

## Integration Testing

- **Manual**: Run `make dev` and interact with the UI.
- **API**: Use Postman/Curl to hit `localhost:8081/api/*`.

## E2E Testing (Cypress)

We use Cypress for End-to-End testing of the full application flow (Frontend + Backend).

### Prerequisites
- The application must be running (`make dev`).
- Frontend accessible at `http://localhost:5173` (Vite default).

### Running Tests
- **Headless Mode**: `make test-e2e`
- **Interactive Mode**: `cd frontend && npx cypress open`

### Test Files
- `frontend/cypress/e2e/todo.cy.ts`: Covers Todos (Add/Complete), Projects (Create/Filter), and basic navigation.
- `frontend/cypress/e2e/tags.cy.ts`: Covers Tag system (Add tags, Filter by Sidebar/Card).

## Test Coverage Reports (As of Latest Run)

### Backend Coverage (Unit Tests)
| Package | Coverage |
|---------|----------|
| `backend/service` | ~85% |
| `backend/server` | ~38% |
| **Total** | **~48%** |

### Frontend Coverage (Unit Tests)
| Category | Coverage |
|----------|----------|
| Components | ~87% |
| Stores | ~69% |
| **Total** | **~73%** |

### E2E Coverage (Frontend)
| Type | Coverage |
|------|----------|
| Lines | ~68% |
| Statements | ~66% |
| **Key Flows** | **100% Pass** |
