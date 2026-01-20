package db

import "time"

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
}

type Subtask struct {
	ID        int       `json:"id"`
	TodoID    int       `json:"todo_id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type Todo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Priority    string    `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	RemindAt    *time.Time `json:"remind_at"`
	Repeat      string    `json:"repeat"`
	Tags        []string  `json:"tags"`
	ProjectID   *int      `json:"project_id"` // Nullable
	Subtasks    []Subtask `json:"subtasks,omitempty"` // For API response
	CreatedAt   time.Time `json:"created_at"`
}
