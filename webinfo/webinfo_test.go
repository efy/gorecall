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

	if info.MediaType != "text/plain" {
		t.Error("expected", "text/plain")
		t.Error("got     ", info.MediaType)
	}

	if info.Size < 1 {
		t.Error("expected", ".Size to be non-zero")
		t.Error("got     ", info.Size)
	}

	if info.Ext != ".txt" || info.Ext != ".asc" {
		t.Error("expected", ".Ext to be .txt or .asc")
		t.Error("got     ", info.Ext)
	}
}

func TestWebInfoExtractsTextContent(t *testing.T) {
	html := `
			<html>
				<head>
					<title>Test</title>
				</head>
				<body>
					<h2>Title</h2>
				</body>
			</html>`

	expect := `Title`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.Write([]byte(html))
	}))
	defer ts.Close()

	testurl := ts.URL

	info, err := Get(testurl)
	if err != nil {
		t.Error("expected", "no error")
		t.Error("got     ", err)
	}

	if info.TextContent != expect {
		t.Error("expected", expect)
		t.Error("got     ", info.TextContent)
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

func TestWebInfoGetImageTitle(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "image/png")
		w.Write([]byte("imagedata"))
	}))
	defer ts.Close()

	info, err := Get(ts.URL + "/test-file.png")
	if err != nil {
		t.Fatal(err)
	}

	if info.Title != "test-file.png" {
		t.Error("expected", "test-file.png")
		t.Error("got     ", info.Title)
	}

	info, err = Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if info.Title != "Unknown" {
		t.Error("expected", "Unknown")
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

func TestInvalidUrl(t *testing.T) {
	_, err := Get("not a url")
	if err != ErrInvalidURL {
		t.Error("expected", ErrInvalidURL)
		t.Error("got     ", err)
	}
}

func TestOpenGraphParsing(t *testing.T) {
	html := `
		<title>The Rock (2006)</title>
		<meta property="og:title" content="The Rock" />
		<meta property="og:type" content="video.movie" />
		<meta property="og:url" content="http://www.imdb.com/title/tt0117500/" />
		<meta property="og:image" content="http://ia.media-imdb.com/images/rock.jpg" />
	`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	}))
	defer ts.Close()

	info, err := Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if info.OpenGraph == nil {
		t.Error("expected open graph tags to be parsed")
	}
}

func TestPartialOpenGraphFails(t *testing.T) {
	html := `
		<title>The Rock (2006)</title>
		<meta property="og:title" content="The Rock" />
	`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	}))
	defer ts.Close()

	info, err := Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	if info.OpenGraph != nil {
		t.Error("expected open graph field to be nil when required attributes are missing")
	}
}
