package webinfo

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

var client = http.Client{
	Timeout: time.Second * 5,
}

var (
	ErrInvalidURL = fmt.Errorf("invalid url")
)

// Map of mime type to function for handling response body
var bodyHandlers = map[string]func(*Info, io.Reader) error{
	"text/html": func(i *Info, r io.Reader) error {
		// Create a goquery doc from reader
		doc, err := createDoc(r)
		if err != nil {
			return err
		}

		title, err := extractHtmlTitle(doc)
		if err == nil {
			i.Title = title
		}

		og, err := parseOpenGraph(doc)
		if err == nil {
			i.OpenGraph = og
		}

		image, err := extractOpenGraphImage(doc)
		if err == nil {
			i.Cover = image
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

	// HTML specific
	OpenGraph *OG    `json:"opengraph"`
	Keywords  string `json:"keywords"`
	Cover     string `json:"cover"`
}

// Get takes a URL and returns the releted information
func Get(uri string) (*Info, error) {
	info := Info{}

	url, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, ErrInvalidURL
	}

	resp, err := client.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	info.StatusCode = resp.StatusCode
	extractHeaders(&info, resp.Header)

	// Handle setting a title for non html urls
	// defaults to the file name extracted from the url
	if info.MediaType != "text/html" {
		name := path.Base(url.Path)
		if name != "." && name != "/" {
			info.Title = name
		} else {
			info.Title = "Unknown"
		}
	}

	// Extract info from body if there is a suitable
	// handler in the map
	h, ok := bodyHandlers[info.MediaType]
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
