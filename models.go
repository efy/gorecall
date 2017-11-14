package main

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

// Bookmarks

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

type User struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

const (
	userInsert = `
    INSERT INTO users (username, password)
    VALUES (?, ?)
  `
	userSelectBase = `
    SELECT username, password FROM users
  `
	userSelectByID       = userSelectBase + `WHERE id = ? LIMIT 1`
	userSelectByUsername = userSelectBase + `WHERE username = ? LIMIT 1`
)

type UserRepo interface {
	Create(user *User) (*User, error)
	GetByID(id int) (*User, error)
	GetAll() ([]*User, error)
	GetByUsername(username string) (*User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func (ur *userRepo) GetByID(id int64) (*User, error) {
	u := User{}
	if err := ur.db.Get(&u, userSelectByID, id); err != nil {
		return nil, err
	}
	return &u, nil
}

func (ur *userRepo) Create(u *User) (*User, error) {
	result, err := ur.db.Exec(userInsert, u.Username, u.Password)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	u.ID = id
	return u, nil
}

func (ur *userRepo) GetAll() ([]User, error) {
	var us []User
	if err := ur.db.Select(&us, userSelectBase); err != nil {
		return nil, err
	}
	return us, nil
}

func (ur *userRepo) GetByUsername(username string) (*User, error) {
	u := User{}
	if err := ur.db.Get(&u, userSelectByUsername, username); err != nil {
		return nil, err
	}
	return &u, nil
}

func NewUserRepo(database *sqlx.DB) (*userRepo, error) {
	if database == nil {
		return nil, errors.New("must use valid database connection")
	}
	ur := userRepo{
		db: database,
	}
	return &ur, nil
}
