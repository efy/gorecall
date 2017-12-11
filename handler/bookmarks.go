package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/templates"
	"github.com/gorilla/mux"
)

func (app *App) NewBookmarkHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templates.RenderTemplate(w, "newbookmark.html", nil)
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

		pagination := generatePagination(count, opts)

		templates.RenderTemplate(w, "bookmarks.html", struct {
			Bookmarks  []datastore.Bookmark
			Count      int
			Pagination Pagination
		}{
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
			Bookmark *datastore.Bookmark
			Tags     []datastore.Tag
			AllTags  []datastore.Tag
		}{
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
		err = app.index.Index(id, bm)
		if err != nil {
			log.Printf("Error indexing bookmark %s", id)
		}

		http.Redirect(w, r, "/bookmarks/"+id, 302)
	})
}

func (app *App) SearchBookmarksHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")

		q := bleve.NewMatchQuery(query)
		s := bleve.NewSearchRequest(q)
		result, err := app.index.Search(s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var bookmarks []datastore.Bookmark
		for _, i := range result.Hits {
			id, err := strconv.ParseInt(i.ID, 10, 64)
			if err != nil {
				continue
			}
			bm, err := app.br.GetByID(id)
			if err != nil {
				continue
			}
			bookmarks = append(bookmarks, *bm)
		}

		templates.RenderTemplate(w, "searchbookmarks.html", struct {
			SearchQuery  string
			SearchResult *bleve.SearchResult
			Bookmarks    []datastore.Bookmark
		}{
			query,
			result,
			bookmarks,
		})
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
