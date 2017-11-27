package templates

import (
	"fmt"
	"html/template"
	"io"
	"net/url"
	"time"

	humanize "github.com/dustin/go-humanize"
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
			<div class="columns vh">
				<div class="column col-2 vh">
					<div class="rc-sidebar">
						<img src="/public/logo.svg" class="logo">
						<ul class="nav">
							{{ if .Authenticated }}
								<li class="nav-item">
									<a href="/bookmarks" class="btn btn-link">Bookmarks</a>
								</li>
								<li class="nav-item">
									<a href="/import" class="btn btn-link">Import</a>
								</li>
								<li class="divider"></li>
								<li class="nav-item">
									<a href="/logout" class="btn btn-link">Logout</a>
								</li>
								<li class="nav-item">
									<a href="/account" class="btn btn-link">
										{{ .Username }}
										<figure class="avatar avatar-sm" data-initial="X" style="background-color: #5755d;">
											<img>
										</figure>
									</a>
								</li>
							{{ end }}
						</ul>
					</div>
				</div>
				<div class="column col-10">
					<div class="content-area container">
						{{ block "content" . }} {{ end }}
					</div>
				</div>
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

<h2>Dashboard</h2>

{{ end }}
`

const bookmarkstmpl = `
{{ define "content" }}

<div class="rc-header">
	<h2>Bookmarks</h2>
	<div class="rc-header-actions">
		<a href="/bookmarks/new" class="btn">New</a>
	</div>
</div>

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
<div class="rc-bm-list">
  {{ range .Bookmarks }}
    <div class="rc-bookmark columns">
      <div class="text-center rc-bm-favicon column col-1">
				{{ if .Icon }}
					<img width="20" height="20" src="{{ .Icon | base64 }}">
				{{ else }}
					<img width="20" height="20" src="" onerror="this.src = '/public/placeholder_favicon.png'">
				{{ end }}
      </div>
      <div class="column col-11">
        <div class="rc-bm-title text-ellipsis">
          <a href="{{ .URL }}" target="_blank" rel="noopener">
            {{ .Title | html }}
          </a>
        </div>
				<div class="rc-bm-details">
					<time>
						{{ .Created | timeago }}
					</time>
					•
					<a href="/bookmarks/{{ .ID }}">
						show
					</a>
					•
					<a href="{{ .URL | website }}" rel="noopener" target="_blank">
						{{ .URL | domain }}
					</a>
				</div>
      </div>
    </div>
  {{ end }}
</div>
<div class="rc-pagination-container">
	<ul class="pagination">
		{{ if lt .Pagination.Prev 0 }}
			<li class="page-item disabled"><a href="#">Previous</a></li>
		{{ else }}
			<li class="page-item"><a href="/bookmarks?per_page={{ .Pagination.PerPage }}&page={{ .Pagination.Prev }}">Previous</a></li>
		{{ end }}

		<li class="page-item active">
			<a href="/bookmarks?per_page={{ .Pagination.PerPage }}&page={{ .Pagination.Current }}">{{ .Pagination.Current }}</a>
		</li>

		{{ $save := .Pagination }}
		{{ range $page := .Pagination.List }}
			{{ if lt $page $save.Last }}
				<li class="page-item">
					<a href="/bookmarks?per_page={{ $save.PerPage }}&page={{ $page}}">{{ $page }}</a>
				</li>
			{{ end }}
		{{ end }}
	
		{{ if ne .Pagination.Last .Pagination.Current }}
			<li class="page-item">
				<span>...</span>
			</li>
			<li class="page-item">
				<a href="/bookmarks?per_page={{ .Pagination.PerPage }}&page={{ .Pagination.Last }}">{{ .Pagination.Last }}</a>
			</li>
		{{ end }}

		{{ if gt .Pagination.Next .Pagination.Last }}
			<li class="page-item disabled"><a href="#">Next</a></li>
		{{ else }}
			<li class="page-item"><a href="/bookmarks?per_page={{ .Pagination.PerPage }}&page={{ .Pagination.Next }}">Next</a></li>
		{{ end }}
	</ul>
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

<h2>Show bookmark</h2>

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
<h2>New Bookmark</h2>

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

<h2>Import</h2>

<div>
  <p>
    Successfully imported {{ len .Bookmarks }} bookmarks.
  </p>
</div>
{{ end }}
`

const importtmpl = `
{{ define "content" }}

<h2>Import</h2>

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

var funcMap = template.FuncMap{
	"base64": func(s string) template.URL {
		return template.URL(s)
	},
	"website": func(s string) template.URL {
		u, err := url.Parse(s)
		if err != nil {
			return ""
		}
		return template.URL(u.Scheme + "://" + u.Host)
	},
	"domain": func(s string) string {
		u, err := url.Parse(s)
		if err != nil {
			return ""
		}
		return u.Host
	},
	"timeago": func(t time.Time) string {
		return humanize.Time(t)
	},
	"html": func(s string) template.HTML {
		return template.HTML(s)
	},
}

func init() {
	registerTemplate("index.html", indextmpl)
	registerTemplate("bookmarks.html", bookmarkstmpl)
	registerTemplate("bookmarksnew.html", bookmarksnewtmpl)
	registerTemplate("bookmarksshow.html", bookmarksshowtmpl)
	registerTemplate("import.html", importtmpl)
	registerTemplate("importsuccess.html", importsuccesstmpl)
	registerTemplate("accountshow.html", accountshowtmpl)
	registerTemplate("accountedit.html", accountedittmpl)
	registerTemplate("servererror.html", servererrortmpl)
	registerTemplate("notfound.html", notfoundtmpl)
	registerTemplate("login.html", logintmpl)
}

// Helper to compile template within a layout context with funcs
func registerTemplate(label string, tmpl string) {
	templates[label] = template.Must(template.New("layout").Parse(layouttmpl))
	template.Must(templates[label].Funcs(funcMap).Parse(tmpl))
}

func RenderTemplate(w io.Writer, t string, data interface{}) {
	err := templates[t].Execute(w, data)
	if err != nil {
		fmt.Println("error rendering", t, err)
	}
}
