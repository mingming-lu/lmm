// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-gist.css'

Vue.config.productionTip = true

Vue.directive('hljs', (el) => {
  let codeBlocks = el.querySelectorAll('pre code')
  Array.prototype.forEach.call(codeBlocks, hljs.highlightBlock)
})

function findAncestor (el, cls) {
  while ((el = el.parentElement) && !el.classList.contains(cls)) {
  }
  return el
}

Vue.directive('nav', (el) => {
  let navs = el.querySelectorAll('.lmm-nav-container .lmm-nav')
  Array.prototype.forEach.call(navs, (el) => {
    let container = findAncestor(el, 'lmm-nav-container')
    if (!container) {
      return
    }
    el.addEventListener('scroll', () => {
      if (container.getBoundingClientRect().top <= 0) {
        el.style = 'position:fixed;'
      } else {
        el.style = 'position:inherit;'
      }
    })
  })
})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  template: '<App/>',
  components: { App }
})
