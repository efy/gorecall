package templates

const bookmarkstmpl = `
{{ define "content" }}

{{ if not .Bookmarks }}
<div class="rc-empty">
	<div class="empty">
		<div class="empty-icon">
			<i class="icon icon-bookmark"></i>
		</div>
		<p class="empty-title h5">Your library is empty</p>
		<p class="empty-subtitle">Choose from the actions below to get started</p>
		<div class="empty-action">
			<a href="/settings/import" class="btn btn-primary">Import</a>
			<a href="/bookmarks/new" class="btn btn-primary">Add</a>
		</div>
	</div>
</div>
{{ else }}
<div class="rc-filter-bar">
	<ul class="breadcrumb">
		<li class="breadcrumb-item">
			<a href="/bookmarks">Bookmarks</a>
		</li>
	</ul>
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
