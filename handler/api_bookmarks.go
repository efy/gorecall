package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/efy/gorecall/datastore"
)

func (app *Api) ApiBookmarksHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opts := datastore.DefaultListOptions
		err := decoder.Decode(&opts, r.URL.Query())
		if err != nil {
			http.Error(w, "Failed to decode request parameters", http.StatusBadRequest)
			return
		}

		bookmarks, err := app.br.List(opts)
		if err != nil {
			http.Error(w, "Error fetching bookmarks", http.StatusInternalServerError)
			return
		}

		payload, err := json.Marshal(bookmarks)
		if err != nil {
			http.Error(w, "Could not marshal json", http.StatusInternalServerError)
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
			jsonResponse(w, 400, "Could not parse body")
			return
		}

		b, err = app.br.Create(b)
		if err != nil {
			jsonResponse(w, 409, "Bookmark exists")
			return
		}

		id := strconv.FormatInt(b.ID, 10)
		err = app.index.Index(id, b)
		if err != nil {
			log.Printf("Error indexing bookmark %s", id)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "Success"}`))
	})
}

func (app *Api) ApiSearchBookmarksHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")

		q := bleve.NewMatchQuery(query)
		s := bleve.NewSearchRequest(q)
		result, err := app.index.Search(s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(payload)
	})
}
