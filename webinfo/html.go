package webinfo

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrNoTitle = fmt.Errorf("could not extract title")
	ErrNoImage = fmt.Errorf("could not extract image")
)

func createDoc(body io.Reader) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func extractHtmlTitle(doc *goquery.Document) (string, error) {
	title := doc.Find("title").First().Text()
	s := strings.Trim(title, " \t\n")
	return s, nil
}

func extractOpenGraphImage(doc *goquery.Document) (string, error) {
	url, exists := doc.Find(`meta[property="og:image"]`).First().Attr("content")

	if !exists {
		return "", ErrNoImage
	}

	return url, nil
}
