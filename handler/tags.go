package handler

import (
	"net/http"
	"strconv"

	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/templates"
	"github.com/gorilla/mux"
)

func (app *App) NewTagHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templates.RenderTemplate(w, "newtag.html", nil)
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

		http.Redirect(w, r, "/tags/"+id, 302)
		return
	})
}

func (app *App) TagHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			renderError(w, err)
			return
		}

		opts := datastore.DefaultListOptions
		err = decoder.Decode(&opts, r.URL.Query())
		if err != nil {
			renderError(w, err)
			return
		}

		tag, err := app.tr.GetByID(id)
		if err != nil {
			renderError(w, err)
			return
		}

		count, err := app.tr.CountBookmarks(id)
		if err != nil {
			renderError(w, err)
			return
		}

		bookmarks, err := app.tr.ListBookmarks(id, opts)
		if err != nil {
			renderError(w, err)
			return
		}

		pagination := generatePagination(count, opts)

		templates.RenderTemplate(w, "tag.html", struct {
			Tag        *datastore.Tag
			Count      int
			Bookmarks  []datastore.Bookmark
			Pagination Pagination
		}{
			tag,
			count,
			bookmarks,
			pagination,
		})
	})
}

func (app *App) DeleteTagHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)

		err = app.tr.Delete(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func (app *App) TagsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tags, err := app.tr.GetAll()
		if err != nil {
			renderError(w, err)
			return
		}

		templates.RenderTemplate(w, "tags.html", struct {
			Tags []datastore.Tag
		}{
			tags,
		})
	})
}
