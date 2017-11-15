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
	"github.com/jmoiron/sqlx"
)

// AppData that can be built and passed to templates
type AppData struct {
	Authenticated bool
	Username      string
	User          *User
	Bookmarks     []Bookmark
	Bookmark      *Bookmark
}

type AppHandler struct {
	db   *sqlx.DB
	Data *AppData
}

func (h AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func CreateBookmarkHandler(app AppHandler) http.Handler {
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

func LoginHandler(app AppHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			RenderTemplate(w, "login.html", app.Data)
			return
		}

		name := r.FormValue("username")
		name = strings.TrimSpace(strings.ToLower(name))

		pass := r.FormValue("password")

		if name != "" && pass != "" {
			check := authenticate(name, pass)

			fmt.Println(check)

			if !check {
				RenderTemplate(w, "login.html", app.Data)
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
		RenderTemplate(w, "login.html", app.Data)
		return
	})
}

func LogoutHandler(app AppHandler) http.Handler {
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

func ImportHandler(app AppHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initAppData(r, app.Data)
		RenderTemplate(w, "import.html", app.Data)
	})
}

func AccountShowHandler(app AppHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initAppData(r, app.Data)
		RenderTemplate(w, "accountshow.html", app.Data)
	})
}

func AccountEditHandler(app AppHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initAppData(r, app.Data)
		if r.Method == "POST" {
			fmt.Println("Account update not implemented")
			RenderTemplate(w, "accountedit.html", app.Data)
		} else {
			RenderTemplate(w, "accountedit.html", app.Data)
		}
	})
}

func BookmarksHandler(app AppHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initAppData(r, app.Data)
		bookmarks, err := bmRepo.GetAll()
		if err != nil {
			renderError(w, err)
			return
		}
		app.Data.Bookmarks = bookmarks

		RenderTemplate(w, "bookmarks.html", app.Data)
	})
}

func HomeHandler(app AppHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initAppData(r, app.Data)
		RenderTemplate(w, "index.html", app.Data)
	})
}

func BookmarksShowHandler(app AppHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)

		initAppData(r, app.Data)

		bookmark, err := bmRepo.GetByID(id)
		if err != nil {
			renderError(w, err)
			return
		}
		app.Data.Bookmark = bookmark

		RenderTemplate(w, "bookmarksshow.html", app.Data)
	})
}

func BookmarksNewHandler(app AppHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initAppData(r, app.Data)
		if r.Method == "GET" {
			RenderTemplate(w, "bookmarksnew.html", app.Data)
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

func NotFoundHandler(app AppHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initAppData(r, app.Data)
		RenderTemplate(w, "notfound.html", app.Data)
	})
}

func ApiPingHandler(app AppHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"version": "v0", "status": "ok"}`))
	})
}

func PreflightHandler(app AppHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Date, Username, Password")
	})
}

func CreateTokenHandler(app AppHandler) http.Handler {
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

func ApiBookmarksHandler(app AppHandler) http.Handler {
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

// Helper to populate app data object from the session
// TODO: move this into auth middleware
func initAppData(r *http.Request, data *AppData) {
	session, err := store.Get(r, "sesh")
	if err != nil {
		fmt.Println("error retrieving session")
	}

	auth, ok := session.Values["authenticated"].(bool)
	if ok {
		data.Authenticated = auth
	}

	username, ok := session.Values["username"].(string)
	if ok {
		data.Username = username
	}

	user, err := uRepo.GetByUsername(username)
	if err != nil {
		fmt.Println(err)
	}

	data.User = user
}
