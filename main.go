package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/cgi"
	"net/http/fcgi"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	usefcgi = flag.Bool("fcgi", false, "Use Fast CGI")
	usecgi  = flag.Bool("cgi", false, "Use CGI")
	addr    = flag.String("addr", ":8080", "Bind address")
	dbname  = flag.String("dbname", "gorecall.db", "Path to the SQLite database file")
	migrate = flag.Bool("migrate", false, "Run database migrations")
	bmRepo  *bookmarkRepo
)

func main() {
	var err error

	flag.Parse()

	db, err := InitDatabase(*dbname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *migrate {
		MigrateDatabase(db)
		os.Exit(0)
	}

	// Initialize repositories
	bmRepo, err = NewBookmarkRepo(db)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/bookmarks", BookmarksHandler)
	r.HandleFunc("/bookmarks/new", BookmarksNewHandler)
	r.HandleFunc("/import", ImportHandler)
	r.HandleFunc("/login", LoginHandler)

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/bookmarks", CreateBookmarkHandler).Methods("POST")

	// Static file handler
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	if *usefcgi {
		err = fcgi.Serve(nil, r)
	} else if *usecgi {
		err = cgi.Serve(r)
	} else {
		err = http.ListenAndServe(*addr, handlers.LoggingHandler(os.Stdout, r))
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func CreateBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateBookmarkHandler")
	decoder := json.NewDecoder(r.Body)
	var b Bookmark
	err := decoder.Decode(&b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500"))
		return
	}

	_, err = bmRepo.Create(&b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500"))
		return
	}

	w.Write([]byte("bookmark created"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "login.html", "")
}

func ImportHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "import.html", "")
}

func BookmarksHandler(w http.ResponseWriter, r *http.Request) {
	bookmarks, err := bmRepo.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500"))
		return
	}
	RenderTemplate(w, "bookmarks.html", bookmarks)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "index.html", "")
}

func BookmarksNewHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "bookmarksnew.html", "")
}
