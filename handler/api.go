package handler

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func (app *Api) ApiPingHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"version": "v0", "status": "ok"}`))
	})
}

func (app *Api) CreateTokenHandler() http.Handler {
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
				log.Println(err)
			}
			w.Write([]byte(tokenString))
			return
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Authentication failure, please check your credentails"))
			return
		}
	})
}
