package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	bookmarksMigration = `
    CREATE TABLE bookmarks (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      title TEXT,
      url TEXT,
      icon TEXT,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
    );
  `

	usersMigration = `
    CREATE TABLE users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      username VARCHAR(50),
      password CHAR(60),
			created DATETIME DEFAULT CURRENT_TIMESTAMP
    );
  `

	tagsMigration = `
    CREATE TABLE tags (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      label VARCHAR(50),
			color VARCHAR(16),
			description TEXT,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
    );

		CREATE TABLE bookmarks_tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tag_id INTEGER,
			bookmark_id INTEGER,
			created DATETIME DEFAULT CURRENT_TIMESTAMP
		);
  `
)

func InitDatabase(file string) (*sqlx.DB, error) {
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

func MigrateDatabase(db *sqlx.DB) {
	fmt.Println("migrating database")

	_, err := db.Exec(bookmarksMigration)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success: bookmarks migration")
	}

	_, err = db.Exec(usersMigration)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success: users migration")
	}

	_, err = db.Exec(tagsMigration)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success: tags migration")
	}

	fmt.Println("migration complete")
}
