package templates

const newbookmarktmpl = `
{{ define "content" }}

<div class="rc-flex-center">
	<div class="columns">
		<div class="column col-6 col-mx-auto">
			<form action="/bookmarks/new" method="post">
				<div class="form-group">
					<div class="input-group">
						<span class="input-group-addon addon-lg">URL</span>
						<input id="url" class="form-input input-lg" type="text" name="url">
						<button type="submit" disabled class="btn btn-primary btn-lg input-group-btn">Create</button>
					</div>
				</div>
			</form>
		</div>
	</div>
</div>

{{ end }}
`
