# Code Patterns

## Frontend (Vue/TS)

### Store Pattern (Pinia)
```typescript
export const useTodoStore = defineStore('todo', () => {
  const todos = ref<Todo[]>([])
  const API_URL = 'http://localhost:8081/api/todos'

  const fetchTodos = async () => { /* ... */ }
  
  return { todos, fetchTodos }
})
```

### Component Pattern
- Use `<script setup lang="ts">`.
- Import icons from `@phosphor-icons/vue` with `Ph` prefix.
- Use `useI18n()` for translations.

### E2E Testing Pattern (Cypress)
- **Select Elements**: Use `cy.contains('Text')` for user-visible text, or standard CSS selectors.
- **Interactions**: `.click()`, `.type('{enter}')`.
- **Assertions**: `.should('be.visible')`, `.should('not.exist')`, `.should('have.class', '...')`.

```typescript
it('should perform action', () => {
  cy.get('input').type('Value{enter}')
  cy.contains('Value').should('be.visible')
})
```

## Backend (Go)

### Handler Pattern
```go
func GetHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Call Service
    data, err := service.GetData()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // 2. Return JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
```

### Service Pattern
- **Separation of Concerns**: Logic split into `project.go`, `todo.go`, `subtask.go`.
- **Domain Models**: Returns structs defined in `db/models.go`.
- **Error Handling**: Returns standard Go errors, propagated to Handler.
- **Transactions**: (Future) Use `db.DB.Begin()` for complex multi-table updates.

### DB Pattern
- **Schema Management**: Tables created in `db.go` `InitDB`.
- **Migrations**: Simple `ALTER TABLE` statements in `InitDB` to support schema evolution.
- **Queries**: Parameterized queries `?` to prevent SQL injection.
