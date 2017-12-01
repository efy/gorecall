package handler

import (
	"log"
	"net/http"

	"github.com/efy/gorecall/auth"
	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/templates"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// AppCtx that can be built and passed to templates
// should probably be page specific
type AppCtx struct {
	Authenticated bool
	Username      string
	User          *datastore.User
	Tag           *datastore.Tag
	Tags          []datastore.Tag
	Bookmarks     []datastore.Bookmark
	Bookmark      *datastore.Bookmark
	Pagination    Pagination
}

func NewAppCtx() *AppCtx {
	return &AppCtx{}
}

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
		ctx := app.initAppCtx(r)
		templates.RenderTemplate(w, "index.html", ctx)
	})
}

func (app *App) NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.initAppCtx(r)
		templates.RenderTemplate(w, "notfound.html", ctx)
	})
}

func renderError(w http.ResponseWriter, err error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

// Builds app data from the request
// TODO: move this into auth middleware
func (app *App) initAppCtx(r *http.Request) *AppCtx {
	ctx := NewAppCtx()
	session, err := app.store.Get(r, "sesh")
	if err != nil {
		log.Println("error retrieving session")
	}

	auth, ok := session.Values["authenticated"].(bool)
	if ok {
		ctx.Authenticated = auth
	}

	username, ok := session.Values["username"].(string)
	if ok {
		ctx.Username = username
	}

	user, err := app.ur.GetByUsername(username)
	if err != nil {
		log.Println(err)
	}

	ctx.User = user

	return ctx
}
