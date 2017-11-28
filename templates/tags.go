package templates

const tagstmpl = `
{{ define "content" }}

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

{{ end }}
`
