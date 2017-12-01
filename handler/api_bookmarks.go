package handler

import (
	"encoding/json"
	"net/http"

	"github.com/efy/gorecall/datastore"
)

func (app *App) ApiBookmarksHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bookmarks listing not implemented"))
		return
	})
}

func (app *App) ApiCreateBookmarkHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		b := &datastore.Bookmark{}
		err := decoder.Decode(b)
		if err != nil {
			renderError(w, err)
			return
		}

		_, err = app.br.Create(b)
		if err != nil {
			renderError(w, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "Success"}`))
	})
}
