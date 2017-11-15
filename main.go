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

	appHandler := AppHandler{
		db: db,
	}

	r := mux.NewRouter()
	r.Handle("/", AuthMiddleware(HomeHandler(appHandler)))
	r.Handle("/bookmarks", AuthMiddleware(BookmarksHandler(appHandler)))
	r.Handle("/bookmarks/new", AuthMiddleware(BookmarksNewHandler(appHandler)))
	r.Handle("/bookmarks/{id:[0-9]+}", AuthMiddleware(BookmarksShowHandler(appHandler)))
	r.Handle("/import", AuthMiddleware(ImportHandler(appHandler)))
	r.Handle("/account", AuthMiddleware(AccountShowHandler(appHandler)))
	r.Handle("/account/edit", AuthMiddleware(AccountEditHandler(appHandler)))

	r.Handle("/login", LoginHandler(appHandler))
	r.Handle("/logout", LogoutHandler(appHandler))

	r.NotFoundHandler = NotFoundHandler(appHandler)

	api := r.PathPrefix("/api").Subrouter()
	api.Handle("/bookmarks", TokenAuthMiddleware(CreateBookmarkHandler(appHandler))).Methods("POST")
	api.Handle("/bookmarks", TokenAuthMiddleware(ApiBookmarksHandler(appHandler))).Methods("GET")
	api.Handle("/ping", CORSMiddleware(ApiPingHandler(appHandler))).Methods("GET")
	api.Handle("/auth", CORSMiddleware(CreateTokenHandler(appHandler))).Methods("POST")
	api.Handle("/auth", PreflightHandler(appHandler)).Methods("OPTIONS")

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
