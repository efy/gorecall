package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/cgi"
	"net/http/fcgi"
	"os"

	"github.com/gorilla/mux"
)

var (
	usefcgi = flag.Bool("fcgi", false, "Use Fast CGI")
	usecgi  = flag.Bool("cgi", false, "Use CGI")
	addr    = flag.String("addr", ":8080", "Bind address")
)

func main() {
	flag.Parse()

	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/bookmarks", BookmarksHandler)
	r.HandleFunc("/import", ImportHandler)
	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/add", AddHandler).Methods("POST")

	var err error

	if *usefcgi {
		err = fcgi.Serve(nil, r)
	} else if *usecgi {
		err = cgi.Serve(r)
	} else {
		err = http.ListenAndServe(*addr, r)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Bookmark struct {
	Title string
	URL   string
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var b Bookmark
	err := decoder.Decode(&b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500"))
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "login.html", "")
}

func ImportHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "import.html", "")
}

func BookmarksHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "bookmarks.html", "")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "index.html", "")
}
