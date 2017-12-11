package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApiBookmarkHandler(t *testing.T) {
	h := muxWrapper("/api/bookmarks", mockApi.ApiBookmarksHandler())

	req, err := http.NewRequest("GET", "/api/bookmarks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error("expected %d got %d", http.StatusOK, rr.Code)
	}
}

func TestApiCreateBookmarkHandler(t *testing.T) {
	h := muxWrapper("/api/bookmarks", mockApi.ApiCreateBookmarkHandler())

	body := `{"title": "Test", "url": "http://test.com"}`

	req, err := http.NewRequest("POST", "/api/bookmarks", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Error("expected", http.StatusCreated)
		t.Error("got     ", rr.Code)
	}
}
