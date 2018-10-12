import Vue from 'vue'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-gist.css'

Vue.directive('hljs', el => {
  let codeBlocks = el.querySelectorAll('pre code')
  console.log(codeBlocks)
  Array.prototype.forEach.call(codeBlocks, hljs.highlightBlock)
})
