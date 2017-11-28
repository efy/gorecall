package templates

const newbookmarktmpl = `
{{ define "content" }}

<form action="/bookmarks/new" method="post">
  <div class="form-group">
    <label class="form-label" for="url">URL</label>
    <input id="url" class="form-input" type="text" name="url">
  </div>

  <div class="form-group">
    <button type="submit" class="btn btn-primary">Create</button>
    <button type="reset" class="btn">Reset</button>
  </div>
</form>

{{ end }}
`