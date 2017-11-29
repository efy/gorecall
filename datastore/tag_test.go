package datastore

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewTagRepo(t *testing.T) {
	expect := ErrInvalidDB
	_, err := NewTagRepo(nil)

	if err != expect {
		t.Error("expected", expect, "got", err)
	}
}

func TestTagValidate(t *testing.T) {
	tt := []struct {
		tag   Tag
		errs  []error
		valid bool
	}{
		{
			Tag{},
			[]error{ErrEmptyLabel},
			false,
		},
		{
			Tag{Label: "not empty"},
			[]error{},
			true,
		},
		{
			Tag{Label: strings.Repeat("X", 51)},
			[]error{ErrLongLabel},
			false,
		},
	}

	for _, tr := range tt {
		valid, errs := tr.tag.Validate()
		if valid != tr.valid {
			t.Error("expected", tr.valid, "got", valid)
		}

		if !reflect.DeepEqual(errs, tr.errs) {
			t.Error("expected", tr.errs)
			t.Error("got     ", errs)
		}
	}
}
