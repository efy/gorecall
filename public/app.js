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
        console.log(res)
        return res.json()
      })
      .then(function(json) {
        webinfo.innerHTML = webinfo_tmpl(json)
        button.disabled = false
      })
      .catch(function(error){
        console.log(error)
        webinfo.innerHTML = webinfo_error_tmpl(error)
        button.disabled = true
      })
  }, 300))

  function webinfo_error_tmpl(obj) {
    return `
      <div class="rc-link-preview">
        <div class="rc-link-preview__error">
          ${obj}
        </div>
      </div>
    `
  }

  function webinfo_tmpl(obj) {
    return `
    <div class="rc-link-preview">
      <div class="rc-link-preview__image">
      </div>
      <div class="rc-link-preview__info">
        <div class="text-bold text-ellipsis">${obj.title}</div>
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

