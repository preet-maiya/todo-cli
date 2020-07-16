package database

import (
	log "github.com/sirupsen/logrus"
)

func (db *DB) InitDB() error {
	createTables := `
		CREATE TABLE groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_name TEXT,
			deleted_at TEXT,
			updated_at TEXT
		);
		CREATE TABLE notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at TEXT,
			content TEXT,
			end_date TEXT,
			updated_at TEXT,
			deleted_at TEXT,
			group_id INTEGER,
			status TEXT NOT NULL DEFAULT "CREATED",
			FOREIGN KEY(group_id) REFERENCES groups(id)
		);
	`

	_, err := db.conn.Exec(createTables)
	if err != nil {
		log.Errorf("Creating tables failed: %v", err)
		return err
	}
	return nil
}
