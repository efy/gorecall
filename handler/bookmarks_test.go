package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewBookmarkHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/bookmarks/new", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := mockApp.NewBookmarkHandler()
	h.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("expected 200 response got %d", rr.Code)
	}
}

func TestBookmarksHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/bookmarks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := mockApp.BookmarksHandler()
	h.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("expected 200 response got %d", rr.Code)
	}
}

func TestBookmarkHandler(t *testing.T) {
	h := muxWrapper("/bookmarks/{id}", mockApp.BookmarkHandler())

	req, err := http.NewRequest("GET", "/bookmarks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("expected 200 response got %d", rr.Code)
	}

	req, err = http.NewRequest("GET", "/bookmarks/1000", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	h.ServeHTTP(rr, req)
	if rr.Code != 404 {
		t.Errorf("expected 404 response got %d", rr.Code)
	}
}

func TestBookmarkAddTagHandler(t *testing.T) {
	t.Log("Test not implemented")
}

func TestBookmarkRemoveTagHandler(t *testing.T) {
	t.Log("Test not implemented")
}
