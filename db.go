package main

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
      url TEXT
    );
  `

	usersMigration = `
    CREATE TABLE users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      username VARCHAR(50),
      password CHAR(60)
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
		fmt.Println("created table bookmarks")
	}

	_, err = db.Exec(usersMigration)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("created table user")
	}

	fmt.Println("migration complete")
}
