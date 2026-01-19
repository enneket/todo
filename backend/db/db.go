package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		return err
	}

	createProjectsTableSQL := `CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		color TEXT DEFAULT '#64748B',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := DB.Exec(createProjectsTableSQL); err != nil {
		return err
	}

	createSubtasksTableSQL := `CREATE TABLE IF NOT EXISTS subtasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		todo_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		completed BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(todo_id) REFERENCES todos(id) ON DELETE CASCADE
	);`
	if _, err := DB.Exec(createSubtasksTableSQL); err != nil {
		return err
	}

	// Migrations: Try to add new columns if they don't exist
	// We ignore errors here because "duplicate column name" is expected if run multiple times
	DB.Exec(`ALTER TABLE todos ADD COLUMN priority TEXT DEFAULT 'medium'`)
	DB.Exec(`ALTER TABLE todos ADD COLUMN due_date DATETIME`)
	DB.Exec(`ALTER TABLE todos ADD COLUMN description TEXT DEFAULT ''`)
	DB.Exec(`ALTER TABLE todos ADD COLUMN tags TEXT DEFAULT '[]'`)
	DB.Exec(`ALTER TABLE todos ADD COLUMN project_id INTEGER REFERENCES projects(id) ON DELETE SET NULL`)

	return nil
}
