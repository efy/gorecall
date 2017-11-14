package main

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type Bookmark struct {
	ID    int64  `db:"id"`
	Title string `db:"title"`
	URL   string `db:"url"`
}

const (
	bookmarkInsert = `
    INSERT INTO bookmarks (title, url)
    VALUES (?, ?)
  `

	bookmarkSelectBase = `
    SELECT * FROM bookmarks
  `

	bookmarkSelectByID = bookmarkSelectBase + `WHERE id = $1`
)

type BookmarkRepo interface {
	Create(bookmark *Bookmark) (*Bookmark, error)
	GetByID(id int) (*Bookmark, error)
	GetAll() ([]*Bookmark, error)
}

type bookmarkRepo struct {
	db *sqlx.DB
}

func (b *bookmarkRepo) GetByID(id int64) (*Bookmark, error) {
	bm := Bookmark{}
	if err := b.db.Get(&bm, bookmarkSelectByID, id); err != nil {
		return nil, err
	}
	return &bm, nil
}

func (b *bookmarkRepo) Create(bm *Bookmark) (*Bookmark, error) {
	result, err := b.db.Exec(bookmarkInsert, bm.Title, bm.URL)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	bm.ID = id
	return bm, nil
}

func (b *bookmarkRepo) GetAll() ([]Bookmark, error) {
	var bms []Bookmark
	if err := b.db.Select(&bms, bookmarkSelectBase); err != nil {
		return nil, err
	}
	return bms, nil
}

func NewBookmarkRepo(database *sqlx.DB) (*bookmarkRepo, error) {
	if database == nil {
		return nil, errors.New("must use valid database connection")
	}
	bmr := bookmarkRepo{
		db: database,
	}
	return &bmr, nil
}
