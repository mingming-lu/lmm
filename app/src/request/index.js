function get (url, onSuccess, isSync = true) {
  var xhr = new XMLHttpRequest()
  xhr.onload = (event) => {
    var response = JSON.parse(xhr.responseText)
    onSuccess(response, event)
  }
  xhr.open('GET', url, isSync)
  xhr.send()
}

export { get }
