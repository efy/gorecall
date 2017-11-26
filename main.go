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
	"github.com/efy/gorecall/router"
	"github.com/efy/gorecall/server"
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

	m := router.App()
	m.Get(router.Dashboard).Handler(app.AuthMiddleware(app.HomeHandler()))
	m.Get(router.Bookmarks).Handler(app.AuthMiddleware(app.BookmarksHandler()))
	m.Get(router.NewBookmark).Handler(app.AuthMiddleware(app.BookmarksNewHandler()))
	m.Get(router.Bookmark).Handler(app.AuthMiddleware(app.BookmarksShowHandler()))
	m.Get(router.Import).Handler(app.AuthMiddleware(app.ImportHandler()))
	m.Get(router.Account).Handler(app.AuthMiddleware(app.AccountShowHandler()))
	m.Get(router.EditAccount).Handler(app.AuthMiddleware(app.AccountEditHandler()))
	m.Get(router.Login).Handler(app.LoginHandler())
	m.Get(router.Logout).Handler(app.LogoutHandler())
	m.NotFoundHandler = app.NotFoundHandler()

	api := router.Api()
	api.Get(router.CreateBookmark).Handler(handler.TokenAuthMiddleware(app.CreateBookmarkHandler()))
	api.Get(router.Bookmarks).Handler(handler.TokenAuthMiddleware(app.ApiBookmarksHandler()))
	api.Get(router.Ping).Handler(handler.CORSMiddleware(app.ApiPingHandler()))
	api.Get(router.Authenticate).Handler(handler.CORSMiddleware(app.CreateTokenHandler()))
	api.Get(router.Preflight).Handler(app.PreflightHandler())

	// Static file handler
	m.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Mount api router onto app
	m.PathPrefix("/api/").Handler(http.StripPrefix("/api/", api))

	// Build middleware chain that is run for all requests
	chain := alice.New(handler.LoggingMiddleware, handler.TimeoutMiddleware).Then(m)

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
