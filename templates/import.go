package templates

const importtmpl = `
{{ define "content" }}

<ul class="tab">
	<li class="tab-item">
		<a href="/settings/account">
			Account
		</a>
	</li>
	<li class="tab-item">
		<a href="/settings/preferences">
			Preferences
		</a>
	</li>
	<li class="tab-item active">
		<a href="/settings/import">
			Import
		</a>
	</li>
	<li class="tab-item">
		<a href="/settings/export">
			Export
		</a>
	</li>
</ul>

<form enctype="multipart/form-data" method="post" action="/settings/import">
  <div class="form-group">
    <label class="form-label" for="bookmarks">File</label>
    <input id="bookmarks_file" class="form-input" type="file" name="bookmarks">
    <p class="form-input-hint">
      A file exported from your browser or bookmarking service.
    </p>
  </div>

  <div class="form-group">
    <button type="submit" class="btn btn-primary">Import</button>
  </div>
</form>

{{ end }}
`

const importsuccesstmpl = `
{{ define "content" }}

<h2>Import</h2>

<div>
  <p>
    Successfully imported {{ len .Bookmarks }} bookmarks.
  </p>
</div>
{{ end }}
`
