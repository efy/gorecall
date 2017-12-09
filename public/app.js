(function (d, w) {

  var url_input = d.querySelector("#url")
  var webinfo = d.querySelector("#webinfo")
  var button = d.querySelector("#create")

  if(url_input == null || url_input.tagName !== "INPUT") {
    return
  }

  url_input.addEventListener('keyup', debounce(function (e){
    url = e.target.value

    fetch("/api/webinfo?url=" + url)
      .then(function(res){
        return res.json()
      })
      .then(function(json) {
        if(json.code && json.code !== 200) {
          webinfo.innerHTML = webinfo_error_tmpl(json)
          return
        }
        webinfo.innerHTML = webinfo_tmpl(json)

        d.querySelector('input[name="url"]').value = url
        d.querySelector('input[name="title"]').value = json.title

        button.disabled = false
      })
      .catch(function(error){
        webinfo.innerHTML = webinfo_error_tmpl(error)
        button.disabled = true
      })
  }, 300))

  function webinfo_error_tmpl(obj) {
    return `
      <div class="rc-link-preview">
        <div class="rc-link-preview__error">
          ${obj.message}
        </div>
      </div>
    `
  }

  function webinfo_tmpl(obj) {
    if(obj.opengraph) {
      return og_tmpl(obj.opengraph)
    }
    return default_tmpl(obj)
  }

  function default_tmpl(obj) {
    return `
    <div class="rc-link-preview">
      <div class="rc-link-preview__image">
        <img src="${obj.cover}">
      </div>
      <div class="rc-link-preview__info">
        <div class="text-bold text-ellipsis">${obj.title}</div>
      </div>
    </div>
    `
  }

  function og_tmpl(obj) {
    return `
    <div class="rc-link-preview">
      <div class="rc-link-preview__image">
        <img src="${obj.image.source}">
      </div>
      <div class="rc-link-preview__info">
        <div class="text-bold text-ellipsis">${obj.title}</div>
        <p>
          ${obj.description}
        </p>
      </div>
    </div>
    `
  }

  function debounce(fn, w) {
    var timeout

    return function() {
      var ctx = this
      var args = arguments
      var later = function () {
        timeout = null
        fn.apply(ctx,args)
      }
      clearTimeout(timeout)
      timeout = setTimeout(later, w)
    }
  }

})(document, window);


// Delete request handling.
// Requests with the data-delete attribute will prompt the 
// user and send a delete request to the supplied endpoint.
// If a data-redirect is supplied and the response code of
// the delete request 2xx the url will be replaced with the
// given redirect.
(function(d, w) {

  var actions = d.querySelectorAll("[data-delete]")

  for(var i = 0; i < actions.length; i++) {
    (function(action) {
      action.addEventListener("click", function(e){
        if(confirm("Are you sure you want to delete this item")) {
          var endpoint = e.target.dataset.delete
          var redirect = e.target.dataset.redirect || "/"

          fetch(endpoint, {
            credentials: 'same-origin',
            method: 'delete'
          }).then(function(res) {
            if(res.status >= 200 && res.status <= 299) {
              window.location = redirect
            }
          })
        }
      })
    })(actions[i]);
  }

})(document, window);

