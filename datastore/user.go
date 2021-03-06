package datastore

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID       int64     `json:"id" db:"id" schema:"id"`
	Username string    `json:"username" db:"username" schema:"username"`
	Password string    `json:"-" db:"password" schema:"-"`
	Email    string    `json:"email" db:"email" schema:"email"`
	Created  time.Time `json:"created" db:"created" schema:"created"`
}

const (
	userInsert = `INSERT INTO users (username, password, email) VALUES ($1, $2, $3)`
	userUpdate = `
		UPDATE users
		SET
		username = $1,
		email = $2
		WHERE
		id = $3
	`
	userSelectBase       = `SELECT * FROM users `
	userSelectByID       = userSelectBase + `WHERE id = $1 LIMIT 1`
	userSelectByUsername = userSelectBase + `WHERE username = $1 LIMIT 1`
	userLastInsert       = `SELECT id FROM users ORDER BY id DESC LIMIT 1`
)

type UserRepo interface {
	Create(user *User) (*User, error)
	GetByID(id int64) (*User, error)
	GetAll() ([]User, error)
	GetByUsername(username string) (*User, error)
	Update(user *User) (*User, error)
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
	tx, err := ur.db.Beginx()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(userInsert, u.Username, u.Password, u.Email)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var id int64
	if err = tx.Get(&id, userLastInsert); err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	u, err = ur.GetByID(id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (ur *userRepo) Update(user *User) (*User, error) {
	_, err := ur.db.Exec(userUpdate, user.Username, user.Email, user.ID)
	if err != nil {
		return user, err
	}
	return user, nil
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
