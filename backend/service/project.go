package service

import (
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
	rows, err := db.DB.Query("SELECT id, name, description, color, created_at FROM projects ORDER BY created_at ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []db.Project
	for rows.Next() {
		var p db.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Color, &p.CreatedAt); err != nil {
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

func DeleteProject(id int) error {
	_, err := db.DB.Exec("DELETE FROM projects WHERE id = ?", id)
	return err
}
