package handler

import (
	"net/http"

	"github.com/efy/gorecall/templates"
)

func (app *App) PreferencesHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templates.RenderTemplate(w, "preferences.html", struct {
			Authenticated bool
		}{
			true,
		})
	})
}
