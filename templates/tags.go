package templates

const tagstmpl = `
{{ define "content" }}

{{ if not .Tags }}
<div class="rc-empty">
	<div class="empty">
		<div class="empty-icon">
			<i class="icon icon-flag"></i>
		</div>
		<p class="empty-title h5">You have no Tags</p>
		<p class="empty-subtitle">Choose from the actions below to get started</p>
		<div class="empty-action">
			<a href="/tags/new" class="btn btn-primary">Add</a>
		</div>
	</div>
</div>
{{ else }}
	<div class="rc-filter-bar">
		<ul class="breadcrumb">
			<li class="breadcrumb-item">
				<a href="/tags">Tags</a>
			</li>
		</ul>
	</div>
	<div class="rc-tags columns">
	{{ range .Tags }}
		<div class="column col-3">
			<div class="rc-tag">
				<a href="/tags/{{ .Tag.ID }}">
					{{ .Tag.Label }}
				</a>
				<div>
					<span class="label">{{ .Count }}</span>
				</div>
			</div>
		</div>
	{{ end }}
	</div>
{{ end }}

{{ end }}
`
