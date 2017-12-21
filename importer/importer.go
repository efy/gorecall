// Package importer handles persisting items to the database and generating a status report
// of successfully imported items returning a list of failed imports and their respective errors
package importer

import (
	"io"
	"log"
	"strconv"
	"sync"

	"github.com/blevesearch/bleve"
	"github.com/efy/bookmark"
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
	RunWebinfo    bool `schema:"webinfo"`
	RunIndex      bool `schema:"index"`
	FoldersAsTags bool `schema:"folders_as_tags"`
	ImportTags    bool `schema:"import_tags"`
	Concurrency   int
	Index         bleve.Index
	TagRepo       datastore.TagRepo
}

var DefaultOptions = Options{
	RunWebinfo:    false,
	Concurrency:   10,
	RunIndex:      false,
	FoldersAsTags: false,
	ImportTags:    true,
	Index:         nil,
	TagRepo:       nil,
}

// Import takes an io.Reader and handles parsing and persisting to the datastore
func Import(file io.Reader, dst datastore.BookmarkRepo, opts Options) (*Report, error) {
	var bookmarks []datastore.Bookmark
	tagMap := make(map[string][]string)

	// Set bookmark parsing options from options
	parseopts := bookmark.DefaultParseOptions
	parseopts.FoldersAsTags = opts.FoldersAsTags

	parsed, err := bookmark.ParseWithOptions(file, parseopts)
	if err != nil {
		return nil, err
	}

	// Convert from bookmark.Bookmark to datastore.Bookmark
	for _, v := range parsed {
		bookmarks = append(bookmarks, datastore.Bookmark{
			Title:   v.Title,
			URL:     v.Url,
			Icon:    v.Icon,
			Created: v.Created,
		})

		// Add to tag map
		for _, t := range v.Tags {
			tagMap[t] = append(tagMap[t], v.Url)
		}
	}

	if opts.RunWebinfo {
		bookmarks = BatchWebinfo(bookmarks, opts.Concurrency)
	}

	report := &Report{}

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

	if opts.TagRepo != nil && opts.ImportTags {
		for k, v := range tagMap {
			log.Printf("%s has %d bookmarks", k, len(v))
			tag := &datastore.Tag{
				Label: k,
			}
			tag, err := opts.TagRepo.Create(tag)
			if err != nil {
				log.Println("import:", err)
			}

			if tag, err = opts.TagRepo.GetByLabel(k); err == nil {
				for _, url := range v {
					bm, err := dst.GetByURL(url)
					if err != nil {
						continue
					}
					err = dst.AddTag(bm.ID, tag.ID)
					if err != nil {
						continue
					}
				}
			}
		}
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
