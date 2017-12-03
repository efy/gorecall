package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestCreateTagHandler(t *testing.T) {
	form := url.Values{}

	form.Add("title", "test create handler")
	form.Add("url", "http://testcreatehandler.com")

	req, err := http.NewRequest("POST", "/tags", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := mockApp.CreateTagHandler()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Error("expected", http.StatusFound)
		t.Error("got     ", rr.Code)
	}
}

func TestNewTagHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/tags/new", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := mockApp.NewTagHandler()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error("expected", http.StatusOK)
		t.Error("got     ", rr.Code)
	}
}

func TestTagHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/tags/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := muxWrapper("/tags/{id}", mockApp.TagHandler())
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error("expected", http.StatusOK)
		t.Error("got		 ", rr.Code)
	}
}

func TestTagsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/tags", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := mockApp.TagsHandler()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error("expected", http.StatusOK)
		t.Error("got		 ", rr.Code)
	}
}

// Mux wrapper takes the path and handler and returns
// a new mux router with the handler mounted. This is
// a work around to test handlers that use url params.
func muxWrapper(path string, h http.Handler) http.Handler {
	m := mux.NewRouter()
	m.Handle(path, h)
	return m
}
