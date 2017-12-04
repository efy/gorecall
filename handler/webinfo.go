package handler

import (
	"encoding/json"
	"net/http"

	"github.com/efy/gorecall/webinfo"
)

func (app *Api) WebInfoHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "must provide url", http.StatusBadRequest)
			return
		}

		info, err := webinfo.Get(url)
		if err != nil {
			renderError(w, err)
			return
		}

		payload, err := json.Marshal(info)
		if err != nil {
			renderError(w, err)
			return
		}

		w.Write(payload)
	})
}
