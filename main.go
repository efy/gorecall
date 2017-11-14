package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/cgi"
	"net/http/fcgi"
	"os"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
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
	r.Handle("/", AuthMiddleware(http.HandlerFunc(HomeHandler)))
	r.Handle("/bookmarks", AuthMiddleware(http.HandlerFunc(BookmarksHandler)))
	r.Handle("/bookmarks/new", AuthMiddleware(http.HandlerFunc(BookmarksNewHandler)))
	r.Handle("/bookmarks/{id:[0-9]+}", AuthMiddleware(http.HandlerFunc(BookmarksShowHandler)))
	r.Handle("/import", AuthMiddleware(http.HandlerFunc(ImportHandler)))

	r.Handle("/login", http.HandlerFunc(LoginHandler))
	r.Handle("/logout", http.HandlerFunc(LogoutHandler))

	r.NotFoundHandler = http.HandlerFunc(http.HandlerFunc(NotFoundHandler))

	api := r.PathPrefix("/api").Subrouter()
	api.Handle("/bookmarks", http.HandlerFunc(CreateBookmarkHandler)).Methods("POST")
	api.HandleFunc("/ping", http.HandlerFunc(ApiPingHandler)).Methods("GET")

	// Static file handler
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Build middleware chain that is run for all requests
	chain := alice.New(LoggingMiddleware).Then(r)

	if *usefcgi {
		err = fcgi.Serve(nil, r)
	} else if *usecgi {
		err = cgi.Serve(r)
	} else {
		err = http.ListenAndServe(*addr, chain)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
