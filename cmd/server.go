package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/efy/gorecall/database"
	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/handler"
	"github.com/efy/gorecall/router"
	"github.com/efy/gorecall/server"
	"github.com/efy/gorecall/subcmd"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
)

var serve = subcmd.Command{
	UsageLine: "serve",
	Short:     "serve the web app using one of the supported methods",
	Run: func(cmd *subcmd.Command, args []string) {
		addr := cmd.Flag.String("addr", ":8080", "the address to bind to when using the http server")
		dbname := cmd.Flag.String("dbname", "gorecall.db", "path to database file")
		usecgi := cmd.Flag.Bool("cgi", false, "Serve app using cgi")
		usefcgi := cmd.Flag.Bool("fcgi", false, "Serve app using fastcgi")
		cmd.ParseFlags(args)

		db, err := database.Init(*dbname)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		uRepo, err := datastore.NewUserRepo(db)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bmRepo, err := datastore.NewBookmarkRepo(db)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		trRepo, err := datastore.NewTagRepo(db)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		store := sessions.NewCookieStore([]byte("something-very-secret"))

		app := handler.NewApp(db, uRepo, bmRepo, trRepo, store)

		m := router.App()
		m.Get(router.Dashboard).Handler(app.AuthMiddleware(app.HomeHandler()))

		m.Get(router.Bookmarks).Handler(app.AuthMiddleware(app.BookmarksHandler()))
		m.Get(router.NewBookmark).Handler(app.AuthMiddleware(app.BookmarksNewHandler()))
		m.Get(router.CreateBookmark).Handler(app.AuthMiddleware(app.CreateBookmarkHandler()))
		m.Get(router.Bookmark).Handler(app.AuthMiddleware(app.BookmarksShowHandler()))
		m.Get(router.BookmarkAddTag).Handler(app.AuthMiddleware(app.BookmarkAddTagHandler()))
		m.Get(router.BookmarkRemoveTag).Handler(app.AuthMiddleware(app.BookmarkRemoveTagHandler()))

		m.Get(router.Tags).Handler(app.AuthMiddleware(app.TagsHandler()))
		m.Get(router.NewTag).Handler(app.AuthMiddleware(app.NewTagHandler()))
		m.Get(router.Tag).Handler(app.AuthMiddleware(app.TagHandler()))
		m.Get(router.CreateTag).Handler(app.AuthMiddleware(app.CreateTagHandler()))

		m.Get(router.Import).Handler(app.AuthMiddleware(app.ImportHandler()))
		m.Get(router.Export).Handler(app.AuthMiddleware(app.ExportHandler()))
		m.Get(router.Account).Handler(app.AuthMiddleware(app.AccountShowHandler()))
		m.Get(router.Preferences).Handler(app.AuthMiddleware(app.PreferencesHandler()))
		m.Get(router.EditAccount).Handler(app.AuthMiddleware(app.AccountEditHandler()))
		m.Get(router.Login).Handler(app.LoginHandler())
		m.Get(router.Logout).Handler(app.LogoutHandler())

		m.NotFoundHandler = app.NotFoundHandler()

		// Static file handler
		m.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

		api := router.Api()
		api.Get(router.CreateBookmark).Handler(handler.TokenAuthMiddleware(app.ApiCreateBookmarkHandler()))
		api.Get(router.Bookmarks).Handler(handler.TokenAuthMiddleware(app.ApiBookmarksHandler()))
		api.Get(router.Ping).Handler(handler.CORSMiddleware(app.ApiPingHandler()))
		api.Get(router.Authenticate).Handler(handler.CORSMiddleware(app.CreateTokenHandler()))
		api.Get(router.Preflight).Handler(app.PreflightHandler())
		api.Get(router.WebInfo).Handler(app.WebInfoHandler())

		r := mux.NewRouter()

		// Mount api router onto app
		r.PathPrefix("/api/").Handler(http.StripPrefix("/api", api))

		// Mount the web app
		r.PathPrefix("/").Handler(m)

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
	},
}
