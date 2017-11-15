package main

import (
	"log"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func CookieMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "sesh")

		session.Values["username"] = "shaun"
		session.Values["authenticated"] = true
		session.Save(r, w)
		h.ServeHTTP(w, r)
	})
}

// Wrap the gorilla handler for use with alice
func LoggingMiddleware(h http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, h)
}

func TokenAuthMiddleware(h http.Handler) http.Handler {
	mw := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	return mw.Handler(h)
}

// basic middleware example
func NotifyMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("before")
		defer log.Println("after")
		h.ServeHTTP(w, r)
	})
}

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "sesh")

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			// Pay the troll toll
			http.Redirect(w, r, "/login", 302)
			return
		} else {
			// Move along
			h.ServeHTTP(w, r)
		}
	})
}

func getUsername(r *http.Request) string {
	session, err := store.Get(r, "sesh")
	if err == nil {
		un, ok := session.Values["username"].(string)
		if ok && un != "" {
			return un
		}
		return ""
	}
	return ""
}
