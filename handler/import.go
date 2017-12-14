package handler

import (
	"net/http"

	"github.com/efy/bookmark"
	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/importer"
	"github.com/efy/gorecall/templates"
)

func (app *App) ImportHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			templates.RenderTemplate(w, "import.html", nil)
			return
		}
		r.ParseMultipartForm(32 << 20)

		file, _, err := r.FormFile("bookmarks")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		parsed, err := bookmark.Parse(file)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		bookmarks := make([]datastore.Bookmark, 0)

		// Convert from bookmark.Bookmark to datastore.Bookmark and populate ctx.Bookmarks
		for _, v := range parsed {
			bookmarks = append(bookmarks, datastore.Bookmark{
				Title:   v.Title,
				URL:     v.Url,
				Icon:    v.Icon,
				Created: v.Created,
			})
		}

		report, err := importer.Import(bookmarks, app.br, importer.DefaultOptions)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		templates.RenderTemplate(w, "importsuccess.html", struct {
			Report importer.Report
		}{
			*report,
		})
	})
}
