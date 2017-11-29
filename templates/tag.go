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
</div>
<dl>
  <dt>ID</dt>
  <dd>{{ .Tag.ID }}</dd>
  <dt>Label</dt>
  <dd>{{ .Tag.Label}}</dd>
  <dt>Color</dt>
  <dd>{{ .Tag.Color }}</dd>
  <dt>Description</dt>
  <dd>{{ .Tag.Description }}</dd>
  <dt>Created</dt>
  <dd>{{ .Tag.Created }}</dd>
</dl>
{{ end }}

{{ end }}
`
