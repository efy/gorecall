package templates

const tagtmpl = `
{{ define "content" }}

<h2>Show tag</h2>

{{ if .Tag }}
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
