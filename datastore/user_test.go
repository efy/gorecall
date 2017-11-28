package datastore

import "testing"

func TestNewUserRepo(t *testing.T) {
	expect := ErrInvalidDB
	_, err := NewUserRepo(nil)

	if err != expect {
		t.Error("expected", expect, "got", err)
	}
}
