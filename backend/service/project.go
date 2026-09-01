package service

import (
	"time"

	"todo/backend/db"
)

func CreateProject(name, description, color string) (int64, error) {
	if color == "" {
		color = "#64748B"
	}
	res, err := db.DB.Exec("INSERT INTO projects (name, description, color) VALUES (?, ?, ?)", name, description, color)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetProjects() ([]db.Project, error) {
	rows, err := db.DB.Query("SELECT id, name, description, color, created_at, deleted_at FROM projects WHERE deleted_at IS NULL ORDER BY created_at ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []db.Project
	for rows.Next() {
		var p db.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Color, &p.CreatedAt, &p.DeletedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func UpdateProject(id int, name, description, color string) error {
	_, err := db.DB.Exec("UPDATE projects SET name = ?, description = ?, color = ? WHERE id = ?", name, description, color, id)
	return err
}

// DeleteProject soft-deletes the project and every one of its todos in a
// single transaction with one shared timestamp. The whole group then appears
// in the trash. Restoration is independent per item — RestoreProject brings
// the project back alone, RestoreTodo brings each todo back individually.
func DeleteProject(id int) error {
	now := time.Now().UTC()
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("UPDATE projects SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL", now, id); err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE todos SET deleted_at = ? WHERE project_id = ? AND deleted_at IS NULL", now, id); err != nil {
		return err
	}
	return tx.Commit()
}

// GetTrashedProjects lists every soft-deleted project, newest deletion first.
func GetTrashedProjects() ([]db.Project, error) {
	rows, err := db.DB.Query("SELECT id, name, description, color, created_at, deleted_at FROM projects WHERE deleted_at IS NOT NULL ORDER BY deleted_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []db.Project
	for rows.Next() {
		var p db.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Color, &p.CreatedAt, &p.DeletedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

// RestoreProject brings a trashed project back into the normal projects
// list. Its todos are left in the trash — each one keeps its own deleted_at
// and is restored individually via RestoreTodo, so restoring a project never
// moves its tasks out of the trash.
func RestoreProject(id int) error {
	_, err := db.DB.Exec("UPDATE projects SET deleted_at = NULL WHERE id = ?", id)
	return err
}

// PurgeProject permanently deletes a trashed project and its todos (whose
// subtasks are removed by the subtasks table's ON DELETE CASCADE). Because a
// project in the trash always carries every one of its todos with it, all
// rows under project_id are safe to remove here. Deleting a project that is
// not in the trash is a no-op by design.
func PurgeProject(id int) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM todos WHERE project_id = ?", id); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM projects WHERE id = ? AND deleted_at IS NOT NULL", id); err != nil {
		return err
	}
	return tx.Commit()
}
