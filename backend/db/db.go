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

	// Migrations: Try to add new columns if they don't exist
	// We ignore errors here because "duplicate column name" is expected if run multiple times
	DB.Exec(`ALTER TABLE todos ADD COLUMN priority TEXT DEFAULT 'medium'`)
	DB.Exec(`ALTER TABLE todos ADD COLUMN due_date DATETIME`)

	return nil
}
