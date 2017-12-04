package handler

import (
	"net/http"
	"strconv"

	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/templates"
	"github.com/gorilla/mux"
)

func (app *App) NewBookmarkHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templates.RenderTemplate(w, "newbookmark.html", struct {
			Authenticated bool
		}{
			true,
		})
	})
}

func (app *App) DeleteBookmarkHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)

		err = app.br.Delete(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func (app *App) BookmarksHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opts := datastore.DefaultListOptions
		err := decoder.Decode(&opts, r.URL.Query())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		bookmarks, err := app.br.List(opts)
		if err != nil {
			renderError(w, err)
			return
		}

		count, err := app.br.Count()
		if err != nil {
			renderError(w, err)
			return
		}

		pagination := Pagination{
			Current: opts.Page,
			Next:    opts.Page + 1,
			Prev:    opts.Page - 1,
			Last:    count / opts.PerPage,
			List: []int{
				opts.Page + 1,
				opts.Page + 2,
				opts.Page + 3,
				opts.Page + 4,
				opts.Page + 5,
			},
			PerPage: opts.PerPage,
		}

		templates.RenderTemplate(w, "bookmarks.html", struct {
			Authenticated bool
			Bookmarks     []datastore.Bookmark
			Count         int
			Pagination    Pagination
		}{
			true,
			bookmarks,
			count,
			pagination,
		})
	})
}

func (app *App) BookmarkHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)

		bookmark, err := app.br.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		tags, err := app.br.ListTags(id)
		if err != nil {
			renderError(w, err)
			return
		}

		allTags, err := app.tr.GetAll()
		if err != nil {
			renderError(w, err)
			return
		}

		templates.RenderTemplate(w, "bookmark.html", struct {
			Authenticated bool
			Bookmark      *datastore.Bookmark
			Tags          []datastore.Tag
			AllTags       []datastore.Tag
		}{
			true,
			bookmark,
			tags,
			allTags,
		})
	})
}

func (app *App) CreateBookmarkHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		bm := &datastore.Bookmark{}

		err := decoder.Decode(bm, r.PostForm)
		if err != nil {
			renderError(w, err)
		}

		bm, err = app.br.Create(bm)
		if err != nil {
			renderError(w, err)
			return
		}

		id := strconv.FormatInt(bm.ID, 10)

		http.Redirect(w, r, "/bookmarks/"+id, 302)
	})
}

func (app *App) BookmarkAddTagHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bid, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			renderError(w, err)
		}

		tid, err := strconv.ParseInt(r.FormValue("tag_id"), 10, 64)
		if err != nil {
			renderError(w, err)
			return
		}

		err = app.br.AddTag(bid, tid)
		if err != nil {
			renderError(w, err)
			return
		}

		http.Redirect(w, r, "/bookmarks/"+vars["id"], 302)
	})
}

func (app *App) BookmarkRemoveTagHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bid, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			renderError(w, err)
		}

		tid, err := strconv.ParseInt(r.FormValue("tag_id"), 10, 64)
		if err != nil {
			renderError(w, err)
			return
		}

		err = app.br.RemoveTag(bid, tid)
		if err != nil {
			renderError(w, err)
			return
		}

		http.Redirect(w, r, "/bookmarks/"+vars["id"], 302)
	})
}
