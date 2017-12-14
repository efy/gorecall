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
			<th>Media Type</th>
			<td> {{ .Bookmark.MediaType }}</td>
		</tr>
		<tr>
			<th>Description</th>
			<td> {{ .Bookmark.Description }}</td>
		</tr>
		<tr>
			<th>Keywords</th>
			<td> {{ .Bookmark.Keywords }}</td>
		</tr>
		<tr>
			<th>Status</th>
			<td> {{ .Bookmark.Status }}</td>
		</tr>
		<tr>
			<th>Created</th>
			<td>{{ .Bookmark.Created}}</td>
		</tr>
		<tr>
			<th>Tags</th>
			<td>
				{{ range .Tags }}
					<span class="chip">
						<a href="/tags/{{.ID}}">
							{{ .Label }}
						</a>
						<form action="/bookmarks/{{ $.Bookmark.ID }}/removetag" method="post">
							<input type="hidden" value="{{ .ID }}" name="tag_id">
							<button type="submit" class="btn btn-clear"></button>
						</form>
					</span>
				{{ end }}
			</td>
		</tr>
		<tr>
			<th>Add tag</th>
			<td>
				<form action="/bookmarks/{{ .Bookmark.ID }}/addtag" method="post">
					<div class="form-group">
						<div class="input-group">
							<select class="form-select" name="tag_id">
								{{ range .AllTags }}
									<option value="{{ .ID }}">{{ .Label }}</option>
								{{ end }}
							</select>
							<button type="submit" class="btn input-group-btn btn-primary">Add</button>
						</div>
					</div>
				</form>
			</td>
		</tr>
	</tbody>
</table>
{{ end }}

{{ end }}
`
