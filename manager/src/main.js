// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import hljs from 'highlight.js'

// stylesheets
import 'highlight.js/styles/github-gist.css'
import '@/../../app/src/assets/styles/common.css'
import '@/assets/styles/common.css'

Vue.config.productionTip = true
Vue.directive('hljs', el => {
  let blocks = el.querySelectorAll('pre code')
  Array.prototype.forEach.call(blocks, hljs.highlightBlock)
})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  template: '<App/>',
  components: { App }
})
