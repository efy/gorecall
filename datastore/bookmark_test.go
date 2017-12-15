package datastore

import (
	"testing"
)

func TestNewBookmarkRepo(t *testing.T) {
	expect := ErrInvalidDB
	_, err := NewBookmarkRepo(nil)

	if err != expect {
		t.Error("expected", expect, "got", err)
	}
}

func TestBookmarkRepoCreate(t *testing.T) {
	db, bookmarkRepo := bookmarkRepoTestDeps()
	defer db.Close()

	bm := &Bookmark{
		Title: "Create",
	}

	bm, err := bookmarkRepo.Create(bm)
	if err != nil {
		t.Error("expected", nil)
		t.Error("got     ", err)
	}
	if bm.Created.IsZero() {
		t.Error("expected", "Created to be set")
		t.Error("got     ", "Zero value date")
	}
	if bm.ID == 0 {
		t.Error("expected", "ID to be set by database")
		t.Error("got     ", "Zero value ID")
	}
	if bm.Title != "Create" {
		t.Error("expected", "Title to be Create")
		t.Error("got     ", bm.Title)
	}
}

func TestBookmarkRepoGetByID(t *testing.T) {
	db, bookmarkRepo := bookmarkRepoTestDeps()
	loadDefaultFixture(db)
	defer db.Close()

	bm, err := bookmarkRepo.GetByID(1)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if bm.Title != "bm1" {
		t.Error("expected", "bm1")
		t.Error("got     ", bm.Title)
	}
}

func TestBookmarkRepoGetAll(t *testing.T) {
	db, bookmarkRepo := bookmarkRepoTestDeps()
	loadDefaultFixture(db)
	defer db.Close()

	bms, err := bookmarkRepo.GetAll()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if len(bms) != 5 {
		t.Error("expected", "5 Bookmarks")
		t.Error("got     ", len(bms))
	}
}

func TestBookmarkRepoList(t *testing.T) {
	db, bookmarkRepo := bookmarkRepoTestDeps()
	loadDefaultFixture(db)
	defer db.Close()

	opts := ListOptions{
		PerPage: 3,
		Page:    0,
		OrderBy: "created",
		Order:   "asc",
	}

	bms, err := bookmarkRepo.List(opts)
	if err != nil {
		t.Error("expected", "no error")
		t.Error("got     ", err)
	}
	if len(bms) != 3 {
		t.Error("expected", "3 bookmarks per page")
		t.Error("got     ", len(bms))
	}
}

func TestBookmarkCount(t *testing.T) {
	db, bookmarkRepo := bookmarkRepoTestDeps()
	loadDefaultFixture(db)
	defer db.Close()

	count, err := bookmarkRepo.Count()
	if err != nil {
		t.Error("expected", "no error")
		t.Error("got     ", err)
	}

	if count != 5 {
		t.Error("expected", "5")
		t.Error("got     ", count)
	}
}

func TestBookmarkListTags(t *testing.T) {
	db, bookmarkRepo := bookmarkRepoTestDeps()
	loadDefaultFixture(db)
	defer db.Close()

	tags, err := bookmarkRepo.ListTags(1)
	if err != nil {
		t.Error("expected", "no error")
		t.Error("got     ", err)
	}
	if len(tags) != 3 {
		t.Error("expected", "bookmark to have 3 tags")
		t.Error("got     ", len(tags))
	}
	if tags[0].Label != "tag1" {
		t.Error("expected", "first tag to have label tag1")
		t.Error("got     ", tags[0].Label)
	}
}

func testBookmarkCountTags(t *testing.T) {
	db, bookmarkRepo := bookmarkRepoTestDeps()
	loadDefaultFixture(db)
	defer db.Close()

	count, err := bookmarkRepo.CountTags(2)
	if err != nil {
		t.Error("expected", "no error")
		t.Error("got     ", err)
	}
	if count != 2 {
		t.Error("expected", "count = 2")
		t.Error("got     ", count)
	}
}

func testBookmarkAddTag(t *testing.T) {
	db, bookmarkRepo := bookmarkRepoTestDeps()
	loadDefaultFixture(db)
	defer db.Close()

	err := bookmarkRepo.AddTag(5, 3)
	if err != nil {
		t.Error("expected", "no error")
		t.Error("got     ", err)
	}
}

func testBookmarkRemoveTag(t *testing.T) {
	db, bookmarkRepo := bookmarkRepoTestDeps()
	loadDefaultFixture(db)
	defer db.Close()

	err := bookmarkRepo.RemoveTag(1, 3)
	if err != nil {
		t.Error("expected", "no error")
		t.Error("got     ", err)
	}

	count, err := bookmarkRepo.CountTags(1)
	if err != nil {
		t.Error(err)
	}
	if count != 2 {
		t.Error("expected", "1 less tag")
		t.Error("got     ", count)
	}
}

func TestBookmarkDelete(t *testing.T) {
	db, bookmarkRepo := bookmarkRepoTestDeps()
	loadDefaultFixture(db)
	defer db.Close()

	err := bookmarkRepo.Delete(1)
	if err != nil {
		t.Error(err)
	}

	err = bookmarkRepo.Delete(1000)
	if err == nil {
		t.Error("expected error got nil")
	}
}

func TestBookmarkRepoUpdate(t *testing.T) {
	db, bookmarkRepo := bookmarkRepoTestDeps()
	loadDefaultFixture(db)
	defer db.Close()

	bm := &Bookmark{
		ID:    1,
		Title: "Updated bookmark",
		URL:   "http://updated.com",
	}

	bm, err := bookmarkRepo.Update(bm)
	if err != nil {
		t.Error(err)
	}

	bookmark, err := bookmarkRepo.GetByID(1)
	if err != nil {
		t.Error(err)
	}

	if bookmark.Title != "Updated bookmark" {
		t.Error("expected", "Updated bookmark")
		t.Error("got     ", bookmark.Title)
	}
}
