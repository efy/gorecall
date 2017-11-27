package datastore

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Bookmark struct {
	ID     int64     `db:"id"`
	Title  string    `db:"title"`
	URL    string    `db:"url"`
	Icon   string    `db:"icon"`
	Create time.Time `db:"created"`
}

const (
	bookmarkInsert = `
    INSERT INTO bookmarks (title, url, icon)
    VALUES (?, ?, ?)
  `
	bookmarkSelectBase = `
    SELECT * FROM bookmarks
  `
	bookmarkListBase = `
    SELECT * FROM bookmarks ORDER BY %s %s LIMIT ? OFFSET ?
  `

	bookmarkSelectByID = bookmarkSelectBase + `WHERE id = $1`

	bookmarkCount = `
		SELECT COUNT(*) as count FROM bookmarks
	`
)

type BookmarkRepo interface {
	Create(bookmark *Bookmark) (*Bookmark, error)
	GetByID(id int64) (*Bookmark, error)
	GetAll() ([]Bookmark, error)
	List(opts ListOptions) ([]Bookmark, error)
	Count() (int, error)
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
	result, err := b.db.Exec(bookmarkInsert, bm.Title, bm.URL, bm.Icon)
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

func (b *bookmarkRepo) Count() (int, error) {
	var count int
	if err := b.db.Get(&count, bookmarkCount); err != nil {
		return 0, err
	}
	return count, nil
}

func (b *bookmarkRepo) List(opts ListOptions) ([]Bookmark, error) {
	var bms []Bookmark
	// Potentially unsafe
	query := fmt.Sprintf(bookmarkListBase, opts.OrderBy, opts.Order)
	offset := opts.PerPage * opts.Page
	if err := b.db.Select(&bms, query, opts.PerPage, offset); err != nil {
		return bms, err
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
