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
	tt := map[string]struct {
		tag   Tag
		errs  []error
		valid bool
	}{
		"empty label": {
			Tag{},
			[]error{ErrEmptyLabel},
			false,
		},
		"valid label": {
			Tag{Label: "not empty"},
			[]error{},
			true,
		},
		"long label": {
			Tag{Label: strings.Repeat("X", 51)},
			[]error{ErrLongLabel},
			false,
		},
	}

	for k, tr := range tt {
		t.Log("running test case:", k)

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
