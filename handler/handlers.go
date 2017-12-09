package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/efy/gorecall/auth"
	"github.com/efy/gorecall/templates"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

type Pagination struct {
	Current int
	Next    int
	Prev    int
	Last    int
	List    []int
	PerPage int
}

func (app *App) HomeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/bookmarks", http.StatusFound)
	})
}

func (app *App) NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templates.RenderTemplate(w, "notfound.html", nil)
	})
}

func renderError(w http.ResponseWriter, err error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type apiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Helper function for rendering a standard api reponses
func jsonResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	resp := apiResponse{
		status,
		message,
	}
	payload, err := json.Marshal(resp)
	if err != nil {
		panic("cannot marshal apiResponse into json")
	}
	log.Println(string(payload))
	w.Write(payload)
}

func (app *App) authenticate(username string, password string) bool {
	u, err := app.ur.GetByUsername(username)

	if err != nil {
		return false
	}

	match := auth.CheckPasswordHash(password, u.Password)
	if !match {
		return false
	}

	return true
}

func (app *Api) authenticate(username string, password string) bool {
	u, err := app.ur.GetByUsername(username)

	if err != nil {
		return false
	}

	match := auth.CheckPasswordHash(password, u.Password)
	if !match {
		return false
	}

	return true
}
