package handler

import (
	"net/http"

	"github.com/efy/gorecall/datastore"
	"github.com/efy/gorecall/router"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
)

type App struct {
	ur    datastore.UserRepo
	br    datastore.BookmarkRepo
	tr    datastore.TagRepo
	store *sessions.CookieStore
}

func (app *App) Handler() http.Handler {
	r := router.App()
	r.Get(router.Dashboard).Handler(app.AuthMiddleware(app.HomeHandler()))

	r.Get(router.Bookmarks).Handler(app.AuthMiddleware(app.BookmarksHandler()))
	r.Get(router.NewBookmark).Handler(app.AuthMiddleware(app.NewBookmarkHandler()))
	r.Get(router.CreateBookmark).Handler(app.AuthMiddleware(app.CreateBookmarkHandler()))
	r.Get(router.DeleteBookmark).Handler(app.AuthMiddleware(app.DeleteBookmarkHandler()))
	r.Get(router.Bookmark).Handler(app.AuthMiddleware(app.BookmarkHandler()))
	r.Get(router.BookmarkAddTag).Handler(app.AuthMiddleware(app.BookmarkAddTagHandler()))
	r.Get(router.BookmarkRemoveTag).Handler(app.AuthMiddleware(app.BookmarkRemoveTagHandler()))

	r.Get(router.Tags).Handler(app.AuthMiddleware(app.TagsHandler()))
	r.Get(router.NewTag).Handler(app.AuthMiddleware(app.NewTagHandler()))
	r.Get(router.Tag).Handler(app.AuthMiddleware(app.TagHandler()))
	r.Get(router.CreateTag).Handler(app.AuthMiddleware(app.CreateTagHandler()))

	r.Get(router.Import).Handler(app.AuthMiddleware(app.ImportHandler()))
	r.Get(router.Export).Handler(app.AuthMiddleware(app.ExportHandler()))
	r.Get(router.Account).Handler(app.AuthMiddleware(app.AccountShowHandler()))
	r.Get(router.Preferences).Handler(app.AuthMiddleware(app.PreferencesHandler()))
	r.Get(router.EditAccount).Handler(app.AuthMiddleware(app.AccountEditHandler()))
	r.Get(router.Login).Handler(app.LoginHandler())
	r.Get(router.Logout).Handler(app.LogoutHandler())

	r.NotFoundHandler = app.NotFoundHandler()

	// Static file handler
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Build middleware chain for app requests
	appchain := alice.New(LoggingMiddleware, TimeoutMiddleware).Then(r)

	return appchain
}

func NewApp(ur datastore.UserRepo, br datastore.BookmarkRepo, tr datastore.TagRepo, store *sessions.CookieStore) App {
	return App{
		ur,
		br,
		tr,
		store,
	}
}

type Api struct {
	ur datastore.UserRepo
	br datastore.BookmarkRepo
	tr datastore.TagRepo
}

func (api *Api) Handler() http.Handler {
	r := router.Api()
	r.Get(router.CreateBookmark).Handler(TokenAuthMiddleware(api.ApiCreateBookmarkHandler()))
	r.Get(router.Bookmarks).Handler(TokenAuthMiddleware(api.ApiBookmarksHandler()))
	r.Get(router.Ping).Handler(api.ApiPingHandler())
	r.Get(router.Authenticate).Handler(api.CreateTokenHandler())
	r.Get(router.WebInfo).Handler(api.WebInfoHandler())

	// Build middleware chaun for api requests
	apichain := alice.New(LoggingMiddleware, TimeoutMiddleware, CORSMiddleware, PreflightMiddleware).Then(r)

	return apichain
}

func NewApi(ur datastore.UserRepo, br datastore.BookmarkRepo, tr datastore.TagRepo) Api {
	return Api{
		ur,
		br,
		tr,
	}
}
