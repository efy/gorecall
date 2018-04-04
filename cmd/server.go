package cmd

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/blevesearch/bleve"
	"github.com/efy/gorecall/database"
	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/handler"
	"github.com/efy/gorecall/server"
	"github.com/efy/gorecall/subcmd"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var serve = subcmd.Command{
	UsageLine: "serve",
	Short:     "serve the web app using one of the supported methods",
	Run: func(cmd *subcmd.Command, args []string) {
		addr := cmd.Flag.String("addr", ":8080", "the address to bind to when using the http server")
		dbdsn := cmd.Flag.String("dsn", "postgres://recall:recall@localhost/recall?sslmode=disable", "data source name")
		indexname := cmd.Flag.String("indexname", "gorecall.idx", "path to index directory")
		usecgi := cmd.Flag.Bool("cgi", false, "Serve app using cgi")
		usefcgi := cmd.Flag.Bool("fcgi", false, "Serve app using fastcgi")
		secret := cmd.Flag.String("secret", "", "Session secret, a key can be generated with the `secret` command. The key should be base64 encoded.")

		cmd.ParseFlags(args)

		db, err := database.Connect(database.Options{
			Driver: "postgres",
			DSN:    *dbdsn,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		index, err := bleve.Open(*indexname)
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

		if *secret == "" {
			fmt.Println("A session secret must be provided")
			os.Exit(1)
		}

		dec, err := base64.StdEncoding.DecodeString(*secret)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		store := sessions.NewCookieStore(dec)

		app := handler.NewApp(uRepo, bmRepo, trRepo, store, index)
		api := handler.NewApi(uRepo, bmRepo, trRepo, index)
		r := mux.NewRouter()

		r.PathPrefix("/api/").Handler(http.StripPrefix("/api", api.Handler()))
		r.PathPrefix("/").Handler(app.Handler())

		if *usefcgi {
			err = server.FCGI(nil, r)
		} else if *usecgi {
			err = server.CGI(r)
		} else {
			srv := server.HTTP
			srv.Addr = *addr
			srv.Handler = r
			err = srv.ListenAndServe()
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
