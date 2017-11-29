package templates

const bookmarktmpl = `
{{ define "content" }}

{{ if .Bookmark }}
<div class="rc-filter-bar">
	<ul class="breadcrumb">
		<li class="breadcrumb-item">
			<a href="/bookmarks">Bookmarks</a>
		</li>
		<li class="breadcrumb-item">
			<a href="/bookmarks/{{ .Bookmark.ID }}">{{ .Bookmark.Title }}</a>
		</li>
	</ul>
</div>
<dl>
  <dt>ID</dt>
  <dd>{{ .Bookmark.ID }}</dd>
  <dt>Title</dt>
  <dd>{{ .Bookmark.Title }}</dd>
  <dt>URL</dt>
  <dd>{{ .Bookmark.URL }}</dd>
</dl>
{{ end }}

{{ end }}
`
