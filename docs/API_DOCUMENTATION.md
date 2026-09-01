# API Documentation

The Todo App uses a RESTful HTTP API for communication between the Frontend (Vue) and Backend (Go).
The server runs on `http://localhost:8081`.

## Base URL
`http://localhost:8081/api`

## Endpoints

### Todos

#### `GET /api/todos`
- **Description**: Fetch all todos, optionally filtered by a search keyword.
- **Query Parameters**:
  - `q` (optional): Search keyword. Case-insensitive substring match against `title`, `description`, and `tags` (e.g. `GET /api/todos?q=milk`). The keyword is matched literally — `%`, `_`, and `\` are not treated as wildcards. Omit (or send empty) for the full list.
- **Response**: `200 OK`
  ```json
  [
    {
      "id": 1,
      "title": "Buy Milk",
      "description": "Go to store",
      "completed": false,
      "priority": "high",
      "due_date": "2023-10-01T10:00:00Z",
      "tags": ["personal"],
      "project_id": 1,
      "subtasks": [
        { "id": 1, "todo_id": 1, "title": "Get Wallet", "completed": true }
      ],
      "created_at": "..."
    }
  ]
  ```

#### `POST /api/todos`
- **Body**:
  ```json
  {
    "title": "Task Title",
    "description": "Optional Desc",
    "priority": "medium",
    "due_date": "2023-10-01T10:00:00Z",
    "tags": ["tag1"],
    "project_id": 1
  }
  ```
- **Response**: `200 OK` `{"id": 1}`

#### `PUT /api/todos/{id}`
- **Description**: Update todo details or status.
- **Body**: (Partial updates allowed)
  ```json
  {
    "title": "New Title",
    "completed": true,
    "project_id": 2
  }
  ```

#### `DELETE /api/todos/{id}`
- **Description**: Soft delete — moves the todo to the trash (see Trash section). The row stays in the database until purged.
- **Response**: `200 OK`

---

### Projects

#### `GET /api/projects`
- **Description**: Fetch all projects.
- **Response**: `200 OK`
  ```json
  [
    {
      "id": 1,
      "name": "Work",
      "description": "Office tasks",
      "color": "#3B82F6",
      "created_at": "..."
    }
  ]
  ```

#### `POST /api/projects`
- **Body**:
  ```json
  {
    "name": "Project Name",
    "description": "Optional Desc",
    "color": "#EF4444"
  }
  ```
- **Response**: `200 OK` `{"id": 1}`

#### `PUT /api/projects/{id}`
- **Body**:
  ```json
  {
    "name": "Updated Name",
    "description": "Updated Desc",
    "color": "#10B981"
  }
  ```
- **Response**: `200 OK`

#### `DELETE /api/projects/{id}`
- **Description**: Soft delete — moves the project and every todo under it to the trash (see Trash section). The rows stay in the database until purged.
- **Response**: `200 OK`

---

### Trash

#### `GET /api/trash/todos`
- **Description**: Fetch all soft-deleted todos, newest deletion first.
- **Response**: `200 OK` — array of todos (same shape as `GET /api/todos`, without subtask batching), each carrying `deleted_at`.

#### `GET /api/trash/projects`
- **Description**: Fetch all soft-deleted projects, newest deletion first.
- **Response**: `200 OK` — array of projects, each carrying `deleted_at`.

#### `POST /api/trash/todos/{id}/restore`
- **Description**: Restore a trashed todo to the normal list. If its project is also in the trash, the project is restored first so the todo never points at a deleted project; the project's other todos stay in the trash.
- **Response**: `200 OK`

#### `POST /api/trash/projects/{id}/restore`
- **Description**: Restore a trashed project. Its todos are left in the trash and can be restored individually via `POST /api/trash/todos/{id}/restore`.
- **Response**: `200 OK`

#### `DELETE /api/trash/todos/{id}`
- **Description**: Permanently delete a trashed todo (subtasks are cascade-deleted). Only works on rows already in the trash.
- **Response**: `200 OK`

#### `DELETE /api/trash/projects/{id}`
- **Description**: Permanently delete a trashed project and all its todos.
- **Response**: `200 OK`

---

### Subtasks

#### `POST /api/todos/{id}/subtasks`
- **Description**: Create a subtask for a specific todo.
- **Body**:
  ```json
  { "title": "Subtask Title" }
  ```
- **Response**: `200 OK` `{"id": 1}`

#### `PUT /api/subtasks/{id}`
- **Description**: Update subtask (e.g., toggle completion). Partial update — only the fields present in the body are changed; omitted fields keep their current values.
- **Body**:
  ```json
  { "completed": true }
  ```
  or `{ "title": "New Title" }`, or both.
- **Response**: `200 OK`

#### `DELETE /api/subtasks/{id}`
- **Response**: `200 OK`

## Data Model

### Todo
| Field | Type | Description |
|-------|------|-------------|
| `id` | `int` | Unique identifier (Auto-increment) |
| `title` | `string` | The task description |
| `completed` | `boolean` | Status of the task |
| `deleted_at` | `datetime` or `null` | Soft-delete timestamp; `null` for active todos, set while the todo is in the trash |
