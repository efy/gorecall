package datastore

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	ErrEmptyLabel = fmt.Errorf("tag label cannot be empty")
	ErrLongLabel  = fmt.Errorf("tag label cannot be longer than 50 characters")
)

type Tag struct {
	ID          int64     `db:"id" schema:"-"`
	Label       string    `db:"label" schema:"label"`
	Description string    `db:"description" schema:"description"`
	Color       string    `db:"color" schema:"color"`
	Created     time.Time `db:"created" schema:"-"`
}

func (t *Tag) Validate() (bool, []error) {
	valid := true
	errs := make([]error, 0)

	if t.Label == "" {
		errs = append(errs, ErrEmptyLabel)
		valid = false
	}

	if len(t.Label) > 50 {
		errs = append(errs, ErrLongLabel)
		valid = false
	}

	return valid, errs
}

type TagRepo interface {
	Create(tag *Tag) (*Tag, error)
	GetByID(id int64) (*Tag, error)
	GetByLabel(label string) (*Tag, error)
	GetAll() ([]Tag, error)
	List(opts ListOptions) ([]Tag, error)
	Count() (int, error)
	Delete(id int64) error
	ListBookmarks(tid int64, opts ListOptions) ([]Bookmark, error)
	CountBookmarks(tid int64) (int, error)
}

const (
	tagInsert          = `INSERT INTO tags (label, description, color) VALUES ($1, $2, $3)`
	tagSelectBase      = `SELECT * FROM tags`
	tagSelectCount     = `SELECT COUNT(*) FROM tags`
	tagSelectByID      = tagSelectBase + ` WHERE id = $1 LIMIT 1`
	tagSelectByLabel   = tagSelectBase + ` WHERE label = $1 LIMIT 1`
	tagListBase        = tagSelectBase + ` ORDER BY %s %s LIMIT $1 OFFSET $2 `
	tagLastInsert      = `SELECT id FROM tags ORDER BY id DESC limit 1`
	tagDelete          = `DELETE FROM tags WHERE id = $1`
	tagDeleteDependant = `DELETE FROM bookmark_tags WHERE tag_id = $1;`

	tagListBookmarks = `
		SELECT
		bookmarks.id,
		bookmarks.title,
		bookmarks.url,
		bookmarks.created
		FROM bookmarks
		INNER JOIN bookmark_tags
		ON bookmarks.id = bookmark_tags.bookmark_id
		WHERE bookmark_tags.tag_id = $1
		LIMIT $2 OFFSET $3
	`

	tagBookmarksCount = `
		SELECT COUNT(*) FROM bookmarks
		INNER JOIN bookmark_tags
		ON bookmarks.id = bookmark_tags.bookmark_id
		WHERE bookmark_tags.tag_id = $1
	`
)

type tagRepo struct {
	db *sqlx.DB
}

func (t *tagRepo) Create(tag *Tag) (*Tag, error) {
	tx, err := t.db.Beginx()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(tagInsert, tag.Label, tag.Description, tag.Color)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var id int64
	if err = tx.Get(&id, tagLastInsert); err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	tag, err = t.GetByID(id)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (t *tagRepo) GetByID(id int64) (*Tag, error) {
	tag := Tag{}
	if err := t.db.Get(&tag, tagSelectByID, id); err != nil {
		log.Println(err)
		return nil, err
	}
	return &tag, nil
}

func (t *tagRepo) GetByLabel(label string) (*Tag, error) {
	tag := Tag{}
	if err := t.db.Get(&tag, tagSelectByLabel, label); err != nil {
		log.Println(err)
		return nil, err
	}
	return &tag, nil
}

func (t *tagRepo) GetAll() ([]Tag, error) {
	var ts []Tag
	if err := t.db.Select(&ts, tagSelectBase); err != nil {
		return nil, err
	}
	return ts, nil
}

func (t *tagRepo) List(opts ListOptions) ([]Tag, error) {
	var ts []Tag
	// Potentially unsafe
	query := fmt.Sprintf(tagListBase, opts.OrderBy, opts.Order)
	offset := opts.PerPage * opts.Page
	if err := t.db.Select(&ts, query, opts.PerPage, offset); err != nil {
		return ts, err
	}
	return ts, nil
}

func (t *tagRepo) Count() (int, error) {
	var count int
	if err := t.db.Get(&count, tagSelectCount); err != nil {
		return 0, err
	}
	return count, nil
}

func (t *tagRepo) Delete(id int64) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.Exec(tagDelete, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(tagDeleteDependant, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (t *tagRepo) ListBookmarks(id int64, opts ListOptions) ([]Bookmark, error) {
	var bms []Bookmark
	offset := opts.PerPage * opts.Page
	if err := t.db.Select(&bms, tagListBookmarks, id, opts.PerPage, offset); err != nil {
		return bms, err
	}
	return bms, nil
}

func (t *tagRepo) CountBookmarks(tid int64) (int, error) {
	var count int
	if err := t.db.Get(&count, tagBookmarksCount, tid); err != nil {
		return 0, err
	}
	return count, nil
}

func NewTagRepo(database *sqlx.DB) (*tagRepo, error) {
	if database == nil {
		return nil, ErrInvalidDB
	}
	tr := tagRepo{
		db: database,
	}
	return &tr, nil
}
