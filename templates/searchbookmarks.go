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
		{{ if .SearchQuery }}
		<li class="breadcrumb-item">
			query:
			<a href="/bookmarks/search?q={{ .SearchQuery }}">{{ .SearchQuery }}</a>
		</li>
		{{ end }}
	</ul>
</div>

<div>
	<form action="/bookmarks/search" method="get">
		<div class="form-group">
			<div class="input-group">
				<input type="text" value="{{ .SearchQuery }}" class="form-input" name="q">
				<button type="submit" class="btn input-group-btn btn-primary">Search</button>
			</div>
		</div>
	</form>
</div>

<div>
	<p class="mt-2">{{ .SearchResult.Total }} results ({{ .SearchResult.Took | seconds }} seconds)</p>
</div>

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
		<div class="column col-9">
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
		<div class="column col-2 text-right">
			<button data-delete="/bookmarks/{{ .ID }}" data-redirect="/bookmarks" class="btn btn-sm btn-default">Delete</button>
		</div>
	</div>
{{ end }}
</div>

{{ end }}
`
