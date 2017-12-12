package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/efy/gorecall/datastore"
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

func TestApiCreateGetsInfo(t *testing.T) {
	requested := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requested = true
		w.Header().Set("content-type", "image/png")
		w.Write([]byte("dummy image"))
	}))
	defer ts.Close()

	url := ts.URL + "/test.png"

	body := `{"url": "` + url + `"}`

	h := muxWrapper("/api/bookmarks", mockApi.ApiCreateBookmarkHandler())

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

	if !requested {
		t.Error("expected external url to be requested")
	}

	bm := datastore.Bookmark{}
	err = json.Unmarshal(rr.Body.Bytes(), &bm)
	if err != nil {
		t.Fatal(err)
	}

	if bm.Title != "test.png" {
		t.Error("expected", "test.png")
		t.Error("got     ", bm.Title)
	}
}
