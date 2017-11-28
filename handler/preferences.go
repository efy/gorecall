package handler

import (
	"net/http"

	"github.com/efy/gorecall/templates"
)

func (app *App) PreferencesHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.initAppCtx(r)
		templates.RenderTemplate(w, "preferences.html", ctx)
	})
}
