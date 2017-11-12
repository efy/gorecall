package main

import (
	"database/sql"
	"fmt"
)

const (
	bookmarksMigration = `
    CREATE TABLE bookmarks (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      title TEXT,
      uri TEXT
    );
  `

	usersMigration = `
    CREATE TABLE users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      name VARCHAR(250),
      password_hash VARCHAR(250)
    );
  `
)

func InitDatabase(file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MigrateDatabase(db *sql.DB) {
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
