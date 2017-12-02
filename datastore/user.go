package datastore

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID       int64     `db:"id" schema:"id"`
	Username string    `db:"username" schema:"username"`
	Password string    `db:"password" schema:"-"`
	Created  time.Time `db:"created" schema:"created"`
}

const (
	userInsert = `
    INSERT INTO users (username, password)
    VALUES (?, ?)
  `
	userSelectBase = `
    SELECT * FROM users
  `
	userSelectByID       = userSelectBase + `WHERE id = ? LIMIT 1`
	userSelectByUsername = userSelectBase + `WHERE username = ? LIMIT 1`
)

type UserRepo interface {
	Create(user *User) (*User, error)
	GetByID(id int64) (*User, error)
	GetAll() ([]User, error)
	GetByUsername(username string) (*User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func (ur *userRepo) GetByID(id int64) (*User, error) {
	u := &User{}
	if err := ur.db.Get(u, userSelectByID, id); err != nil {
		return nil, err
	}
	return u, nil
}

func (ur *userRepo) Create(u *User) (*User, error) {
	result, err := ur.db.Exec(userInsert, u.Username, u.Password)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	u, err = ur.GetByID(id)
	if err != nil {
		return nil, err
	}
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
		return nil, ErrInvalidDB
	}
	ur := userRepo{
		db: database,
	}
	return &ur, nil
}
