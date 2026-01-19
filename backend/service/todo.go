package service

import (
	"encoding/json"
	"time"
	"todo/backend/db"
)

func CreateTodo(title, description, priority string, dueDate *time.Time, tags []string) (int64, error) {
	if priority == "" {
		priority = "medium"
	}
	tagsJSON, _ := json.Marshal(tags)
	if tags == nil {
		tagsJSON = []byte("[]")
	}

	res, err := db.DB.Exec("INSERT INTO todos (title, description, priority, due_date, tags) VALUES (?, ?, ?, ?, ?)", title, description, priority, dueDate, string(tagsJSON))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTodos() ([]db.Todo, error) {
	rows, err := db.DB.Query("SELECT id, title, description, completed, priority, due_date, tags, created_at FROM todos ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []db.Todo
	for rows.Next() {
		var t db.Todo
		var tagsJSON string
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.Priority, &t.DueDate, &tagsJSON, &t.CreatedAt); err != nil {
			return nil, err
		}
		if tagsJSON != "" {
			json.Unmarshal([]byte(tagsJSON), &t.Tags)
		}
		if t.Tags == nil {
			t.Tags = []string{}
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func UpdateTodoStatus(id int, completed bool) error {
	_, err := db.DB.Exec("UPDATE todos SET completed = ? WHERE id = ?", completed, id)
	return err
}

func UpdateTodoDetails(id int, title, description, priority string, dueDate *time.Time, tags []string) error {
	tagsJSON, _ := json.Marshal(tags)
	if tags == nil {
		tagsJSON = []byte("[]")
	}
	_, err := db.DB.Exec("UPDATE todos SET title = ?, description = ?, priority = ?, due_date = ?, tags = ? WHERE id = ?", title, description, priority, dueDate, string(tagsJSON), id)
	return err
}

func DeleteTodo(id int) error {
	_, err := db.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}
