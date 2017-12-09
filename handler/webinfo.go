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
			jsonResponse(w, http.StatusBadRequest, "Missing required parameter: url")
			return
		}

		info, err := webinfo.Get(url)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		payload, err := json.Marshal(info)
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, "Failed to encode")
			return
		}

		w.Write(payload)
	})
}
