package templates

const searchbookmarkstmpl = `
{{ define "content" }}
<div class="rc-filter-bar">
	<ul class="breadcrumb">
		<li class="breadcrumb-item">
			<a href="/bookmarks">Bookmarks</a>
		</li>
		<li class="breadcrumb-item">
			<a href="/bookmarks/search">Search</a>
		</li>
		<li class="breadcrumb-item">
			query:
			<a href="/bookmarks/search?q={{ .SearchQuery }}">{{ .SearchQuery }}</a>
		</li>
	</ul>
</div>

<div>
	<form action="/bookmarks/search" method="get">
		<div class="form-group">
			<div class="input-group">
				<input type="text" class="form-input" name="q">
				<button type="submit" class="btn input-group-btn btn-primary">Search</button>
			</div>
		</div>
	</form>
</div>
{{ end }}
`
