package datastore

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Bookmark struct {
	ID      int64     `db:"id"`
	Title   string    `db:"title"`
	URL     string    `db:"url"`
	Icon    string    `db:"icon"`
	Created time.Time `db:"created"`
}

const (
	bookmarkInsert = `
    INSERT INTO bookmarks (title, url, icon, created)
    VALUES (?, ?, ?, ?)
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

	tagList = `
		SELECT tags.* FROM tags
		INNER JOIN bookmark_tags
		ON tags.id = bookmark_tags.tag_id
		WHERE bookmark_tags.bookmark_id = ?
	`

	tagCount = `
		SELECT COUNT(*) FROM tags
		INNER JOIN bookmark_tags
		ON tags.id = bookmark_tags.tag_id
		WHERE bookmark_tags.bookmark_id = ?
	`

	addTag = `
		INSERT INTO bookmark_tags (bookmark_id, tag_id) VALUES (?, ?)
	`

	removeTag = `
		DELETE FROM bookmark_tags WHERE bookmark_id = ? AND tag_id = ?
	`
)

type BookmarkRepo interface {
	Create(bookmark *Bookmark) (*Bookmark, error)
	GetByID(id int64) (*Bookmark, error)
	GetAll() ([]Bookmark, error)
	List(opts ListOptions) ([]Bookmark, error)
	Count() (int, error)
	ListTags(bid int64) ([]Tag, error)
	CountTags(bid int64) (int, error)
	AddTag(bid int64, tid int64) error
	RemoveTag(bid int64, tid int64) error
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
	if bm.Created.IsZero() {
		bm.Created = time.Now()
	}
	result, err := b.db.Exec(bookmarkInsert, bm.Title, bm.URL, bm.Icon, bm.Created)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	bm, err = b.GetByID(id)
	if err != nil {
		return nil, err
	}
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

func (b *bookmarkRepo) ListTags(id int64) ([]Tag, error) {
	var tags []Tag
	if err := b.db.Select(&tags, tagList, id); err != nil {
		return tags, err
	}
	return tags, nil
}

func (b *bookmarkRepo) CountTags(id int64) (int, error) {
	var count int
	if err := b.db.Get(&count, tagCount, id); err != nil {
		return count, err
	}
	return count, nil
}

func (b *bookmarkRepo) AddTag(bid int64, tid int64) error {
	_, err := b.db.Exec(addTag, bid, tid)
	if err != nil {
		return err
	}
	return nil
}

func (b *bookmarkRepo) RemoveTag(bid int64, tid int64) error {
	_, err := b.db.Exec(removeTag, bid, tid)
	if err != nil {
		return err
	}
	return nil
}

func NewBookmarkRepo(database *sqlx.DB) (*bookmarkRepo, error) {
	if database == nil {
		return nil, ErrInvalidDB
	}
	bmr := bookmarkRepo{
		db: database,
	}
	return &bmr, nil
}
