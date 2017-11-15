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
)

func CreateBookmarkHandler(w http.ResponseWriter, r *http.Request) {
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
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		RenderTemplate(w, "login.html", "")
		return
	}

	name := r.FormValue("username")
	name = strings.TrimSpace(strings.ToLower(name))

	pass := r.FormValue("password")

	if name != "" && pass != "" {
		check := authenticate(name, pass)

		fmt.Println(check)

		if !check {
			RenderTemplate(w, "login.html", "")
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
	RenderTemplate(w, "login.html", "")
	return
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sesh")
	if err != nil {
		fmt.Println("error getting session")
	}
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		renderError(w, err)
		return
	}
	http.Redirect(w, r, "/login", 302)
}

func ImportHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "import.html", "")
}

func BookmarksHandler(w http.ResponseWriter, r *http.Request) {
	bookmarks, err := bmRepo.GetAll()
	if err != nil {
		renderError(w, err)
		return
	}
	RenderTemplate(w, "bookmarks.html", bookmarks)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "index.html", "")
}

func BookmarksShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)

	bookmark, err := bmRepo.GetByID(id)
	if err != nil {
		renderError(w, err)
		return
	}
	RenderTemplate(w, "bookmarksshow.html", bookmark)
}

func BookmarksNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		RenderTemplate(w, "bookmarksnew.html", "")
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
}

func ApiPingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"version": "v0", "status": "ok"}`))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "notfound.html", "")
}

func PreflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Date, Username, Password")
}

func CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
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
}

func ApiBookmarksHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bookmarks listing not implemented"))
	return
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
