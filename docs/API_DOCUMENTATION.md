# API Documentation

The Todo App uses a RESTful HTTP API for communication between the Frontend (Vue) and Backend (Go).
The server runs on `http://localhost:8081`.

## Base URL
`http://localhost:8081/api`

## Endpoints

### 1. Get All Todos
Fetch the list of all todo items.

- **URL**: `/todos`
- **Method**: `GET`
- **Response**: `200 OK`
- **Body**: JSON Array of Todo objects

```json
[
  {
    "id": 1,
    "title": "Buy milk",
    "completed": false
  },
  {
    "id": 2,
    "title": "Walk the dog",
    "completed": true
  }
]
```

### 2. Create Todo
Create a new todo item.

- **URL**: `/todos`
- **Method**: `POST`
- **Content-Type**: `application/json`
- **Request Body**:

```json
{
  "title": "New Task"
}
```

- **Response**: `201 Created`
- **Body**: JSON Object of the created Todo (with ID)

### 3. Update Todo
Update an existing todo item (e.g., mark as completed).

- **URL**: `/todos/{id}`
- **Method**: `PUT`
- **Content-Type**: `application/json`
- **Request Body**:

```json
{
  "title": "Updated Task Title",
  "completed": true
}
```

- **Response**: `200 OK`
- **Body**: JSON Object of the updated Todo

### 4. Delete Todo
Delete a todo item.

- **URL**: `/todos/{id}`
- **Method**: `DELETE`
- **Response**: `200 OK`

## Data Model

### Todo
| Field | Type | Description |
|-------|------|-------------|
| `id` | `int` | Unique identifier (Auto-increment) |
| `title` | `string` | The task description |
| `completed` | `boolean` | Status of the task |
