package handler

import (
	"log"
	"net/http"

	"github.com/efy/gorecall/templates"
)

func (app *App) AccountShowHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.initAppCtx(r)
		templates.RenderTemplate(w, "account.html", ctx)
	})
}

func (app *App) AccountEditHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.initAppCtx(r)
		if r.Method == "POST" {
			log.Println("Account update not implemented")
			templates.RenderTemplate(w, "account.html", ctx)
		} else {
			templates.RenderTemplate(w, "account.html", ctx)
		}
	})
}
