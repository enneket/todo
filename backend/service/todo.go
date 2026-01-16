package service

import (
	"todo/backend/db"
)

func CreateTodo(title string) (int64, error) {
	res, err := db.DB.Exec("INSERT INTO todos (title) VALUES (?)", title)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTodos() ([]db.Todo, error) {
	rows, err := db.DB.Query("SELECT id, title, completed, created_at FROM todos ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []db.Todo
	for rows.Next() {
		var t db.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func UpdateTodo(id int, completed bool) error {
	_, err := db.DB.Exec("UPDATE todos SET completed = ? WHERE id = ?", completed, id)
	return err
}

func DeleteTodo(id int) error {
	_, err := db.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}
