package handler

import (
	"net/http"

	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/templates"
)

func (app *App) AccountShowHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := app.store.Get(r, "sesh")
		if err != nil {
			renderError(w, err)
			return
		}
		username, ok := session.Values["username"].(string)
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		user, err := app.ur.GetByUsername(username)
		if err != nil {
			renderError(w, err)
		}

		templates.RenderTemplate(w, "account.html", struct {
			User *datastore.User
		}{
			user,
		})
	})
}

func (app *App) AccountEditHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templates.RenderTemplate(w, "account.html", nil)
	})
}
