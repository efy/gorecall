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

<table class="table">
  <thead>
    <th>Title</th>
    <th>URL</th>
  </thead>
  <tbody>
    {{ range . }}
      <tr>
        <td>
          {{ .Title }}
        </td>
        <td>
          {{ .URL }}
        </td>
      </tr>
    {{ end }}
  </tbody>
</table>

{{ end }}

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

const importtmpl = `
{{ define "content" }}

<h2 class="text-center">Import</h2>

<form method="post" action="/import">
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

	templates["login.html"] = template.Must(template.New("layout").Parse(layouttmpl))
	template.Must(templates["login.html"].Parse(logintmpl))

	templates["bookmarksnew.html"] = template.Must(template.New("layout").Parse(layouttmpl))
	template.Must(templates["bookmarksnew.html"].Parse(bookmarksnewtmpl))

	templates["bookmarksshow.html"] = template.Must(template.New("layout").Parse(layouttmpl))
	template.Must(templates["bookmarksshow.html"].Parse(bookmarksshowtmpl))

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
