package templates

const tagtmpl = `
{{ define "content" }}

{{ if .Tag }}
<div class="rc-filter-bar">
	<ul class="breadcrumb">
		<li class="breadcrumb-item">
			<a href="/tags">Tags</a>
		</li>
		<li class="breadcrumb-item">
			<a href="/tags/{{ .Tag.ID }}">{{ .Tag.Label }}</a>
		</li>
	</ul>
	<div>
		<button data-delete="/tags/{{ .Tag.ID }}" data-redirect="/tags" class="btn btn-sm btn-default">Delete</button>
	</div>
</div>
{{ end }}

{{ if .Bookmarks }}
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
{{ end }}

{{ if .Pagination }}
<div class="rc-pagination-container">
<ul class="pagination">
	{{ if lt .Pagination.Prev 0 }}
		<li class="page-item disabled"><a href="#">Previous</a></li>
	{{ else }}
		<li class="page-item"><a href="/tags/{{ .Tag.ID }}?per_page={{ .Pagination.PerPage }}&page={{ .Pagination.Prev }}">Previous</a></li>
	{{ end }}

	<li class="page-item active">
		<a href="/tags/{{ .Tag.ID }}?per_page={{ .Pagination.PerPage }}&page={{ .Pagination.Current }}">{{ .Pagination.Current }}</a>
	</li>

	{{ $save := .Pagination }}
	{{ $root := . }}
	{{ range $page := .Pagination.List }}
		{{ if lt $page $save.Last }}
			<li class="page-item">
				<a href="/tags/{{ $root.Tag.ID }}?per_page={{ $save.PerPage }}&page={{ $page}}">{{ $page }}</a>
			</li>
		{{ end }}
	{{ end }}

	{{ if ne .Pagination.Last .Pagination.Current }}
		<li class="page-item">
			<span>...</span>
		</li>
		<li class="page-item">
			<a href="/tags/{{ .Tag.ID }}?per_page={{ .Pagination.PerPage }}&page={{ .Pagination.Last }}">{{ .Pagination.Last }}</a>
		</li>
	{{ end }}

	{{ if gt .Pagination.Next .Pagination.Last }}
		<li class="page-item disabled"><a href="#">Next</a></li>
	{{ else }}
		<li class="page-item"><a href="/tags/{{ .Tag.ID }}?per_page={{ .Pagination.PerPage }}&page={{ .Pagination.Next }}">Next</a></li>
	{{ end }}
</ul>
</div>
{{ end }}

{{ end }}
`
