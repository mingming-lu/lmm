import hljs from 'highlight.js'
import Vue from 'vue'

Vue.directive('hljs', (el) => {
  let codeBlocks = el.querySelectorAll('pre code')
  Array.prototype.forEach.call(codeBlocks, hljs.highlightBlock)
})
