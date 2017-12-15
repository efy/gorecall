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
	
	<div class="columns">
		<div class="column col-4">
			<div class="form-group">
				<label class="form-switch">
					<input type="checkbox" name="folders_to_tags">
					<i class="form-icon"></i>
					Folders to tags?
				</label>
				<p class="form-input-hint">
					Faltten out hierarchical folders and add them as tags
				</p>
			</div>
		</div>

		<div class="column col-4">
			<div class="form-group">
				<label class="form-switch">
					<input type="checkbox" name="webinfo">
					<i class="form-icon"></i>
					Run webinfo for each bookmark
				</label>
				<p class="form-input-hint">
					Requests each URL and adds additonal information
					Will significantly increase time to import
				</p>
			</div>
		</div>

		<div class="column col-4">
			<div class="form-group">
				<label class="form-switch">
					<input type="checkbox" name="index">
					<i class="form-icon"></i>
					Index each bookmark
				</label>
				<p class="form-input-hint">
					Will increase time to import
				</p>
			</div>
		</div>
	</div>


  <div class="form-group">
    <button type="submit" class="btn btn-primary">Import</button>
  </div>
</form>

<div class="divider" data-content="Help"></div>
<div class="columns">
	<div class="column col-6">
		<div class="rc-import-help-block">
			<h4>Exporting from Google Chrome</h4>
			<p>
				<ol>
					<li>Go to your bookmarks manager (CTRL+SHIFT+O)</li>
					<li>Open "Organize" dropdown</li>
					<li>Click "Export bookmarks to HTML file..."</li>
				</ol>
			</p>
		</div>
	</div>
	<div class="column col-6">
		<div class="rc-import-help-block">
			<h4>Exporting from Mozilla Firefox</h4>
			<p>
				<ol>
					<li>Go to your bookmarks library (CTRL+SHIFT+B)</li>
					<li>Open "Import and Backup" dropdown</li>
					<li>Click "Export bookmarks to HTML..."</li>
				</ol>
			</p>
		</div>
	</div>
	<div class="column col-6">
		<div class="rc-import-help-block">
			<h4>Exporting from Safari</h4>
			<p>
				<ol>
					<li>Open your "File" menu</li>
					<li>Click "Export Bookmarks..."</li>
				</ol>
			</p>
		</div>
		<div class="rc-import-help-block">
			<h4>Note about Microsoft Edge & Opera</h4>
			<p>
				Unfortunately these browsers do not support exporting bookmarks natively - although in the case of Opera there's a possibility to use 3rd party browser plugin.
			</p>
		</div>
	</div>
	<div class="column col-6">
		<div class="rc-import-help-block">
			<h4>Exporting from Internet Explorer 11</h4>
			<p>
				<ol>
					<li>Open your favorites (button with star icon)</li>
					<li>Open dropdown right next to "Add to favorites"</li>
					<li>Click "Import and export"</li>
					<li>Select "Export to a file"</li>
					<li>In next step select "Favorites"</li>
					<li>In next step select which bookmarks folders you want to export</li>
					<li>In final step choose folder where you want your file and click on "Export"</li>
				</ol>
			</p>
		</div>
	</div>
	<div class="column col-6">
	</div>
</div>

{{ end }}
`

const importsuccesstmpl = `
{{ define "content" }}

<h2>Import</h2>

<div>
  <p>
    Successfully imported {{ .Report.SuccessCount }} bookmarks. {{ .Report.FailureCount }} failed (see table below).
  </p>
	<div class="rc-import-errors">
		<table class="table">
			<thead>
				<th>
					Title
				</th>
				<th>
					Error
				</th>
			</thead>
			<tbody>
			{{ range .Report.Results }}
				{{ if .Error }}
				<tr>
					<td>
						{{ .Bookmark.Title }}
					</td>
					<td>
						{{ .Error.Error }}
					</td>
				</tr>
				{{ end }}
			{{ end }}
			</tbody>
		</table>
	</div>
</div>
{{ end }}
`
