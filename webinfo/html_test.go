package webinfo

import (
	"strings"
	"testing"
)

func TestExtractOpenGraphImage(t *testing.T) {
	body := strings.NewReader(`<meta property="og:image" content="http://ia.media-imdb.com/images/rock.jpg" />`)
	doc, err := createDoc(body)
	if err != nil {
		t.Fatal(err)
	}

	expect := "http://ia.media-imdb.com/images/rock.jpg"

	image, err := extractOpenGraphImage(doc)
	if err != nil {
		t.Error("expected", "no errors")
		t.Error("got     ", err)
	}
	if image != expect {
		t.Error("expected", expect)
		t.Error("got     ", image)
	}
}

func TestExtractTextContent(t *testing.T) {
	html := `<body><h1>a &amp; b</h1></body>`
	doc, err := createDoc(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}

	expect := "a & b"
	got, err := extractTextContent(doc)
	if err != nil {
		t.Error(err)
	}
	if got != expect {
		t.Error("expected", expect)
		t.Error("got     ", got)
	}
}

func TestDocToString(t *testing.T) {
	html := `<html><head></head><body>Hello world</body></html>`
	doc, err := createDoc(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}

	result, err := docToString(doc)
	if err != nil {
		t.Fatal(err)
	}

	if result != html {
		t.Error("expected", html)
		t.Error("got     ", result)
	}
}
