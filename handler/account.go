package handler

import (
	"net/http"

	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/templates"
)

func (app *App) AccountHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := app.store.Get(r, "sesh")
		if err != nil {
			renderError(w, err)
			return
		}
		userid, ok := session.Values["user_id"].(int64)
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		user, err := app.ur.GetByID(userid)
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

func (app *App) UpdateAccountHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := app.store.Get(r, "sesh")
		if err != nil {
			renderError(w, err)
			return
		}

		r.ParseForm()

		userid, ok := session.Values["user_id"].(int64)
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		user, err := app.ur.GetByID(userid)
		if err != nil {
			renderError(w, err)
			return
		}

		err = decoder.Decode(user, r.PostForm)

		if err != nil {
			renderError(w, err)
			return
		}

		user, err = app.ur.Update(user)
		if err != nil {
			renderError(w, err)
			return
		}

		http.Redirect(w, r, "/settings/account", http.StatusFound)
	})
}
