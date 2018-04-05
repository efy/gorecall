package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	schemasqlite = `
		CREATE TABLE bookmarks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			url TEXT UNIQUE,
			icon TEXT,
			media_type TEXT DEFAULT "",
			description TEXT DEFAULT "",
			keywords TEXT DEFAULT "",
			text_content TEXT DEFAULT "",
			status INTEGER DEFAULT 0,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(50) UNIQUE,
			password CHAR(60),
			email VARCHAR(255) UNIQUE,
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

	schemapostgres = `
		CREATE TABLE bookmarks (
			id SERIAL PRIMARY KEY,
			title TEXT,
			url TEXT UNIQUE,
			icon TEXT,
			media_type TEXT DEFAULT '',
			description TEXT DEFAULT '',
			keywords TEXT DEFAULT '',
			text_content TEXT DEFAULT '',
			status INTEGER DEFAULT 0,
			created TIMESTAMP DEFAULT NOW()
		);

		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE,
			password CHAR(60),
			email VARCHAR(255) UNIQUE,
			created TIMESTAMP DEFAULT NOW()
		);

    CREATE TABLE tags (
      id SERIAL PRIMARY KEY,
      label VARCHAR(50) UNIQUE,
			color VARCHAR(16),
			description TEXT,
			created TIMESTAMP DEFAULT NOW()
    );

		CREATE TABLE bookmark_tags (
			id SERIAL PRIMARY KEY,
			tag_id INTEGER,
			bookmark_id INTEGER,
			created TIMESTAMP DEFAULT NOW()
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

// Apply the correct database schema based on options.Driver
func Setup(opts Options, db *sqlx.DB) error {
	switch opts.Driver {
	case "sqlite3":
		_, err := db.Exec(schemasqlite)
		if err != nil {
			return err
		}
	case "postgres":
		_, err := db.Exec(schemapostgres)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("database: no schema for driver %s", opts.Driver)
	}
	return nil
}
