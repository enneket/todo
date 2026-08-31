package db

import (
	"database/sql"
	"log"
	"strings"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	// Enable WAL + busy_timeout via DSN pragmas:
	//   - WAL lets readers proceed concurrently with the single writer, which
	//     is what removes the "database is locked" errors under HTTP load.
	//   - busy_timeout makes any write that hits a lock wait up to 5s instead
	//     of failing immediately, smoothing over short-lived contention.
	dsn := dbPath + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"
	DB, err = sql.Open("sqlite", dsn)
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

	// Migrations: add new columns if they don't already exist.
	// "duplicate column name" errors are expected on every run after the first
	// and are silently ignored; anything else gets logged so we notice real
	// schema failures instead of discovering them at INSERT time.
	for _, stmt := range []string{
		`ALTER TABLE todos ADD COLUMN priority TEXT DEFAULT 'medium'`,
		`ALTER TABLE todos ADD COLUMN due_date DATETIME`,
		`ALTER TABLE todos ADD COLUMN remind_at DATETIME`,
		`ALTER TABLE todos ADD COLUMN notified_at DATETIME`,
		`ALTER TABLE todos ADD COLUMN repeat TEXT DEFAULT ''`,
		`ALTER TABLE todos ADD COLUMN description TEXT DEFAULT ''`,
		`ALTER TABLE todos ADD COLUMN tags TEXT DEFAULT '[]'`,
		`ALTER TABLE todos ADD COLUMN project_id INTEGER REFERENCES projects(id) ON DELETE SET NULL`,
	} {
		if _, err := DB.Exec(stmt); err != nil {
			if !isDuplicateColumnErr(err) {
				log.Printf("migration failed (%s): %v", stmt, err)
			}
		}
	}

	return nil
}

// isDuplicateColumnErr reports whether an ALTER TABLE error is the harmless
// "column already exists" case (which is expected on every run after the
// first migration).
func isDuplicateColumnErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "duplicate column") || strings.Contains(msg, "already exists")
}
