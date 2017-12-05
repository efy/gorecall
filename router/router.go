package router

import "github.com/gorilla/mux"

// Route Names
const (
	// Shared
	Bookmarks         = "bookmarks:index"
	CreateBookmark    = "bookmarks:create"
	Bookmark          = "bookmark"
	DeleteBookmark    = "bookmark:delete"
	Tags              = "tags"
	CreateTag         = "tags:create"
	Tag               = "tag"
	BookmarkAddTag    = "bookmark:addtag"
	BookmarkRemoveTag = "bookmark:removetag"
	SearchBookmarks   = "bookmarks:search"

	// Webapp only
	Dashboard   = "dashboard:index"
	Import      = "import"
	Export      = "export"
	Account     = "account"
	EditAccount = "account:edit"
	Login       = "login"
	Logout      = "logout"
	NewBookmark = "bookmarks:new"
	NewTag      = "tag:new"
	Preferences = "preferences"

	// Api only
	Ping         = "ping"
	Authenticate = "auth"
	Preflight    = "preflight"
	WebInfo      = "webinfo"
)

func App() *mux.Router {
	m := mux.NewRouter()

	m.Path("/").Methods("GET").Name(Dashboard)
	m.Path("/bookmarks").Methods("GET").Name(Bookmarks)
	m.Path("/bookmarks").Methods("POST").Name(CreateBookmark)
	m.Path("/bookmarks/search").Methods("GET").Name(SearchBookmarks)
	m.Path("/bookmarks/new").Methods("GET").Name(NewBookmark)
	m.Path("/bookmarks/{id:[0-9]+}").Methods("GET").Name(Bookmark)
	m.Path("/bookmarks/{id:[0-9]+}").Methods("DELETE").Name(DeleteBookmark)
	m.Path("/bookmarks/{id:[0-9]+}/addtag").Methods("POST").Name(BookmarkAddTag)
	m.Path("/bookmarks/{id:[0-9]+}/removetag").Methods("POST").Name(BookmarkRemoveTag)
	m.Path("/tags").Methods("GET").Name(Tags)
	m.Path("/tags").Methods("POST").Name(CreateTag)
	m.Path("/tags/new").Methods("GET").Name(NewTag)
	m.Path("/tags/{id:[0-9]+}").Methods("GET").Name(Tag)
	m.Path("/settings/preferences").Methods("GET").Name(Preferences)
	m.Path("/settings/account").Methods("GET").Name(Account)
	m.Path("/settings/account/edit").Methods("GET").Name(EditAccount)
	m.Path("/settings/import").Methods("GET", "POST").Name(Import)
	m.Path("/settings/export").Methods("GET").Name(Export)
	m.Path("/login").Methods("GET", "POST").Name(Login)
	m.Path("/logout").Methods("GET", "POST").Name(Logout)

	return m
}

func Api() *mux.Router {
	m := mux.NewRouter()
	m.Path("/bookmarks").Methods("GET").Name(Bookmarks)
	m.Path("/bookmarks/search").Methods("GET").Name(SearchBookmarks)
	m.Path("/bookmarks").Methods("POST").Name(CreateBookmark)
	m.Path("/bookmarks/{id:[0-9]+}").Methods("GET").Name(Bookmark)
	m.Path("/tags").Methods("GET").Name(Tags)
	m.Path("/tags").Methods("POST").Name(CreateTag)
	m.Path("/tags/new").Methods("GET").Name(NewTag)
	m.Path("/tags/{id:[0-9]+}").Methods("GET").Name(Tag)
	m.Path("/ping").Methods("GET").Name(Ping)
	m.Path("/auth").Methods("POST").Name(Authenticate)
	m.Path("/webinfo").Methods("GET").Name(WebInfo)

	// Should be an api middleware
	m.Path("/auth").Methods("OPTIONS").Name(Preflight)

	return m
}
