package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	schema = `
		CREATE TABLE bookmarks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			url TEXT,
			icon TEXT,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(50),
			password CHAR(60),
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		);

    CREATE TABLE tags (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      label VARCHAR(50),
			color VARCHAR(16),
			description TEXT,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
    );

		CREATE TABLE bookmark_tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tag_id INTEGER,
			bookmark_id INTEGER,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`
)

// Create a new database and ensure it is connected
func Init(file string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Apply the database schema to a *sqlx.DB
func Setup(db *sqlx.DB) error {
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}
	return nil
}
