package datastore

import (
	"github.com/efy/gorecall/database"
	"github.com/jmoiron/sqlx"
)

// Returns dependencies required for testing
// tag repo
func tagRepoTestDeps() (*sqlx.DB, *tagRepo) {
	db := testDB()
	tagRepo, err := NewTagRepo(db)
	if err != nil {
		panic(err)
	}
	return db, tagRepo
}

// Returns dependencies required for testing
// bookmark repo
func bookmarkRepoTestDeps() (*sqlx.DB, *bookmarkRepo) {
	db := testDB()
	bookmarkRepo, err := NewBookmarkRepo(db)
	if err != nil {
		panic(err)
	}
	return db, bookmarkRepo
}

// Returns dependencies required for testing
// user repo
func userRepoTestDeps() (*sqlx.DB, *userRepo) {
	db := testDB()
	userRepo, err := NewUserRepo(db)
	if err != nil {
		panic(err)
	}
	return db, userRepo
}

// Returns in memory database with schema applied
func testDB() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	database.Setup(database.Options{Driver: "sqlite3"}, db)

	return db
}

// Fill the database with test data
func loadDefaultFixture(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec(tx.Rebind("INSERT INTO tags (label, color, description) VALUES(?, ? ,?)"), "tag1", "#000", "")
	tx.MustExec(tx.Rebind("INSERT INTO tags (label, color, description) VALUES(?, ? ,?)"), "tag2", "#000", "")
	tx.MustExec(tx.Rebind("INSERT INTO tags (label, color, description) VALUES(?, ? ,?)"), "tag3", "#000", "")

	tx.MustExec(tx.Rebind("INSERT INTO bookmarks (title, url, icon) VALUES(?, ?, ?)"), "bm1", "bmurl1", "")
	tx.MustExec(tx.Rebind("INSERT INTO bookmarks (title, url, icon) VALUES(?, ?, ?)"), "bm2", "bmurl2", "")
	tx.MustExec(tx.Rebind("INSERT INTO bookmarks (title, url, icon) VALUES(?, ?, ?)"), "bm3", "bmurl3", "")
	tx.MustExec(tx.Rebind("INSERT INTO bookmarks (title, url, icon) VALUES(?, ?, ?)"), "bm4", "bmurl4", "")
	tx.MustExec(tx.Rebind("INSERT INTO bookmarks (title, url, icon) VALUES(?, ?, ?)"), "bm5", "bmurl5", "")

	// Populate join table
	// bm1 tags: 1, 2, 3
	tx.MustExec(tx.Rebind("INSERT INTO bookmark_tags (bookmark_id, tag_id) VALUES(?, ?)"), 1, 1)
	tx.MustExec(tx.Rebind("INSERT INTO bookmark_tags (bookmark_id, tag_id) VALUES(?, ?)"), 1, 2)
	tx.MustExec(tx.Rebind("INSERT INTO bookmark_tags (bookmark_id, tag_id) VALUES(?, ?)"), 1, 3)

	// bm2 tags: 1, 3
	tx.MustExec(tx.Rebind("INSERT INTO bookmark_tags (bookmark_id, tag_id) VALUES(?, ?)"), 2, 1)
	tx.MustExec(tx.Rebind("INSERT INTO bookmark_tags (bookmark_id, tag_id) VALUES(?, ?)"), 2, 3)

	// bm3 tags: 1
	tx.MustExec(tx.Rebind("INSERT INTO bookmark_tags (bookmark_id, tag_id) VALUES(?, ?)"), 3, 1)

	// users
	tx.MustExec(tx.Rebind("INSERT INTO users (username, password, email) VALUES(?, ?, ?)"), "testuser1", "testpw", "testuser1@test.com")
	tx.MustExec(tx.Rebind("INSERT INTO users (username, password, email) VALUES(?, ?, ?)"), "testuser2", "testpw", "testuser2@test.com")

	tx.Commit()
}
