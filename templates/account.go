package templates

const accounttmpl = `
{{ define "content" }}

<ul class="tab">
	<li class="tab-item active">
		<a href="/settings/account">
			Account
		</a>
	</li>
	<li class="tab-item">
		<a href="/settings/preferences">
			Preferences
		</a>
	</li>
	<li class="tab-item">
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

<form action="/account/edit" method="post">
  <div class="form-group">
    <label class="form-label" for="username">Username</label>
    <input id="username" class="form-input" value="{{ .User.Username }}" type="text" name="username">
  </div>

  <div class="form-group">
    <button type="submit" class="btn btn-primary">Update</button>
  </div>
</form>

{{ end }}
`
