package handler

import (
	"net/http"
	"net/http/httptest"
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
