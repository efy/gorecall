package datastore

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Tag struct {
	ID          int64     `db:"id"`
	Label       string    `db:"label"`
	Description string    `db:"description"`
	Color       string    `db:"color"`
	Created     time.Time `db:"created"`
}

type TagRepo interface {
	Create(tag *Tag) (*Tag, error)
	GetByID(id int64) (*Tag, error)
	GetAll() ([]Tag, error)
	List(opts ListOptions) ([]Tag, error)
	Count() (int, error)
}

const (
	tagInsert = `
		INSERT INTO tags (label, description, color, created)
		VALUES (?, ?, ?, ?)
	`
	tagSelectBase  = `SELECT * FROM tags`
	tagSelectCount = `SELECT COUNT(*) FROM tags`
	tagSelectByID  = tagSelectBase + `WHERE id = $1`
	tagListBase    = tagSelectBase + `ORDER BY %s %s LIMIT ? OFFSET ?`
)

type tagRepo struct {
	db *sqlx.DB
}

func (t *tagRepo) Create(tag *Tag) (*Tag, error) {
	result, err := t.db.Exec(tagInsert, tag.Label, tag.Description, tag.Color, tag.Created)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	tag.ID = id
	return tag, nil
}

func (t *tagRepo) GetByID(id int64) (*Tag, error) {
	tag := Tag{}
	if err := t.db.Get(&tag, tagSelectByID, id); err != nil {
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

func NewTagRepo(database *sqlx.DB) (*tagRepo, error) {
	if database == nil {
		return nil, errors.New("must use valid database connection")
	}
	tr := tagRepo{
		db: database,
	}
	return &tr, nil
}
