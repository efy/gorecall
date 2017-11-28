package templates

const logintmpl = `
{{ define "content" }}

<h2 class="text-center">Login</h2>

<form method="post" action="/login">
  <div class="form-group">
    <label class="form-label" for="username">Username</label>
    <input id="username" class="form-input" type="text" name="username">
  </div>

  <div class="form-group">
    <label class="form-label" for="password">Password</label>
    <input id="password" class="form-input" type="password" name="password">
  </div>

  <div class="form-group">
    <label for="remember_me" class="form-checkbox">
      <input id="remember_me" type="checkbox" name="remember_me">
      <i class="form-icon"></i>Remember me?
    </label>
  </div>

  <div class="form-group">
    <button type="submit" class="btn btn-primary">Login</button>
  </div>
</form>

{{ end }}
`
