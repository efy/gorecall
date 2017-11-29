package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/templates"
	"github.com/gorilla/mux"
)

func (app *App) BookmarksNewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.initAppCtx(r)
		if r.Method == "GET" {
			templates.RenderTemplate(w, "newbookmark.html", ctx)
		}

		if r.Method == "POST" {
			r.ParseForm()

			bm := datastore.Bookmark{
				Title: strings.Join(r.Form["title"], ""),
				URL:   strings.Join(r.Form["url"], ""),
			}

			_, err := app.br.Create(&bm)
			if err != nil {
				renderError(w, err)
				return
			}

			http.Redirect(w, r, "/bookmarks", 302)
		}
	})
}

func (app *App) BookmarksShowHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)

		bookmark, err := app.br.GetByID(id)
		if err != nil {
			renderError(w, err)
			return
		}

		tags, err := app.br.ListTags(id)
		if err != nil {
			renderError(w, err)
			return
		}

		templates.RenderTemplate(w, "bookmark.html", struct {
			Authenticated bool
			Bookmark      *datastore.Bookmark
			Tags          []datastore.Tag
		}{
			true,
			bookmark,
			tags,
		})
	})
}

func (app *App) BookmarksHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.initAppCtx(r)

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

		ctx.Bookmarks = bookmarks

		p := Pagination{
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

		ctx.Pagination = p

		templates.RenderTemplate(w, "bookmarks.html", ctx)
	})
}

func (app *App) CreateBookmarkHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var b datastore.Bookmark
		err := decoder.Decode(&b)
		if err != nil {
			renderError(w, err)
		}

		_, err = app.br.Create(&b)
		if err != nil {
			renderError(w, err)
			return
		}

		w.Write([]byte("bookmark created"))
	})
}
