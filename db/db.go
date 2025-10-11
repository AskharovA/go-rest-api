package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbName)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTables(db *sql.DB) error {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)
	`
	_, err := db.Exec(createUsersTable)
	if err != nil {
		return err
	}

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		userId INTEGER,
		FOREIGN KEY(userId) REFERENCES users(id)
	)
	`
	_, err = db.Exec(createEventsTable)
	if err != nil {
		return err
	}

	createRegistrationsTable := `
	CREATE TABLE IF NOT EXISTS registrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		eventId INTEGER,
		userId INTEGER,
		FOREIGN KEY(eventId) REFERENCES events(id),
		FOREIGN KEY(userId) REFERENCES users(id)
	)
	`
	_, err = db.Exec(createRegistrationsTable)

	return err
}
