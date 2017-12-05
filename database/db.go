package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	schema = `
		CREATE TABLE bookmarks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			url TEXT UNIQUE,
			icon TEXT,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(50) UNIQUE,
			password CHAR(60),
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		);

    CREATE TABLE tags (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      label VARCHAR(50) UNIQUE,
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

		CREATE UNIQUE INDEX tagging
		ON bookmark_tags (bookmark_id, tag_id);
	`
)

type Options struct {
	Driver string
	DSN    string
}

// Connect returns a new database instance using the provided options
// and ensures it is connected.
func Connect(opts Options) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	switch opts.Driver {
	case "sqlite3":
		db, err = sqlx.Open("sqlite3", opts.DSN)
		if err != nil {
			return nil, err
		}
	case "postgres":
		db, err = sqlx.Open("postgres", opts.DSN)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("database: driver not supported %s", opts.Driver)
	}

	// Test connection
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
