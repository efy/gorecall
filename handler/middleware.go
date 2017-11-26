package handler

import (
	"net/http"
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
)

// Wrap http timeout handler
func TimeoutMiddleware(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, 30*time.Second, "timed out")
}

// Wrap the gorilla handler for use with alice
func LoggingMiddleware(h http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, h)
}

// Set CORS header
func CORSMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		h.ServeHTTP(w, r)
	})
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

func (app *App) AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.store.Get(r, "sesh")

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
