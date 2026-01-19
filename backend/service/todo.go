package service

import (
	"time"
	"todo/backend/db"
)

func CreateTodo(title, description, priority string, dueDate *time.Time) (int64, error) {
	if priority == "" {
		priority = "medium"
	}
	res, err := db.DB.Exec("INSERT INTO todos (title, description, priority, due_date) VALUES (?, ?, ?, ?)", title, description, priority, dueDate)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTodos() ([]db.Todo, error) {
	rows, err := db.DB.Query("SELECT id, title, description, completed, priority, due_date, created_at FROM todos ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []db.Todo
	for rows.Next() {
		var t db.Todo
		// Handle description being potentially NULL if DB was manually messed with, but we set DEFAULT ''
		// Scan automatically handles empty strings for TEXT
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.Priority, &t.DueDate, &t.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func UpdateTodoStatus(id int, completed bool) error {
	_, err := db.DB.Exec("UPDATE todos SET completed = ? WHERE id = ?", completed, id)
	return err
}

func UpdateTodoDetails(id int, title, description, priority string, dueDate *time.Time) error {
	_, err := db.DB.Exec("UPDATE todos SET title = ?, description = ?, priority = ?, due_date = ? WHERE id = ?", title, description, priority, dueDate, id)
	return err
}

func DeleteTodo(id int) error {
	_, err := db.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}
