package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/efy/bookmark"
	"github.com/efy/gorecall/auth"
	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/templates"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

var decoder = schema.NewDecoder()

// AppCtx that can be built and passed to templates
type AppCtx struct {
	Authenticated bool
	Username      string
	User          *datastore.User
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

type App struct {
	db    *sqlx.DB
	ur    datastore.UserRepo
	br    datastore.BookmarkRepo
	store *sessions.CookieStore
}

func (h App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func (app *App) CreateBookmarkHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CreateBookmarkHandler")
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

func (app *App) LoginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewAppCtx()

		if r.Method == "GET" {
			templates.RenderTemplate(w, "login.html", ctx)
			return
		}

		name := r.FormValue("username")
		name = strings.TrimSpace(strings.ToLower(name))

		pass := r.FormValue("password")

		if name != "" && pass != "" {
			check := app.authenticate(name, pass)

			fmt.Println(check)

			if !check {
				templates.RenderTemplate(w, "login.html", ctx)
				return
			}

			session, err := app.store.Get(r, "sesh")
			if err != nil {
				fmt.Println(err)
			}
			session.Values["username"] = name
			session.Values["authenticated"] = true
			session.Save(r, w)

			http.Redirect(w, r, "/", 302)
			return
		}
		templates.RenderTemplate(w, "login.html", ctx)
		return
	})
}

func (app *App) LogoutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := app.store.Get(r, "sesh")
		if err != nil {
			fmt.Println("error retrieving session")
		}
		session.Options.MaxAge = -1
		err = session.Save(r, w)
		if err != nil {
			renderError(w, err)
			return
		}
		http.Redirect(w, r, "/login", 302)
	})
}

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
				Title: v.Title,
				URL:   v.Url,
				Icon:  v.Icon,
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

func (app *App) AccountShowHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.initAppCtx(r)
		templates.RenderTemplate(w, "accountshow.html", ctx)
	})
}

func (app *App) AccountEditHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.initAppCtx(r)
		if r.Method == "POST" {
			fmt.Println("Account update not implemented")
			templates.RenderTemplate(w, "accountedit.html", ctx)
		} else {
			templates.RenderTemplate(w, "accountedit.html", ctx)
		}
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
		ctx.Bookmarks = bookmarks

		p := Pagination{
			Current: opts.Page,
			Next:    opts.Page + 1,
			Prev:    opts.Page - 1,
			Last:    200,
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

func (app *App) HomeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.initAppCtx(r)
		templates.RenderTemplate(w, "index.html", ctx)
	})
}

func (app *App) BookmarksShowHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)

		ctx := app.initAppCtx(r)

		bookmark, err := app.br.GetByID(id)
		if err != nil {
			renderError(w, err)
			return
		}
		ctx.Bookmark = bookmark

		templates.RenderTemplate(w, "bookmarksshow.html", ctx)
	})
}

func (app *App) BookmarksNewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.initAppCtx(r)
		if r.Method == "GET" {
			templates.RenderTemplate(w, "bookmarksnew.html", ctx)
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

func (app *App) NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := app.initAppCtx(r)
		templates.RenderTemplate(w, "notfound.html", ctx)
	})
}

func (app *App) ApiPingHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"version": "v0", "status": "ok"}`))
	})
}

func (app *App) PreflightHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Date, Username, Password")
	})
}

func (app *App) CreateTokenHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.Header.Get("Username")
		password := r.Header.Get("Password")

		if username == "" || password == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authentication failure, please check your credentails"))
			return
		}

		match := app.authenticate(username, password)

		if match {
			token := jwt.New(jwt.SigningMethodHS256)
			claims := token.Claims.(jwt.MapClaims)
			claims["username"] = username
			claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
			tokenString, err := token.SignedString([]byte("secret"))
			if err != nil {
				fmt.Println(err)
			}
			w.Write([]byte(tokenString))
			fmt.Println(tokenString)
			return
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authentication failure, please check your credentails"))
			return
		}
	})
}

func (app *App) ApiBookmarksHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bookmarks listing not implemented"))
		return
	})
}

func renderError(w http.ResponseWriter, err error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		templates.RenderTemplate(w, "servererror.html", err)
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

// Builds app data from the request
// TODO: move this into auth middleware
func (app *App) initAppCtx(r *http.Request) *AppCtx {
	ctx := NewAppCtx()
	session, err := app.store.Get(r, "sesh")
	if err != nil {
		fmt.Println("error retrieving session")
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
		fmt.Println(err)
	}

	ctx.User = user

	return ctx
}

func NewApp(db *sqlx.DB, ur datastore.UserRepo, br datastore.BookmarkRepo, store *sessions.CookieStore) App {
	return App{
		db,
		ur,
		br,
		store,
	}
}
