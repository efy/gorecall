package templates

const exporttmpl = `
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
	<li class="tab-item">
		<a href="/settings/import">
			Import
		</a>
	</li>
	<li class="tab-item active">
		<a href="/settings/export">
			Export
		</a>
	</li>
</ul>

<form enctype="multipart/form-data" method="post" action="/settings/export">
</form>

{{ end }}
`
