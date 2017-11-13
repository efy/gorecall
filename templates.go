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
    <header>
      <nav>
        <a href="/bookmarks">Bookmarks</a>
        <a href="/import">Import</a>
        <a href="/login">Login</a>
      </nav>
    </header>
    {{ block "content" . }} {{ end }}
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
    <th>Icon</th>
    <th>Title</th>
    <th>Created</th>
  </thead>
</table>

{{ end }}
`

const importtmpl = `
{{ define "content" }}

<h2>Import</h2>

<form method="post" action="/import">
  <label for="bookmarks">bookmarks file</label>
  <input type="file" name="bookmarks">

  <button>import</button>
</form>

{{ end }}
`

const logintmpl = `
{{ define "content" }}

<h2>Login</h2>

<form method="post" action="/login">
  <label for="username">Username</label>
  <input type="text" name="username">

  <label for="password">Password</label>
  <input type="password" name="password">

  <button>Login</button>
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
}

func RenderTemplate(w io.Writer, t string, data interface{}) {
	err := templates[t].Execute(w, data)
	if err != nil {
		fmt.Println("error rendering", t, err)
	}
}
