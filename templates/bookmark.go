package templates

const bookmarktmpl = `
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
