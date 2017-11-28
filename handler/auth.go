package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/efy/gorecall/templates"
)

func (app *App) LogoutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := app.store.Get(r, "sesh")
		if err != nil {
			log.Println("error retrieving session")
		}
		session.Options.MaxAge = -1
		err = session.Save(r, w)
		if err != nil {
			renderError(w, err)
			return
		}
		http.Redirect(w, r, "/login", 302)
	})
}

func (app *App) LoginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewAppCtx()

		if r.Method == "GET" {
			templates.RenderTemplate(w, "login.html", ctx)
			return
		}

		name := r.FormValue("username")
		name = strings.TrimSpace(strings.ToLower(name))

		pass := r.FormValue("password")

		if name != "" && pass != "" {
			check := app.authenticate(name, pass)

			if !check {
				templates.RenderTemplate(w, "login.html", ctx)
				return
			}

			session, err := app.store.Get(r, "sesh")
			if err != nil {
				log.Println(err)
			}
			session.Values["username"] = name
			session.Values["authenticated"] = true
			session.Save(r, w)

			http.Redirect(w, r, "/", 302)
			return
		}
		templates.RenderTemplate(w, "login.html", ctx)
		return
	})
}
