# Testing Strategy

## Backend Testing

- **Unit Tests**: Located in `backend/service/service_test.go`.
- **Framework**: Standard Go `testing` package with **In-Memory SQLite** (`:memory:`).
- **Run Tests**: `go test -v ./backend/service` (or `make test`)

## Frontend Testing

- **Unit Tests**: Located in `frontend/src/stores/*.test.ts`.
- **Framework**: **Vitest** + **HappyDOM** + **Axios Mocking**.
- **Run Tests**: `npm run test` (in frontend dir) or `make test` (from root).

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
