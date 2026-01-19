package service

import (
	"todo/backend/db"
)

func CreateSubtask(todoID int, title string) (int64, error) {
	res, err := db.DB.Exec("INSERT INTO subtasks (todo_id, title) VALUES (?, ?)", todoID, title)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetSubtasks(todoID int) ([]db.Subtask, error) {
	rows, err := db.DB.Query("SELECT id, todo_id, title, completed, created_at FROM subtasks WHERE todo_id = ? ORDER BY created_at ASC", todoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subtasks []db.Subtask
	for rows.Next() {
		var s db.Subtask
		if err := rows.Scan(&s.ID, &s.TodoID, &s.Title, &s.Completed, &s.CreatedAt); err != nil {
			return nil, err
		}
		subtasks = append(subtasks, s)
	}
	return subtasks, nil
}

func UpdateSubtask(id int, title string, completed bool) error {
	_, err := db.DB.Exec("UPDATE subtasks SET title = ?, completed = ? WHERE id = ?", title, completed, id)
	return err
}

func DeleteSubtask(id int) error {
	_, err := db.DB.Exec("DELETE FROM subtasks WHERE id = ?", id)
	return err
}
