package service

import (
	"encoding/json"
	"time"
	"todo/backend/db"
)

func CreateTodo(title, description, priority string, dueDate, remindAt *time.Time, repeat string, tags []string, projectID *int) (int64, error) {
	if priority == "" {
		priority = "medium"
	}
	tagsJSON, _ := json.Marshal(tags)
	if tags == nil {
		tagsJSON = []byte("[]")
	}

	res, err := db.DB.Exec("INSERT INTO todos (title, description, priority, due_date, remind_at, repeat, tags, project_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", title, description, priority, dueDate, remindAt, repeat, string(tagsJSON), projectID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTodos() ([]db.Todo, error) {
	rows, err := db.DB.Query("SELECT id, title, description, completed, priority, due_date, remind_at, repeat, tags, project_id, created_at FROM todos ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []db.Todo
	for rows.Next() {
		var t db.Todo
		var tagsJSON string
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.Priority, &t.DueDate, &t.RemindAt, &t.Repeat, &tagsJSON, &t.ProjectID, &t.CreatedAt); err != nil {
			return nil, err
		}
		if tagsJSON != "" {
			json.Unmarshal([]byte(tagsJSON), &t.Tags)
		}
		if t.Tags == nil {
			t.Tags = []string{}
		}
		
		// Fetch subtasks for this todo
		subtasks, err := GetSubtasks(t.ID)
		if err == nil {
			t.Subtasks = subtasks
		} else {
			t.Subtasks = []db.Subtask{}
		}

		todos = append(todos, t)
	}
	return todos, nil
}

func UpdateTodoStatus(id int, completed bool) error {
	_, err := db.DB.Exec("UPDATE todos SET completed = ? WHERE id = ?", completed, id)
	if err != nil {
		return err
	}

	if completed {
		// Check for repeat
		var t db.Todo
		var tagsJSON string
		err := db.DB.QueryRow("SELECT title, description, priority, due_date, remind_at, repeat, tags, project_id FROM todos WHERE id = ?", id).Scan(
			&t.Title, &t.Description, &t.Priority, &t.DueDate, &t.RemindAt, &t.Repeat, &tagsJSON, &t.ProjectID,
		)
		if err == nil && t.Repeat != "" {
			// Calculate next dates
			nextDueDate := calculateNextDate(t.DueDate, t.Repeat)
			nextRemindAt := calculateNextDate(t.RemindAt, t.Repeat)
			
			if tagsJSON != "" {
				json.Unmarshal([]byte(tagsJSON), &t.Tags)
			}
			
			CreateTodo(t.Title, t.Description, t.Priority, nextDueDate, nextRemindAt, t.Repeat, t.Tags, t.ProjectID)
		}
	}

	return nil
}

func calculateNextDate(current *time.Time, repeat string) *time.Time {
	if current == nil {
		return nil
	}
	t := *current
	switch repeat {
	case "daily":
		t = t.AddDate(0, 0, 1)
	case "weekly":
		t = t.AddDate(0, 0, 7)
	case "monthly":
		t = t.AddDate(0, 1, 0)
	case "weekdays":
		// +1 day, then check if Sat/Sun
		t = t.AddDate(0, 0, 1)
		for t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
			t = t.AddDate(0, 0, 1)
		}
	}
	return &t
}

func UpdateTodoDetails(id int, title, description, priority string, dueDate, remindAt *time.Time, repeat string, tags []string, projectID *int) error {
	tagsJSON, _ := json.Marshal(tags)
	if tags == nil {
		tagsJSON = []byte("[]")
	}
	_, err := db.DB.Exec("UPDATE todos SET title = ?, description = ?, priority = ?, due_date = ?, remind_at = ?, repeat = ?, tags = ?, project_id = ? WHERE id = ?", title, description, priority, dueDate, remindAt, repeat, string(tagsJSON), projectID, id)
	return err
}

func DeleteTodo(id int) error {
	_, err := db.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}
