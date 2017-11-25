package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"net/http/cgi"
	"net/http/fcgi"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
)

var (
	addr    = flag.String("addr", ":8080", "Bind address")
	dbname  = flag.String("dbname", "gorecall.db", "Path to the SQLite database file")
	usefcgi = flag.Bool("fcgi", false, "Use Fast CGI")
	usecgi  = flag.Bool("cgi", false, "Use CGI")

	// Command flags
	migrate    = flag.Bool("migrate", false, "Run database migrations")
	createuser = flag.Bool("createuser", false, "Create a user")

	bmRepo *bookmarkRepo
	uRepo  *userRepo
	store  *sessions.CookieStore
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

	uRepo, err = NewUserRepo(db)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *createuser {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Creating user...")

		fmt.Println("Enter username:")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)
		user, err := uRepo.GetByUsername(username)
		if err != nil && user != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if user != nil {
			fmt.Println("User already exists")
			os.Exit(1)
		}

		fmt.Println("Enter password:")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		if len(password) < 8 {
			fmt.Println("Password should be at least 8 characters long")
			os.Exit(1)
		}

		hash, err := HashPassword(password)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		u := &User{
			Username: username,
			Password: hash,
		}

		u, err = uRepo.Create(u)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("created user:", u.Username)
		os.Exit(0)
	}

	// Initize the cookie store
	store = sessions.NewCookieStore([]byte("something-very-secret"))

	app := App{
		db:    db,
		ur:    uRepo,
		br:    bmRepo,
		store: store,
	}

	r := mux.NewRouter()
	r.Handle("/", AuthMiddleware(app.HomeHandler()))
	r.Handle("/bookmarks", AuthMiddleware(app.BookmarksHandler()))
	r.Handle("/bookmarks/new", AuthMiddleware(app.BookmarksNewHandler()))
	r.Handle("/bookmarks/{id:[0-9]+}", AuthMiddleware(app.BookmarksShowHandler()))
	r.Handle("/import", AuthMiddleware(app.ImportHandler())).Methods("GET", "POST")
	r.Handle("/account", AuthMiddleware(app.AccountShowHandler()))
	r.Handle("/account/edit", AuthMiddleware(app.AccountEditHandler()))

	r.Handle("/login", app.LoginHandler())
	r.Handle("/logout", app.LogoutHandler())

	r.NotFoundHandler = app.NotFoundHandler()

	api := r.PathPrefix("/api").Subrouter()
	api.Handle("/bookmarks", TokenAuthMiddleware(app.CreateBookmarkHandler())).Methods("POST")
	api.Handle("/bookmarks", TokenAuthMiddleware(app.ApiBookmarksHandler())).Methods("GET")
	api.Handle("/ping", CORSMiddleware(app.ApiPingHandler())).Methods("GET")
	api.Handle("/auth", CORSMiddleware(app.CreateTokenHandler())).Methods("POST")
	api.Handle("/auth", app.PreflightHandler()).Methods("OPTIONS")

	// Static file handler
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Build middleware chain that is run for all requests
	chain := alice.New(LoggingMiddleware, TimeoutMiddleware).Then(r)

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
