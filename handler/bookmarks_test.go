package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCreateBookmarkHandler(t *testing.T) {
	form := url.Values{}

	form.Add("title", "test create bookmark handler")
	form.Add("url", "http://testcreatebookmarkhandler.com")

	req := &http.Request{
		Method:   "POST",
		URL:      &url.URL{Path: "/bookmarks"},
		PostForm: form,
	}

	rr := httptest.NewRecorder()
	h := mockApp.CreateBookmarkHandler()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Errorf("expected %d response got %d", http.StatusFound, rr.Code)
	}
}

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

func TestDeleteBookmarkHandler(t *testing.T) {
	h := muxWrapper("/bookmarks/{id}", mockApp.DeleteBookmarkHandler())

	req, err := http.NewRequest("DELETE", "/bookmarks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected 204 response got %d", rr.Code)
	}

	req, err = http.NewRequest("DELETE", "/bookmarks/1000", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 response got %d", rr.Code)
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

func TestSearchBookmarksHandler(t *testing.T) {
	h := mockApp.SearchBookmarksHandler()
	req, err := http.NewRequest("GET", "/bookmarks/search", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 got %d", rr.Code)
	}
}

func TestCreateGetsInfo(t *testing.T) {
	requested := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requested = true
		w.Header().Set("content-type", "image/png")
		w.Write([]byte("dummy image"))
	}))
	defer ts.Close()

	form := url.Values{}
	form.Add("url", ts.URL)

	req := &http.Request{
		Method:   "POST",
		URL:      &url.URL{Path: "/bookmarks"},
		PostForm: form,
	}

	rr := httptest.NewRecorder()
	h := mockApp.CreateBookmarkHandler()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Error("expected", http.StatusFound)
		t.Error("got     ", rr.Code)
	}

	if !requested {
		t.Error("expected external url to be requested")
	}
}

func TestBookmarkAddTagHandler(t *testing.T) {
	t.Log("Test not implemented")
}

func TestBookmarkRemoveTagHandler(t *testing.T) {
	t.Log("Test not implemented")
}
