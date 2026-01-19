# API Documentation

The Todo App uses a RESTful HTTP API for communication between the Frontend (Vue) and Backend (Go).
The server runs on `http://localhost:8081`.

## Base URL
`http://localhost:8081/api`

## Endpoints

### Todos

#### `GET /api/todos`
- **Description**: Fetch all todos.
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
- **Description**: Update subtask (e.g., toggle completion).
- **Body**:
  ```json
  { "completed": true, "title": "New Title" }
  ```
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
