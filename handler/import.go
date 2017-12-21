package handler

import (
	"net/http"

	"github.com/efy/gorecall/importer"
	"github.com/efy/gorecall/templates"
)

func (app *App) ImportHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			templates.RenderTemplate(w, "import.html", nil)
			return
		}

		importopts := importer.DefaultOptions
		importopts.Index = app.index
		importopts.TagRepo = app.tr
		importopts.Concurrency = 150

		r.ParseMultipartForm(32 << 20)

		err := decoder.Decode(&importopts, r.PostForm)
		if err != nil {
			http.Error(w, "Could not decode import options", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("bookmarks")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		report, err := importer.Import(file, app.br, importopts)
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
