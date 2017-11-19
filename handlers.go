package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

// AppCtx that can be built and passed to templates
type AppCtx struct {
	Authenticated bool
	Username      string
	User          *User
	Bookmarks     []Bookmark
	Bookmark      *Bookmark
}

func NewAppCtx() *AppCtx {
	return &AppCtx{}
}

type App struct {
	db    *sqlx.DB
	ur    *userRepo
	br    *bookmarkRepo
	store *sessions.CookieStore
}

func (h App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func (app *App) CreateBookmarkHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CreateBookmarkHandler")
		decoder := json.NewDecoder(r.Body)
		var b Bookmark
		err := decoder.Decode(&b)
		if err != nil {
			renderError(w, err)
		}

		_, err = bmRepo.Create(&b)
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
			RenderTemplate(w, "login.html", ctx)
			return
		}

		name := r.FormValue("username")
		name = strings.TrimSpace(strings.ToLower(name))

		pass := r.FormValue("password")

		if name != "" && pass != "" {
			check := authenticate(name, pass)

			fmt.Println(check)

			if !check {
				RenderTemplate(w, "login.html", ctx)
				return
			}

			session, err := store.Get(r, "sesh")
			if err != nil {
				fmt.Println(err)
			}
			session.Values["username"] = name
			session.Values["authenticated"] = true
			session.Save(r, w)

			http.Redirect(w, r, "/", 302)
			return
		}
		RenderTemplate(w, "login.html", ctx)
		return
	})
}

func (app *App) LogoutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "sesh")
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
		ctx := initAppCtx(r)
		RenderTemplate(w, "import.html", ctx)
	})
}

func (app *App) AccountShowHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := initAppCtx(r)
		RenderTemplate(w, "accountshow.html", ctx)
	})
}

func (app *App) AccountEditHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := initAppCtx(r)
		if r.Method == "POST" {
			fmt.Println("Account update not implemented")
			RenderTemplate(w, "accountedit.html", ctx)
		} else {
			RenderTemplate(w, "accountedit.html", ctx)
		}
	})
}

func (app *App) BookmarksHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := initAppCtx(r)
		bookmarks, err := bmRepo.GetAll()
		if err != nil {
			renderError(w, err)
			return
		}
		ctx.Bookmarks = bookmarks

		RenderTemplate(w, "bookmarks.html", ctx)
	})
}

func (app *App) HomeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := initAppCtx(r)
		RenderTemplate(w, "index.html", ctx)
	})
}

func (app *App) BookmarksShowHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)

		ctx := initAppCtx(r)

		bookmark, err := bmRepo.GetByID(id)
		if err != nil {
			renderError(w, err)
			return
		}
		ctx.Bookmark = bookmark

		RenderTemplate(w, "bookmarksshow.html", ctx)
	})
}

func (app *App) BookmarksNewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := initAppCtx(r)
		if r.Method == "GET" {
			RenderTemplate(w, "bookmarksnew.html", ctx)
		}

		if r.Method == "POST" {
			r.ParseForm()

			bm := Bookmark{
				Title: strings.Join(r.Form["title"], ""),
				URL:   strings.Join(r.Form["url"], ""),
			}

			_, err := bmRepo.Create(&bm)
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
		ctx := initAppCtx(r)
		RenderTemplate(w, "notfound.html", ctx)
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

		match := authenticate(username, password)

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
		RenderTemplate(w, "servererror.html", err)
	}
}

func authenticate(username string, password string) bool {
	u, err := uRepo.GetByUsername(username)

	if err != nil {
		return false
	}

	match := CheckPasswordHash(password, u.Password)
	if !match {
		return false
	}

	return true
}

// Builds app data from the request
// TODO: move this into auth middleware
func initAppCtx(r *http.Request) *AppCtx {
	ctx := NewAppCtx()
	session, err := store.Get(r, "sesh")
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

	user, err := uRepo.GetByUsername(username)
	if err != nil {
		fmt.Println(err)
	}

	ctx.User = user

	return ctx
}
