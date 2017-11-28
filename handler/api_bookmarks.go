package handler

import "net/http"

func (app *App) ApiBookmarksHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bookmarks listing not implemented"))
		return
	})
}
