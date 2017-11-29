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
<table class="table">
	<tbody>
		<tr>
			<th>ID</th>
			<td> {{ .Bookmark.ID }}</td>
		</tr>
		<tr>
			<th>Title</th>
			<td> {{ .Bookmark.Title }}</td>
		</tr>
		<tr>
			<th>URL</th>
			<td> {{ .Bookmark.URL }}</td>
		</tr>
		<tr>
			<th>Created</th>
			<td> {{ .Bookmark.Created}}</td>
		</tr>
	</tbody>
</table>
{{ end }}

{{ end }}
`
