package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/webinfo"
)

func (app *Api) ApiBookmarksHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opts := datastore.DefaultListOptions
		err := decoder.Decode(&opts, r.URL.Query())
		if err != nil {
			jsonResponse(w, http.StatusBadRequest, "Could not parse options")
			return
		}

		bookmarks, err := app.br.List(opts)
		if err != nil {
			jsonResponse(w, http.StatusNotFound, "No bookmarks")
			return
		}

		payload, err := json.Marshal(bookmarks)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, "Failed to marshal json")
			return
		}

		w.Write(payload)
	})
}

func (app *Api) ApiCreateBookmarkHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		b := &datastore.Bookmark{}
		err := decoder.Decode(b)
		if err != nil {
			jsonResponse(w, http.StatusBadRequest, "Could not parse options")
			return
		}

		info, err := webinfo.Get(b.URL)
		if err == nil && b.Title == "" {
			b.Title = info.Title
		}

		b, err = app.br.Create(b)
		if err != nil {
			jsonResponse(w, http.StatusConflict, "Bookmark exists")
			return
		}

		id := strconv.FormatInt(b.ID, 10)
		err = app.index.Index(id, b)
		if err != nil {
			log.Printf("Error indexing bookmark %s", id)
		}

		payload, err := json.Marshal(b)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, "Could not marshal json")
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	})
}

func (app *Api) ApiSearchBookmarksHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")

		q := bleve.NewMatchQuery(query)
		s := bleve.NewSearchRequest(q)
		result, err := app.index.Search(s)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, "Failed to execute search")
			return
		}

		var bookmarks []datastore.Bookmark
		for _, i := range result.Hits {
			id, err := strconv.ParseInt(i.ID, 10, 64)
			if err != nil {
				continue
			}
			bm, err := app.br.GetByID(id)
			if err != nil {
				continue
			}
			bookmarks = append(bookmarks, *bm)
		}

		payload, err := json.Marshal(bookmarks)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, "Failed to encode bookmarks as json")
			return
		}

		w.Write(payload)
	})
}
