package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/efy/gorecall/webinfo"
)

func (app *Api) WebInfoHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			renderError(w, fmt.Errorf("no url provided"))
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
