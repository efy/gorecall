package webinfo

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var client = http.Client{
	Timeout: time.Second * 5,
}

var (
	ErrNoTitle = fmt.Errorf("could not extract title")
)

var handlers = map[string]func(*Info, io.Reader) error{
	"text/html": func(i *Info, r io.Reader) error {
		title, err := extractHtmlTitle(r)
		if err == nil {
			i.Title = title
		}
		return nil
	},
	"text/plain": func(i *Info, r io.Reader) error {
		return nil
	},
}

type Info struct {
	Title      string `json:"title"`
	MediaType  string `json:"media_type"`
	Size       int    `json:"size"`
	StatusCode int    `json:"status_code"`
	Ext        string `json:"ext"`
}

// Get takes a URL and returns the releted information
func Get(url string) (*Info, error) {
	info := Info{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	info.StatusCode = resp.StatusCode
	extractHeaders(&info, resp.Header)

	// Extract info from body if there is a suitable
	// handler in the map
	h, ok := handlers[info.MediaType]
	if !ok {
		return &info, nil
	}
	err = h(&info, resp.Body)
	if err != nil {
		return &info, err
	}

	return &info, nil
}

// Takes the http response headers and extracts the ones
// we're interested in adding them to the info struct
func extractHeaders(i *Info, h http.Header) {
	contentType := h.Get("content-type")
	mt, _, err := mime.ParseMediaType(contentType)
	if err == nil {
		i.MediaType = mt
	}

	ext, err := mime.ExtensionsByType(i.MediaType)
	if err == nil && len(ext) > 0 {
		i.Ext = ext[0]
	}

	contentLength := h.Get("content-length")
	size, err := strconv.ParseInt(contentLength, 10, 64)
	if err == nil {
		i.Size = int(size)
	}
}

func extractHtmlTitle(body io.Reader) (string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return "", err
	}

	s, ok := traverse(doc)
	if !ok {
		return "", ErrNoTitle
	}

	s = strings.Trim(s, " \t\n")

	return s, nil
}

func traverse(n *html.Node) (string, bool) {
	if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}
