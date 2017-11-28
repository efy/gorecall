package handler

import (
	"log"
	"net/http"

	"github.com/efy/bookmark"
	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/templates"
)

func (app *App) ImportHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.initAppCtx(r)
		if r.Method != "POST" {
			templates.RenderTemplate(w, "import.html", ctx)
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

		// Convert from bookmark.Bookmark to datastore.Bookmark and populate ctx.Bookmarks
		for _, v := range parsed {
			ctx.Bookmarks = append(ctx.Bookmarks, datastore.Bookmark{
				Title:   v.Title,
				URL:     v.Url,
				Icon:    v.Icon,
				Created: v.Created,
			})
		}

		for _, bm := range ctx.Bookmarks {
			_, err := app.br.Create(&bm)
			if err != nil {
				log.Println(err)
			}
		}

		templates.RenderTemplate(w, "importsuccess.html", ctx)
	})
}
