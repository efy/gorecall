package templates

const layouttmpl = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>Home</title>
    <link rel="stylesheet" href="/public/style.css">
		<script async src="/public/app.js"></script>
  </head>
  <body>
    <div class="wrapper">
			<div class="columns vh">
				<div class="column col-2 vh">
					<div class="rc-sidebar">
						<img src="/public/logo.svg" class="logo">
						<ul class="nav">
							{{ if .Authenticated }}
								<li class="divider" data-content="Links"></li>
								<li class="nav-item">
									<a href="/bookmarks">All</a>
								</li>
								<li class="nav-item">
									<a href="/bookmarks/new">Add Link</a>
									<i class="icon icon-plus"></i>
								</li>
								<li class="divider" data-content="Tags"></li>
								<li class="nav-item">
									<a href="/tags">All</a>
								</li>
								<li class="nav-item">
									<a href="/bookmarks">Untagged</a>
								</li>
								<li class="nav-item">
									<a href="/tags/new">Add Tag</a>
									<i class="icon icon-plus"></i>
								</li>
								<li class="divider"></li>
								<li class="nav-item">
									<a href="/settings/account">Settings</a>
									<i class="icon icon-forward"></i>
								</li>
								<li class="nav-item">
									<a href="/logout">Logout</a>
									<i class="icon icon-shutdown"></i>
								</li>
							{{ end }}
						</ul>
					</div>
				</div>
				<div class="column col-10">
					<div class="content-area">
						{{ block "content" . }} {{ end }}
					</div>
				</div>
			</div>
    </div>
  </body>
</html>
`
