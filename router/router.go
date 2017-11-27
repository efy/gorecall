package router

import "github.com/gorilla/mux"

// Route Names
const (
	// Shared
	Bookmarks      = "bookmarks:index"
	CreateBookmark = "bookmarks:create"
	Bookmark       = "bookmark"

	// Webapp only
	Dashboard   = "dashboard:index"
	Import      = "import"
	Account     = "account"
	EditAccount = "account:edit"
	Login       = "login"
	Logout      = "logout"
	NewBookmark = "bookmarks:new"

	// Api only
	Ping         = "ping"
	Authenticate = "auth"
	Preflight    = "preflight"
)

func App() *mux.Router {
	m := mux.NewRouter()

	m.Path("/").Methods("GET").Name(Dashboard)
	m.Path("/bookmarks").Methods("GET").Name(Bookmarks)
	m.Path("/bookmarks").Methods("POST").Name(CreateBookmark)
	m.Path("/bookmarks/new").Methods("GET").Name(NewBookmark)
	m.Path("/bookmarks/{id:[0-9]+}").Methods("GET").Name(Bookmark)
	m.Path("/import").Methods("GET", "POST").Name(Import)
	m.Path("/account").Methods("GET").Name(Account)
	m.Path("/account/edit").Methods("GET").Name(EditAccount)
	m.Path("/login").Methods("GET", "POST").Name(Login)
	m.Path("/logout").Methods("GET", "POST").Name(Logout)

	return m
}

func Api() *mux.Router {
	m := mux.NewRouter()
	m.Path("/bookmarks").Methods("GET").Name(Bookmarks)
	m.Path("/bookmarks").Methods("POST").Name(CreateBookmark)
	m.Path("/bookmarks/{id:[0-9]+}").Methods("GET").Name(Bookmark)
	m.Path("/ping").Methods("GET").Name(Ping)
	m.Path("/auth").Methods("POST").Name(Authenticate)

	// Should be an api middleware
	m.Path("/auth").Methods("OPTIONS").Name(Preflight)

	return m
}