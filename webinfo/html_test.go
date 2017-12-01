package webinfo

import (
	"strings"
	"testing"
)

func TestExtractOpenGraphImage(t *testing.T) {
	body := strings.NewReader(`<meta property="og:image" content="http://ia.media-imdb.com/images/rock.jpg" />`)
	doc, err := createDoc(body)
	if err != nil {
		t.Error(err)
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
