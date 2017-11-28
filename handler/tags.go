package handler

import (
	"net/http"
	"strconv"

	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/templates"
)

func (app *App) NewTagHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.initAppCtx(r)
		templates.RenderTemplate(w, "newtag.html", ctx)
	})
}

func (app *App) CreateTagHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			renderError(w, err)
			return
		}

		tag := &datastore.Tag{}

		err = decoder.Decode(tag, r.PostForm)
		if err != nil {
			renderError(w, err)
			return
		}

		tag, err = app.tr.Create(tag)

		if err != nil {
			renderError(w, err)
			return
		}

		id := strconv.FormatInt(tag.ID, 10)

		http.Redirect(w, r, "/tag/"+id, 302)
		return
	})
}

func (app *App) TagHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.initAppCtx(r)
		templates.RenderTemplate(w, "tag.html", ctx)
	})
}

func (app *App) TagsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.initAppCtx(r)
		templates.RenderTemplate(w, "tags.html", ctx)
	})
}
