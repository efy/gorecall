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
	registerTemplate("newbookmark.html", newbookmarktmpl)
	registerTemplate("bookmark.html", bookmarktmpl)

	registerTemplate("tags.html", tagstmpl)
	registerTemplate("newtag.html", newtagtmpl)
	registerTemplate("tag.html", tagtmpl)

	registerTemplate("import.html", importtmpl)
	registerTemplate("export.html", exporttmpl)
	registerTemplate("importsuccess.html", importsuccesstmpl)
	registerTemplate("account.html", accounttmpl)
	registerTemplate("preferences.html", preferencestmpl)

	registerTemplate("login.html", logintmpl)

	registerTemplate("servererror.html", servererrortmpl)
	registerTemplate("notfound.html", notfoundtmpl)
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
