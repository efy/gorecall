package importer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/efy/gorecall/database"
	"github.com/efy/gorecall/datastore"
	"github.com/jmoiron/sqlx"
)

func TestImport(t *testing.T) {
	bookmarks := `
		<a href="http://bookmark1.com">Bookmark 1</a>
		<a href="http://bookmark2.com">Bookmark 2</a>
		<a href="http://bookmark3.com">Bookmark 3</a>
		<!-- dup -->
		<a href="http://bookmark3.com">Bookmark 3</a>
	`

	db := testDB()
	defer db.Close()

	bookmarkRepo, err := datastore.NewBookmarkRepo(db)
	if err != nil {
		t.Fatal(err)
	}

	report, err := Import(strings.NewReader(bookmarks), bookmarkRepo, DefaultOptions)
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

func TestImportFile(t *testing.T) {
	file := `
		<DT><H3>One</H3>
		<DL>
			<DT><A HREF="http://bookmark1.com">Bookmark1</A>
			<DT><H3>Two</H3>
			<DL>
				<DT><A HREF="http://bookmark2.com">Bookmark2</A>
			</DL>
			<DT><H3>Three</H3>
			<DL>
				<DT><A HREF="http://bookmark3.com">Bookmark3</A>
			</DL>
		</DL>
	`
	db := testDB()
	defer db.Close()

	bookmarkRepo, err := datastore.NewBookmarkRepo(db)
	if err != nil {
		t.Fatal(err)
	}

	tagRepo, err := datastore.NewTagRepo(db)
	if err != nil {
		t.Fatal(err)
	}

	opts := DefaultOptions
	opts.TagRepo = tagRepo
	opts.ImportTags = true
	opts.FoldersAsTags = true

	report, err := Import(strings.NewReader(file), bookmarkRepo, opts)
	if err != nil {
		t.Fatal(err)
	}

	if report.TaggingCount != 5 {
		t.Error("expected", 5)
		t.Error("got     ", report.TaggingCount)
	}

	tag, err := tagRepo.GetByLabel("Two")
	if err != nil {
		t.Fatal(err)
	}

	count, err := tagRepo.CountBookmarks(tag.ID)
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Error("expected", 1)
		t.Error("got     ", count)
	}
}

func TestBatchWebinfo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status, err := strconv.ParseInt(r.URL.Path[1:], 10, 32)
		if err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(int(status))
		fmt.Fprintln(w, "response")
	}))
	defer ts.Close()

	bookmarks := []datastore.Bookmark{
		{
			Title: "Bookmark 1",
			URL:   ts.URL + "/404",
		},
		{
			Title: "Bookmark 2",
			URL:   ts.URL + "/200",
		},
		{
			Title: "Bookmark 3",
			URL:   ts.URL + "/500",
		},
	}

	bms := BatchWebinfo(bookmarks, 5)
	if len(bms) != 3 {
		t.Error("expected", 3)
		t.Error("got     ", len(bms))
	}
}

func TestBatchWebinfoSerial(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status, err := strconv.ParseInt(r.URL.Path[1:], 10, 32)
		if err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(int(status))
		fmt.Fprintln(w, "response")
	}))
	defer ts.Close()

	bookmarks := []datastore.Bookmark{
		{
			Title: "Bookmark 1",
			URL:   ts.URL + "/404",
		},
		{
			Title: "Bookmark 2",
			URL:   ts.URL + "/200",
		},
		{
			Title: "Bookmark 3",
			URL:   ts.URL + "/500",
		},
	}

	bms := BatchWebinfoSerial(bookmarks)
	if len(bms) != 3 {
		t.Error("Expected 3 bookmarks")
	}
	if bms[0].Status != http.StatusNotFound {
		t.Error("expected", http.StatusNotFound)
		t.Error("got     ", bms[0].Status)
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
