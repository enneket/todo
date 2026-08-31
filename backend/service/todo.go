package service

import (
	"encoding/json"
	"strings"
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
	rows, err := db.DB.Query("SELECT id, title, description, completed, priority, due_date, remind_at, notified_at, repeat, tags, project_id, created_at FROM todos ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}

	// Drain rows into memory and close before opening a second query. Holding
	// rows open while we run another Query / Exec would deadlock the connection
	// pool — see checkReminders for the same pattern.
	var todos []db.Todo
	for rows.Next() {
		var t db.Todo
		var tagsJSON string
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.Priority, &t.DueDate, &t.RemindAt, &t.NotifiedAt, &t.Repeat, &tagsJSON, &t.ProjectID, &t.CreatedAt); err != nil {
			rows.Close()
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
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(todos) == 0 {
		return todos, nil
	}

	// Batch-fetch every subtask for every todo in a single query, then group
	// by todo_id. Replaces the per-todo GetSubtasks call (N+1 → 2 queries).
	placeholders := make([]string, len(todos))
	ids := make([]interface{}, len(todos))
	for i, t := range todos {
		placeholders[i] = "?"
		ids[i] = t.ID
	}
	subQuery := "SELECT id, todo_id, title, completed, created_at FROM subtasks WHERE todo_id IN (" + strings.Join(placeholders, ",") + ") ORDER BY created_at ASC"

	subRows, err := db.DB.Query(subQuery, ids...)
	if err != nil {
		return nil, err
	}
	defer subRows.Close()

	byTodo := make(map[int][]db.Subtask, len(todos))
	for subRows.Next() {
		var s db.Subtask
		if err := subRows.Scan(&s.ID, &s.TodoID, &s.Title, &s.Completed, &s.CreatedAt); err != nil {
			return nil, err
		}
		byTodo[s.TodoID] = append(byTodo[s.TodoID], s)
	}
	if err := subRows.Err(); err != nil {
		return nil, err
	}

	for i := range todos {
		if subs, ok := byTodo[todos[i].ID]; ok {
			todos[i].Subtasks = subs
		} else {
			todos[i].Subtasks = []db.Subtask{}
		}
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

// UpdateTodoParams is a partial-update payload: any nil field is left alone
// in the database. Pass a non-nil pointer to set that field. Passing an empty
// slice for Tags (vs. nil) is meaningful — it clears the tag list.
type UpdateTodoParams struct {
	Title       *string
	Description *string
	Priority    *string
	DueDate     *time.Time
	RemindAt    *time.Time
	Repeat      *string
	Tags        *[]string
	ProjectID   *int
}

func UpdateTodoDetails(id int, p UpdateTodoParams) error {
	sets := []string{}
	args := []interface{}{}

	if p.Title != nil {
		sets = append(sets, "title = ?")
		args = append(args, *p.Title)
	}
	if p.Description != nil {
		sets = append(sets, "description = ?")
		args = append(args, *p.Description)
	}
	if p.Priority != nil {
		sets = append(sets, "priority = ?")
		args = append(args, *p.Priority)
	}
	if p.DueDate != nil {
		sets = append(sets, "due_date = ?")
		args = append(args, *p.DueDate)
	}
	if p.RemindAt != nil {
		sets = append(sets, "remind_at = ?")
		args = append(args, *p.RemindAt)
	}
	if p.Repeat != nil {
		sets = append(sets, "repeat = ?")
		args = append(args, *p.Repeat)
	}
	if p.Tags != nil {
		tagsJSON, err := json.Marshal(*p.Tags)
		if err != nil {
			return err
		}
		sets = append(sets, "tags = ?")
		args = append(args, string(tagsJSON))
	}
	if p.ProjectID != nil {
		sets = append(sets, "project_id = ?")
		args = append(args, *p.ProjectID)
	}

	if len(sets) == 0 {
		return nil
	}

	// Reset notified_at so any reminder reschedules correctly after edits.
	sets = append(sets, "notified_at = NULL")

	query := "UPDATE todos SET " + strings.Join(sets, ", ") + " WHERE id = ?"
	args = append(args, id)
	_, err := db.DB.Exec(query, args...)
	return err
}

func DeleteTodo(id int) error {
	_, err := db.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}

// GetTodo fetches a single todo (with its subtasks). Used by handlers that
// just created/updated a row so they can return the full object to clients
// without the client having to re-GET the whole list.
func GetTodo(id int) (db.Todo, error) {
	row := db.DB.QueryRow("SELECT id, title, description, completed, priority, due_date, remind_at, notified_at, repeat, tags, project_id, created_at FROM todos WHERE id = ?", id)
	var t db.Todo
	var tagsJSON string
	if err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Completed, &t.Priority, &t.DueDate, &t.RemindAt, &t.NotifiedAt, &t.Repeat, &tagsJSON, &t.ProjectID, &t.CreatedAt); err != nil {
		return db.Todo{}, err
	}
	if tagsJSON != "" {
		json.Unmarshal([]byte(tagsJSON), &t.Tags)
	}
	if t.Tags == nil {
		t.Tags = []string{}
	}
	subs, err := GetSubtasks(id)
	if err == nil {
		t.Subtasks = subs
	} else {
		t.Subtasks = []db.Subtask{}
	}
	return t, nil
}
