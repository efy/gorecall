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
      <header>
        <nav class="top-nav">
          <div class="nav-left">
            <a href="/bookmarks">Bookmarks</a>
            <a href="/import">Import</a>
            <a href="/login">Login</a>
          </div>
          <div class="nav-right">
            <a href="/bookmarks/new" class="button">New</a>
          </div>
        </nav>
      </header>
      {{ block "content" . }} {{ end }}
    </div>
  </body>
</html>
`

const indextmpl = `
{{ define "content" }}

<h2>Home</h2>

{{ end }}
`

const bookmarkstmpl = `
{{ define "content" }}

<h2>Bookmarks</h2>

<table>
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
`

const bookmarksshowtmpl = `
{{ define "content" }}

<h2>Show bookmark</h2>

<dl>
  <dt>ID</dt>
  <dd>{{ .ID }}</dd>
  <dt>Title</dt>
  <dd>{{ .Title }}</dd>
  <dt>URL</dt>
  <dd>{{ .URL }}</dd>
</dl>

{{ end }}
`

const bookmarksnewtmpl = `
{{ define "content" }}
<h2>New Bookmark</h2>

<form action="/bookmarks/new" method="post">
  <div class="field">
    <label for="title">Title</label>
    <input type="text" name="title">
  </div>

  <div class="field">
    <label for="url">URL</label>
    <input type="text" name="url">
  </div>

  <div class="field">
    <button type="submit" class="button-primary">Create</button>
    <button type="reset" disabled>Reset</button>
  </div>
</form>

{{ end }}
`

const importtmpl = `
{{ define "content" }}

<h2>Import</h2>

<form method="post" action="/import">
  <div class="field">
    <label for="bookmarks">bookmarks file</label>
    <input type="file" name="bookmarks">
  </div>

  <div class="field">
    <button type="submit" class="button-primary">Import</button>
  </div>
</form>

{{ end }}
`

const logintmpl = `
{{ define "content" }}

<h2>Login</h2>

<form method="post" action="/login">
  <div class="field">
    <label for="username">Username</label>
    <input type="text" name="username">
  </div>

  <div class="field">
    <label for="password">Password</label>
    <input type="password" name="password">
  </div>

  <div class="field">
    <button type="submit" class="button-primary">Login</button>
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
}

func RenderTemplate(w io.Writer, t string, data interface{}) {
	err := templates[t].Execute(w, data)
	if err != nil {
		fmt.Println("error rendering", t, err)
	}
}
