package datastore

import (
	"reflect"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestNewTagRepo(t *testing.T) {
	expect := ErrInvalidDB
	_, err := NewTagRepo(nil)

	if err != expect {
		t.Error("expected", expect)
		t.Error("got     ", err)
	}
}

func TestTagValidate(t *testing.T) {
	tt := map[string]struct {
		tag   Tag
		errs  []error
		valid bool
	}{
		"empty label": {
			Tag{},
			[]error{ErrEmptyLabel},
			false,
		},
		"valid label": {
			Tag{Label: "not empty"},
			[]error{},
			true,
		},
		"long label": {
			Tag{Label: strings.Repeat("X", 51)},
			[]error{ErrLongLabel},
			false,
		},
	}

	for k, tr := range tt {
		t.Run(k, func(t *testing.T) {
			valid, errs := tr.tag.Validate()
			if valid != tr.valid {
				t.Error("expected", tr.valid)
				t.Error("got     ", valid)
			}

			if !reflect.DeepEqual(errs, tr.errs) {
				t.Error("expected", tr.errs)
				t.Error("got     ", errs)
			}
		})
	}
}

func TestTagRepoDelete(t *testing.T) {
	withDatabaseFixtures(t, func(db *sqlx.DB) {
		tagRepo := tagRepoTestDeps(db)

		err := tagRepo.Delete(1)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestTagRepoCreate(t *testing.T) {
	withDatabaseFixtures(t, func(db *sqlx.DB) {
		tagRepo := tagRepoTestDeps(db)

		tag := &Tag{
			Label: "Create",
		}
		tag, err := tagRepo.Create(tag)
		if err != nil {
			t.Error("expected", nil)
			t.Error("got     ", err)
		}
		if tag.Created.IsZero() {
			t.Error("expected", "Create to be set by database")
			t.Error("got     ", "Zero value date")
		}
		if tag.ID == 0 {
			t.Error("expected", "ID to be set by database")
			t.Error("got     ", "Zero value ID")
		}
	})
}

func TestTagRepoGetByID(t *testing.T) {
	withDatabaseFixtures(t, func(db *sqlx.DB) {
		tagRepo := tagRepoTestDeps(db)

		tag, err := tagRepo.GetByID(1)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		if tag.Label != "tag1" {
			t.Error("expected", "tag1")
			t.Error("got     ", tag.Label)
		}
	})
}

func TestTagRepoGetByLabel(t *testing.T) {
	withDatabaseFixtures(t, func(db *sqlx.DB) {
		tagRepo := tagRepoTestDeps(db)

		_, err := tagRepo.GetByLabel("tag1")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestTagRepoCount(t *testing.T) {
	withDatabaseFixtures(t, func(db *sqlx.DB) {
		tagRepo := tagRepoTestDeps(db)

		count, err := tagRepo.Count()
		if err != nil {
			t.Error("expected", "no error")
			t.Error("got     ", err)
		}

		if count != 3 {
			t.Error("expected", 3)
			t.Error("got     ", count)
		}
	})
}

func TestTagRepoList(t *testing.T) {
	withDatabaseFixtures(t, func(db *sqlx.DB) {
		tagRepo := tagRepoTestDeps(db)

		tags, err := tagRepo.List(DefaultListOptions)
		if err != nil {
			t.Error("expected", "no error")
			t.Error("got     ", err)
		}

		if len(tags) != 3 {
			t.Error("expected", 3)
			t.Error("got     ", len(tags))
		}
	})
}

func TestTagRepoGetAll(t *testing.T) {
	withDatabaseFixtures(t, func(db *sqlx.DB) {
		tagRepo := tagRepoTestDeps(db)

		tags, err := tagRepo.GetAll()
		if err != nil {
			t.Error("expected", "no error")
			t.Error("got     ", err)
		}

		if len(tags) != 3 {
			t.Error("expected", 3, "got", len(tags))
		}
	})
}

func TestTagRepoListBookmarks(t *testing.T) {
	withDatabaseFixtures(t, func(db *sqlx.DB) {
		tagRepo := tagRepoTestDeps(db)

		bookmarks, err := tagRepo.ListBookmarks(1, DefaultListOptions)
		if err != nil {
			t.Error("expected", "no error")
			t.Error("got     ", err)
		}

		if len(bookmarks) != 3 {
			t.Error("expected", 3)
			t.Error("got     ", len(bookmarks))
		}

		last := bookmarks[2]

		if last.Title != "bm3" {
			t.Error("expected", "bm3")
			t.Error("got     ", last.Title)
		}
	})
}

func TestTagRepoCountBookmarks(t *testing.T) {
	withDatabaseFixtures(t, func(db *sqlx.DB) {
		tagRepo := tagRepoTestDeps(db)

		count, err := tagRepo.CountBookmarks(1)
		if err != nil {
			t.Error("expected", "no error")
			t.Error("got     ", err)
		}

		if count != 3 {
			t.Error("expected", 3)
			t.Error("got     ", count)
		}
	})
}
