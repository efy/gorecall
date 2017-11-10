package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/fcgi"
	"os"

	"github.com/gorilla/mux"
)

var (
	usefcgi = flag.Bool("fcgi", false, "Use Fast CGI")
	addr    = flag.String("addr", ":8080", "Bind address")
)

func main() {
	flag.Parse()

	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/", HomeView)
	r.HandleFunc("/bookmarks", BookmarksView)
	r.HandleFunc("/import", ImportView)
	r.HandleFunc("/login", LoginView)

	var err error

	if *usefcgi {
		err = fcgi.Serve(nil, r)
	} else {
		err = http.ListenAndServe(*addr, r)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func LoginView(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<h1>Login</h1>")
}

func ImportView(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<h1>Import</h1>")
}

func BookmarksView(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<h1>Bookmarks</h1>")
}

func HomeView(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<h1>Recall</h1>")
}
