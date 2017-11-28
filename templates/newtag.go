package templates

const newtagtmpl = `
{{ define "content" }}

<div class="rc-flex-center">
	<div class="columns">
		<div class="column col-6 col-mx-auto">
			<form action="/tags" method="post">
				<div class="form-group">
					<label class="form-label" for="label">Label</label>
					<input id="label" class="form-input" type="text" name="label">
				</div>

				<div class="form-group">
					<label class="form-label" for="color">Color</label>
					<input id="color" class="form-input" type="color" name="color">
				</div>

				<div class="form-group">
					<button type="submit" class="btn btn-primary">Create</button>
					<button type="reset" class="btn">Reset</button>
				</div>
			</form>
		</div>
	</div>
</div>

{{ end }}
`
