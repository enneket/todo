package service

import (
	"strings"

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

func GetSubtask(id int) (db.Subtask, error) {
	row := db.DB.QueryRow("SELECT id, todo_id, title, completed, created_at FROM subtasks WHERE id = ?", id)
	var s db.Subtask
	if err := row.Scan(&s.ID, &s.TodoID, &s.Title, &s.Completed, &s.CreatedAt); err != nil {
		return db.Subtask{}, err
	}
	return s, nil
}

func UpdateSubtask(id int, title *string, completed *bool) error {
	// Build a dynamic SET clause so partial updates only touch the fields the
	// client actually sent — toggling completion must not blank the title.
	var sets []string
	var args []interface{}
	if title != nil {
		sets = append(sets, "title = ?")
		args = append(args, *title)
	}
	if completed != nil {
		sets = append(sets, "completed = ?")
		args = append(args, *completed)
	}
	if len(sets) == 0 {
		return nil // no-op update
	}
	args = append(args, id)
	_, err := db.DB.Exec("UPDATE subtasks SET "+strings.Join(sets, ", ")+" WHERE id = ?", args...)
	return err
}

func DeleteSubtask(id int) error {
	_, err := db.DB.Exec("DELETE FROM subtasks WHERE id = ?", id)
	return err
}
