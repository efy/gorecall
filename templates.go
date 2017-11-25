package main

import (
	"fmt"
	"html/template"
	"io"
)

var templates map[string]*template.Template = make(map[string]*template.Template)

const layouttmpl = `
<!DOCTYPE html>
<html>
  <head>
    <title>Home</title>
    <link rel="stylesheet" href="/public/style.css">
  </head>
  <body>
    <div class="wrapper">
      <div class="container">
        <header class="navbar">
          <nav class="navbar-section">
            <img src="/public/logo.svg" class="logo">
            {{ if .Authenticated }}
              <a href="/bookmarks" class="btn btn-link">Bookmarks</a>
              <a href="/import" class="btn btn-link">Import</a>
            {{ end }}
          </nav>
          <nav class="navbar-section">
            {{ if .Authenticated }}
              <a href="/logout" class="btn btn-link">Logout</a>
              <a href="/account" class="btn btn-link">
                {{ .Username }}
                <figure class="avatar avatar-sm" data-initial="X" style="background-color: #5755d;">
                  <img>
                </figure>
              </a>
            {{ else }}
              <a href="/login" class="btn btn-link">Login</a>
            {{ end }}
          </nav>
        </header>
      </div>

      <div class="content-area container">
        {{ block "content" . }} {{ end }}
      </div>
    </div>
  </body>
</html>
`
const servererrortmpl = `
{{ define "content" }}

<h2>500 Internal server error</h2>
<pre>{{ .Error }}</pre>

{{ end }}
`

const notfoundtmpl = `
{{ define "content" }}

<h2 class="text-center">404 Not found</h2>

{{ end }}
`

const indextmpl = `
{{ define "content" }}

<h2 class="text-center">Home</h2>

{{ end }}
`

const bookmarkstmpl = `
{{ define "content" }}

<h2 class="text-center">Bookmarks</h2>

{{ if not .Bookmarks }}
<div class="empty">
  <div class="empty-icon">
    <i class="icon icon-bookmark"></i>
  </div>
  <p class="empty-title h5">You have no Bookmarks</p>
  <p class="empty-subtitle">Choose from the actions below to get started</p>
  <div class="empty-action">
    <a href="/import" class="btn btn-primary">Import</a>
    <a href="/bookmarks/new" class="btn btn-primary">Add</a>
  </div>
</div>
{{ else }}
<div class="rc-list">
  {{ range .Bookmarks }}
    <div class="rc-bookmark columns">
      <div class="column col-1">
        <figure class="avatar" data-initial="BM" style="background: #5755d9;">
        </figure>
      </div>
      <div class="column col-11">
        <div class="rc-bm-title">
          <a href="{{ .URL }}" target="_blank" rel="noopener">
            {{ .Title }}
          </a>
        </div>
      </div>
    </div>
  {{ end }}
</div>
{{ end }}

{{ end }}
`

const accountshowtmpl = `
{{ define "content" }}

<header class="navbar">
  <div class="navbar-section">
    <h2>Account</h2>
  </div>
  <div class="navbar-section">
    <a class="btn btn-default" href="/account/edit">Update Account</a>
  <div>
</header>

<dl>
  <dt>ID</dt>
  <dd>{{ .User.ID}}</dd>
  <dt>Username</dt>
  <dd>{{ .User.Username}}</dd>
</dl>

{{ end }}
`

const accountedittmpl = `
{{ define "content" }}

<header class="navbar">
  <div class="navbar-section">
    <h2>Account</h2>
  </div>
  <div class="navbar-section">
    <a class="btn btn-default" href="/account">Show Account</a>
  <div>
</header>

<form action="/account/edit" method="post">
  <div class="form-group">
    <label class="form-label" for="username">Username</label>
    <input id="username" class="form-input" value="{{ .User.Username }}" type="text" name="username">
  </div>

  <div class="form-group">
    <button type="submit" class="btn btn-primary">Update</button>
  </div>
</form>

{{ end }}
`

const bookmarksshowtmpl = `
{{ define "content" }}

<h2 class="text-center">Show bookmark</h2>

<dl>
  <dt>ID</dt>
  <dd>{{ .Bookmark.ID }}</dd>
  <dt>Title</dt>
  <dd>{{ .Bookmark.Title }}</dd>
  <dt>URL</dt>
  <dd>{{ .Bookmark.URL }}</dd>
</dl>

{{ end }}
`

const bookmarksnewtmpl = `
{{ define "content" }}
<h2 class="text-center">New Bookmark</h2>

<form action="/bookmarks/new" method="post">
  <div class="form-group">
    <label class="form-label" for="title">Title</label>
    <input id="title" class="form-input" type="text" name="title">
  </div>

  <div class="form-group">
    <label class="form-label" for="url">URL</label>
    <input id="url" class="form-input" type="text" name="url">
  </div>

  <div class="form-group">
    <button type="submit" class="btn btn-primary">Create</button>
    <button type="reset" class="btn">Reset</button>
  </div>
</form>

{{ end }}
`

const importsuccesstmpl = `
{{ define "content" }}

<h2 class="text-center">Import</h2>

<div>
  <p class="text-center">
    Successfully imported {{ len .Bookmarks }} bookmarks.
  </p>
</div>
{{ end }}
`

const importtmpl = `
{{ define "content" }}

<h2 class="text-center">Import</h2>

<form enctype="multipart/form-data" method="post" action="/import">
  <div class="form-group">
    <label class="form-label" for="bookmarks">File</label>
    <input id="bookmarks_file" class="form-input" type="file" name="bookmarks">
    <p class="form-input-hint">
      A file exported from your browser or bookmarking service.
    </p>
  </div>

  <div class="form-group">
    <button type="submit" class="btn btn-primary">Import</button>
  </div>
</form>

{{ end }}
`

const logintmpl = `
{{ define "content" }}

<h2 class="text-center">Login</h2>

<form method="post" action="/login">
  <div class="form-group">
    <label class="form-label" for="username">Username</label>
    <input id="username" class="form-input" type="text" name="username">
  </div>

  <div class="form-group">
    <label class="form-label" for="password">Password</label>
    <input id="password" class="form-input" type="password" name="password">
  </div>

  <div class="form-group">
    <label for="remember_me" class="form-checkbox">
      <input id="remember_me" type="checkbox" name="remember_me">
      <i class="form-icon"></i>Remember me?
    </label>
  </div>

  <div class="form-group">
    <button type="submit" class="btn btn-primary">Login</button>
  </div>
</form>

{{ end }}
`

func init() {
	templates["index.html"] = template.Must(template.New("layout").Parse(layouttmpl))
	template.Must(templates["index.html"].Parse(indextmpl))

	templates["bookmarks.html"] = template.Must(template.New("layout").Parse(layouttmpl))
	template.Must(templates["bookmarks.html"].Parse(bookmarkstmpl))

	templates["import.html"] = template.Must(template.New("layout").Parse(layouttmpl))
	template.Must(templates["import.html"].Parse(importtmpl))

	templates["importsuccess.html"] = template.Must(template.New("layout").Parse(layouttmpl))
	template.Must(templates["importsuccess.html"].Parse(importsuccesstmpl))

	templates["login.html"] = template.Must(template.New("layout").Parse(layouttmpl))
	template.Must(templates["login.html"].Parse(logintmpl))

	templates["bookmarksnew.html"] = template.Must(template.New("layout").Parse(layouttmpl))
	template.Must(templates["bookmarksnew.html"].Parse(bookmarksnewtmpl))

	templates["bookmarksshow.html"] = template.Must(template.New("layout").Parse(layouttmpl))
	template.Must(templates["bookmarksshow.html"].Parse(bookmarksshowtmpl))

	templates["accountshow.html"] = template.Must(template.New("layout").Parse(layouttmpl))
	template.Must(templates["accountshow.html"].Parse(accountshowtmpl))

	templates["accountedit.html"] = template.Must(template.New("layout").Parse(layouttmpl))
	template.Must(templates["accountedit.html"].Parse(accountedittmpl))

	templates["servererror.html"] = template.Must(template.New("layout").Parse(layouttmpl))
	template.Must(templates["servererror.html"].Parse(servererrortmpl))

	templates["notfound.html"] = template.Must(template.New("layout").Parse(layouttmpl))
	template.Must(templates["notfound.html"].Parse(notfoundtmpl))
}

func RenderTemplate(w io.Writer, t string, data interface{}) {
	err := templates[t].Execute(w, data)
	if err != nil {
		fmt.Println("error rendering", t, err)
	}
}
