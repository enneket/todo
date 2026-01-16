# Testing Strategy

## Backend Testing

- **Unit Tests**: Place `_test.go` files next to source files.
- **Run Tests**: `go test ./...` (or `make test-backend`)

## Frontend Testing

- **Unit Tests**: Use Vitest (if configured).
- **Run Tests**: `npm run test` (or `make test-frontend`)

## Integration Testing

- Currently manual verification via `wails dev`.
- Ensure backend API is responding correctly via `curl` or Postman when running.

## E2E Testing (Cypress)

We use Cypress for End-to-End testing of the frontend application.

### Prerequisites
- The application must be running (e.g., via `make dev` or `wails dev`).
- The frontend dev server should be accessible at `http://localhost:5175`.

### Running Tests
- **Headless Mode**: `make test-e2e` (or `cd frontend && npx cypress run`)
- **Interactive Mode**: `cd frontend && npx cypress open`

### Test Files
- `frontend/cypress/e2e/todo.cy.ts`: Covers CRUD operations, filtering, and interactions.
