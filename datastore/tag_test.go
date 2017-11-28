package datastore

import "testing"

func TestNewTagRepo(t *testing.T) {
	expect := ErrInvalidDB
	_, err := NewTagRepo(nil)

	if err != expect {
		t.Error("expected", expect, "got", err)
	}
}
