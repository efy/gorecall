package webinfo

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/k3a/html2text"
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

// Given a html document returns the text content with html tags stripped out
func extractTextContent(doc *goquery.Document) (string, error) {
	html, err := docToString(doc)
	if err != nil {
		return "", err
	}
	plain := extractTextFromHtml(html)
	plain = strings.TrimSpace(plain)
	return plain, nil
}

func docToString(doc *goquery.Document) (string, error) {
	html, err := doc.Html()
	if err != nil {
		return "", err
	}
	return html, nil
}

// Wrapped content extraction dependency to make it easier
// to replace with custom implmentation / fork. Currently
// the dependency does not function quite as expected and
// looks to contain bugs.
//
// e.g. does not extract any content from:
// https://snook.ca/archives/html_and_css/calendar-css-grid correctly
func extractTextFromHtml(html string) string {
	return html2text.HTML2Text(html)
}
