package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/efy/gorecall/auth"
	"github.com/efy/gorecall/database"
	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/handler"
	"github.com/efy/gorecall/server"
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

	bmRepo datastore.BookmarkRepo
	uRepo  datastore.UserRepo
	store  *sessions.CookieStore
)

func main() {
	var err error

	flag.Parse()

	db, err := database.InitDatabase(*dbname)
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
		database.MigrateDatabase(db)
		os.Exit(0)
	}

	// Initialize repositories
	bmRepo, err = datastore.NewBookmarkRepo(db)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	uRepo, err = datastore.NewUserRepo(db)
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

		hash, err := auth.HashPassword(password)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		u := &datastore.User{
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

	app := handler.NewApp(db, uRepo, bmRepo, store)

	r := mux.NewRouter()
	r.Handle("/", app.AuthMiddleware(app.HomeHandler()))
	r.Handle("/bookmarks", app.AuthMiddleware(app.BookmarksHandler()))
	r.Handle("/bookmarks/new", app.AuthMiddleware(app.BookmarksNewHandler()))
	r.Handle("/bookmarks/{id:[0-9]+}", app.AuthMiddleware(app.BookmarksShowHandler()))
	r.Handle("/import", app.AuthMiddleware(app.ImportHandler())).Methods("GET", "POST")
	r.Handle("/account", app.AuthMiddleware(app.AccountShowHandler()))
	r.Handle("/account/edit", app.AuthMiddleware(app.AccountEditHandler()))

	r.Handle("/login", app.LoginHandler())
	r.Handle("/logout", app.LogoutHandler())

	r.NotFoundHandler = app.NotFoundHandler()

	api := r.PathPrefix("/api").Subrouter()
	api.Handle("/bookmarks", handler.TokenAuthMiddleware(app.CreateBookmarkHandler())).Methods("POST")
	api.Handle("/bookmarks", handler.TokenAuthMiddleware(app.ApiBookmarksHandler())).Methods("GET")
	api.Handle("/ping", handler.CORSMiddleware(app.ApiPingHandler())).Methods("GET")
	api.Handle("/auth", handler.CORSMiddleware(app.CreateTokenHandler())).Methods("POST")
	api.Handle("/auth", app.PreflightHandler()).Methods("OPTIONS")

	// Static file handler
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Build middleware chain that is run for all requests
	chain := alice.New(handler.LoggingMiddleware, handler.TimeoutMiddleware).Then(r)

	if *usefcgi {
		err = server.FCGI(nil, chain)
	} else if *usecgi {
		err = server.CGI(chain)
	} else {
		srv := server.HTTP
		srv.Addr = *addr
		srv.Handler = chain
		err = srv.ListenAndServe()
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
