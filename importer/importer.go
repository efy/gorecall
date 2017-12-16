// Package importer handles persisting items to the database and generating a status report
// of successfully imported items returning a list of failed imports and their respective errors
package importer

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/blevesearch/bleve"
	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/webinfo"
)

type ItemReport struct {
	Bookmark datastore.Bookmark `json:"bookmark"`
	Error    error              `json:"error"`
}

type Report struct {
	SuccessCount int          `json:"success_count"`
	FailureCount int          `json:"failure_count"`
	Results      []ItemReport `json:"results"`
}

type Options struct {
	RunWebinfo  bool `schema:"webinfo"`
	RunIndex    bool `schema:"index"`
	Concurrency int
	Index       bleve.Index
}

var DefaultOptions = Options{
	RunWebinfo:  false,
	Concurrency: 10,
	RunIndex:    false,
	Index:       nil,
}

func Import(src []datastore.Bookmark, dst datastore.BookmarkRepo, opts Options) (*Report, error) {
	report := &Report{}
	var bookmarks []datastore.Bookmark

	if opts.RunWebinfo {
		bookmarks = BatchWebinfo(src, opts.Concurrency)
	} else {
		bookmarks = src
	}

	for _, bm := range bookmarks {
		item := ItemReport{}
		b, err := dst.Create(&bm)
		if err != nil {
			b = &bm
			item.Error = err
			report.FailureCount++
		} else {
			if opts.RunIndex && opts.Index != nil {
				id := strconv.FormatInt(b.ID, 10)
				err := opts.Index.Index(id, b)
				if err != nil {
					log.Println("failed to index bookmark id: %d\n", b.ID)
				}
			}

			report.SuccessCount++
		}

		item.Bookmark = *b

		report.Results = append(report.Results, item)
	}

	return report, nil
}

// Takes a list of bookmarks running a webinfo request for each and filling
// appropriate fields returning a channel of bookmarks
func BatchWebinfo(bms []datastore.Bookmark, cap int) []datastore.Bookmark {
	var mu sync.Mutex
	var wg sync.WaitGroup
	rateLimit := make(chan struct{}, cap)
	bookmarks := make([]datastore.Bookmark, 0)
	for _, v := range bms {
		wg.Add(1)
		rateLimit <- struct{}{}
		go func(bookmark datastore.Bookmark, rl chan struct{}) {
			defer func() { <-rl }()
			fmt.Println("running web info:", bookmark.URL)
			info, err := webinfo.Get(bookmark.URL)
			if err == nil {
				fillBookmarkFromWebinfo(&bookmark, *info)
			}
			mu.Lock()
			bookmarks = append(bookmarks, bookmark)
			mu.Unlock()
			wg.Done()
		}(v, rateLimit)
	}

	wg.Wait()
	return bookmarks
}

// Same as BatchWebinfo but handles requests in serial
func BatchWebinfoSerial(bms []datastore.Bookmark) []datastore.Bookmark {
	bookmarks := make([]datastore.Bookmark, 0)
	for _, v := range bms {
		func(bookmark datastore.Bookmark) {
			info, err := webinfo.Get(bookmark.URL)
			if err == nil {
				fillBookmarkFromWebinfo(&bookmark, *info)
			}
			bookmarks = append(bookmarks, bookmark)
		}(v)
	}

	return bookmarks
}

// Overwrites Bookmark model fields with data from webinfo request where appropriate
func fillBookmarkFromWebinfo(bm *datastore.Bookmark, info webinfo.Info) {
	if bm.Title == "" {
		bm.Title = info.Title
	}
	if bm.MediaType == "" {
		bm.MediaType = info.MediaType
	}
	if bm.Status == 0 {
		bm.Status = info.StatusCode
	}
	if bm.TextContent == "" {
		bm.TextContent = info.TextContent
	}
}
