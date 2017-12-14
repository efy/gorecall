package importer

import (
	"testing"

	"github.com/efy/gorecall/database"
	"github.com/efy/gorecall/datastore"
	"github.com/jmoiron/sqlx"
)

func TestImport(t *testing.T) {
	bookmarks := []datastore.Bookmark{
		{
			Title: "Bookmark 1",
			URL:   "http://bookmark1.com",
		},
		{
			Title: "Bookmark 2",
			URL:   "http://bookmark2.com",
		},
		{
			Title: "Bookmark 3",
			URL:   "http://bookmark3.com",
		},
		{
			Title: "Bookmark 3",
			URL:   "http://bookmark3.com",
		},
	}
	db := testDB()
	defer db.Close()

	bookmarkRepo, err := datastore.NewBookmarkRepo(db)
	if err != nil {
		t.Fatal(err)
	}

	report, err := Import(bookmarks, bookmarkRepo, DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}

	if report.SuccessCount != 3 {
		t.Errorf("expected %d bookmarks to import without error", 3)
		t.Errorf("got      %d bookmarks to import without error", 3)
	}

	if report.FailureCount != 1 {
		t.Errorf("expected %d bookmarks to import with error", 1)
		t.Errorf("got      %d bookmarks to import with error", 1)
	}
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
