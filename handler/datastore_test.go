package handler

import (
	"fmt"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/efy/gorecall/datastore"
	"github.com/gorilla/sessions"
)

// This file mocks out the datastore.* interfaces for testing handlers
// and provides a mockapp

// Uses an actual cookie store
var (
	store = sessions.NewCookieStore([]byte("test"))

	// Uses an in memory index
	mapping  = bleve.NewIndexMapping()
	index, _ = bleve.NewMemOnly(mapping)

	mockApp = App{store: store, br: &bookmarkRepo{}, ur: &userRepo{}, tr: &tagRepo{}, index: index}
	mockApi = Api{br: &bookmarkRepo{}, ur: &userRepo{}, tr: &tagRepo{}, index: index}
)

var bookmarks = []datastore.Bookmark{
	{
		ID:      1,
		Title:   "Bookmark 1",
		URL:     "http://bookmark1.com",
		Icon:    "",
		Created: time.Now(),
	},
	{
		ID:      2,
		Title:   "Bookmark 2",
		URL:     "http://bookmark2.com",
		Icon:    "",
		Created: time.Now(),
	},
}

var users = []datastore.User{
	{
		ID:       1,
		Username: "user1",
		Password: "",
	},
	{
		ID:       2,
		Username: "user2",
		Password: "",
	},
}

var tags = []datastore.Tag{
	{
		ID:          1,
		Label:       "Test label 1",
		Color:       "#000",
		Description: "",
		Created:     time.Now(),
	},
	{
		ID:          2,
		Label:       "Test label 2",
		Color:       "#555",
		Description: "",
		Created:     time.Now(),
	},
}

type userRepo struct{}

func (u *userRepo) Create(user *datastore.User) (*datastore.User, error) {
	user.ID = int64(len(users))
	user.Created = time.Now()
	return user, nil
}

func (u *userRepo) GetByID(id int64) (*datastore.User, error) {
	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("could not find user with id %d", id)
}

func (u *userRepo) GetAll() ([]datastore.User, error) {
	return users, nil
}

func (u *userRepo) GetByUsername(username string) (*datastore.User, error) {
	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("could not find user with username %s", username)
}

type tagRepo struct{}

func (t *tagRepo) Create(tag *datastore.Tag) (*datastore.Tag, error) {
	tag.ID = int64(len(tags))
	tag.Created = time.Now()
	return tag, nil
}

func (t *tagRepo) Delete(id int64) error {
	return nil
}

func (t *tagRepo) GetByID(id int64) (*datastore.Tag, error) {
	return &tags[id], nil
}

func (t *tagRepo) GetAll() ([]datastore.Tag, error) {
	return tags, nil
}

func (t *tagRepo) List(opts datastore.ListOptions) ([]datastore.Tag, error) {
	return tags, nil
}

func (t *tagRepo) Count() (int, error) {
	return len(tags), nil
}

func (t *tagRepo) ListBookmarks(tid int64, opts datastore.ListOptions) ([]datastore.Bookmark, error) {
	return bookmarks, nil
}

func (t *tagRepo) CountBookmarks(tid int64) (int, error) {
	return len(bookmarks), nil
}

type bookmarkRepo struct{}

func (b *bookmarkRepo) Create(bookmark *datastore.Bookmark) (*datastore.Bookmark, error) {
	bookmark.ID = int64(len(bookmarks))
	if bookmark.Created.IsZero() {
		bookmark.Created = time.Now()
	}
	return bookmark, nil
}

func (b *bookmarkRepo) GetByID(id int64) (*datastore.Bookmark, error) {
	if int64(len(bookmarks)) > id {
		return &bookmarks[id], nil
	}
	return nil, fmt.Errorf("bookmark not found")
}

func (b *bookmarkRepo) GetAll() ([]datastore.Bookmark, error) {
	return bookmarks, nil
}

func (b *bookmarkRepo) List(opts datastore.ListOptions) ([]datastore.Bookmark, error) {
	return bookmarks, nil
}

func (b *bookmarkRepo) Delete(id int64) error {
	if id < int64(len(bookmarks)) {
		return nil
	}
	return fmt.Errorf("cannot delete bookmark %d", id)
}

func (b *bookmarkRepo) Count() (int, error) {
	return len(bookmarks), nil
}

func (b *bookmarkRepo) ListTags(bid int64) ([]datastore.Tag, error) {
	return tags, nil
}

func (b *bookmarkRepo) CountTags(bid int64) (int, error) {
	return len(tags), nil
}

func (b *bookmarkRepo) AddTag(bid int64, tid int64) error {
	panic("not implemented")
}

func (b *bookmarkRepo) RemoveTag(bid int64, tid int64) error {
	panic("not implemented")
}
