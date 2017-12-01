package handler

import (
	"net/http"

	"github.com/efy/gorecall/templates"
)

func (app *App) AccountShowHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templates.RenderTemplate(w, "account.html", struct {
			Authenticated bool
		}{
			true,
		})
	})
}

func (app *App) AccountEditHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templates.RenderTemplate(w, "account.html", struct {
			Authenticated bool
		}{
			true,
		})
	})
}
