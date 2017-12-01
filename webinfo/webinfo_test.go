package webinfo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebInfoGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte("plain text"))
	}))
	defer ts.Close()

	testurl := ts.URL

	info, err := Get(testurl)
	if err != nil {
		t.Error("expected", "no error")
		t.Error("got     ", err)
	}

	if info.StatusCode != 200 {
		t.Error("expected", "200 OK")
		t.Error("got     ", info.StatusCode)
	}

	if info.StatusCode != 200 {
		t.Error("expected", "200 OK")
		t.Error("got     ", info.StatusCode)
	}

	if info.MediaType != "text/plain" {
		t.Error("expected", "text/plain")
		t.Error("got     ", info.MediaType)
	}

	if info.Size < 1 {
		t.Error("expected", ".Size to be non-zero")
		t.Error("got     ", info.Size)
	}

	if info.Ext != ".txt" {
		t.Error("expected", ".Ext to be .txt")
		t.Error("got     ", info.Ext)
	}
}

func TestWebInfoGetHtmlTitle(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.Write([]byte("<title>Test Title</title>"))
	}))
	defer ts.Close()

	testurl := ts.URL

	info, err := Get(testurl)
	if err != nil {
		t.Error("expected", "no error")
		t.Error("got     ", err)
	}
	if info.Title != "Test Title" {
		t.Error("expected", ".Title to be 'Test Title'")
		t.Error("got     ", info.Title)
	}
}

func TestWebInfoGetHtmlCoverImage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.Write([]byte(`<meta property="og:image" content="http://ia.media-imdb.com/images/rock.jpg" />`))
	}))
	defer ts.Close()

	testurl := ts.URL

	info, err := Get(testurl)
	if err != nil {
		t.Error("expected", "no error")
		t.Error("got     ", err)
	}
	if info.Cover == "" {
		t.Error("expected", ".Cover to not be blank")
		t.Error("got     ", info.Cover)
	}
}
