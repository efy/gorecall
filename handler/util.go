package handler

import (
	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/webinfo"
)

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
