package router

import "github.com/gorilla/mux"

// Route Names
const (
	// Shared
	Tags      = "tags"
	CreateTag = "tags:create"
	Tag       = "tag"
	DeleteTag = "tag:delete"

	Bookmarks         = "bookmarks:index"
	CreateBookmark    = "bookmarks:create"
	Bookmark          = "bookmark"
	DeleteBookmark    = "bookmark:delete"
	BookmarkAddTag    = "bookmark:addtag"
	BookmarkRemoveTag = "bookmark:removetag"
	BookmarkWebinfo   = "bookmark:webinfo"
	SearchBookmarks   = "bookmarks:search"

	Account       = "account"
	UpdateAccount = "account:update"

	Preferences       = "preferences"
	UpdatePreferences = "preferences:update"

	Import = "import"
	Export = "export"

	// Webapp only
	Dashboard   = "dashboard:index"
	Login       = "login"
	Logout      = "logout"
	NewBookmark = "bookmarks:new"
	NewTag      = "tag:new"

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
	m.Path("/bookmarks/{id:[0-9]+}/webinfo").Methods("PUT", "POST").Name(BookmarkWebinfo)
	m.Path("/tags").Methods("GET").Name(Tags)
	m.Path("/tags").Methods("POST").Name(CreateTag)
	m.Path("/tags/new").Methods("GET").Name(NewTag)
	m.Path("/tags/{id:[0-9]+}").Methods("GET").Name(Tag)
	m.Path("/tags/{id:[0-9]+}").Methods("DELETE").Name(DeleteTag)
	m.Path("/settings/preferences").Methods("GET").Name(Preferences)
	m.Path("/settings/account").Methods("GET").Name(Account)
	m.Path("/settings/account").Methods("POST").Name(UpdateAccount)
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
	m.Path("/bookmarks/{id:[0-9]+}").Methods("DELETE").Name(DeleteBookmark)
	m.Path("/bookmarks/{id:[0-9]+}/addtag").Methods("POST").Name(BookmarkAddTag)
	m.Path("/bookmarks/{id:[0-9]+}/removetag").Methods("DELETE").Name(BookmarkRemoveTag)
	m.Path("/bookmarks/{id:[0-9]+}/webinfo").Methods("PUT").Name(BookmarkWebinfo)
	m.Path("/tags").Methods("GET").Name(Tags)
	m.Path("/tags").Methods("POST").Name(CreateTag)
	m.Path("/tags/{id:[0-9]+}").Methods("GET").Name(Tag)
	m.Path("/tags/{id:[0-9]+}").Methods("DELETE").Name(DeleteTag)
	m.Path("/ping").Methods("GET").Name(Ping)
	m.Path("/auth").Methods("POST").Name(Authenticate)
	m.Path("/webinfo").Methods("GET").Name(WebInfo)

	// Should be an api middleware
	m.Path("/auth").Methods("OPTIONS").Name(Preflight)

	return m
}
