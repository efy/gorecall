package datastore

import "testing"

func TestNewBookmarkRepo(t *testing.T) {
	expect := ErrInvalidDB
	_, err := NewBookmarkRepo(nil)

	if err != expect {
		t.Error("expected", expect, "got", err)
	}
}
