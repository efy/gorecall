package webinfo

import (
	"io"
	"mime"
	"net/http"
	"strconv"
	"time"
)

var client = http.Client{
	Timeout: time.Second * 5,
}

var handlers = map[string]func(*Info, io.Reader) error{
	"text/html": func(i *Info, r io.Reader) error {
		// Create a doc from reader
		doc, err := createDoc(r)
		if err != nil {
			return err
		}

		title, err := extractHtmlTitle(doc)
		if err == nil {
			i.Title = title
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
	Cover      string `json:"cover"`
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
