package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

// Wrap the gorilla handler for use with alice
func LoggingMiddleware(h http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, h)
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
		authenticated := true
		if authenticated {
			// Move along
			h.ServeHTTP(w, r)
		} else {
			// Pay the troll toll
			http.Redirect(w, r, "/login", 302)
			return
		}
	})
}
