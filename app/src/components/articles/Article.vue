<template>
  <div class="lmm-row">

    <!-- Article text -->
    <div class="lmm-left" style="width:75%; display:inline-block">
      <div class="lmm-container lmm-margin lmm-card-4">
        <h2 class="lmm-center">{{ title }}</h2>
        <br>
        <div v-html="text" v-hljs style="text-align:justify"></div>
        <br>
        <p v-if="createdDate === updatedDate" class="lmm-right lmm-opacity">Created at {{ createdDate }}</p>
        <p v-else class="lmm-right lmm-opacity">Updated at {{ updatedDate }}</p>
      </div>
    </div>

    <!-- Article chapters navigation -->
    <div class="lmm-right" style="width:25%; display:inline-block">
      <div class="lmm-container lmm-margin" style="text-align: left">
        <p><b>Chapters</b></p>
        <p v-for="subtitle in subtitles" :key="subtitle.name">
          <router-link :to="subtitle.link" @click.native="jumpToHash(subtitle.link)" class="lmm-white lmm-link lmm-hover">{{ subtitle.name }}</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import Markdownit from 'markdown-it'
export default {
  data () {
    return {
      title: '',
      subtitles: [],
      text: '',
      createdDate: '',
      updatedDate: ''
    }
  },
  created () {
    let match = /^\/articles\/(\d)$/g.exec(this.$route.path)
    let url = 'http://api.lmm.local' + this.$route.path.replace(/^\/articles\/\d$/, '/article?id=' + match[1])
    let md = new Markdownit({
      html: true,
      typographer: true
    })
    axios.get(url).then((res) => {
      let article = res.data
      this.title = article.title
      this.text = md.render(article.text)
      this.createdDate = article.created_date
      this.updatedDate = article.updated_date

      // prepare subtitles and their links
      let results = this.extractSubtitles(this.text, this.$route.path)
      this.text = results[0]
      this.subtitles = results[1]
    }).catch((e) => {
      console.log(e)
    })
  },
  methods: {
    jumpToHash: (hash) => {
      location.href = hash

      // change background color of subtitle for 0.5s
      let match = /^#(.+)$/g.exec(hash)
      if (match !== null && match.length >= 2) {
        let id = match[1]
        console.log(id)
        document.getElementById(id).className = 'lmm-light-grey'
        setTimeout(() => {
          document.getElementById(id).className = 'lmm-white-trans'
        }, 500)
      }
    },
    extractSubtitles: (text, url) => {
      let lines = text.split('\n')
      let subtitles = []

      // regard all h3 as subtitle
      lines.forEach((line, index) => {
        let match = /^<h3>(.+)<\/h3>$/g.exec(line)
        if (match && match.length >= 2) {
          let subtitle = {
            name: match[1],
            link: '#' + match[1]
          }
          subtitles.push(subtitle)
          lines[index] = '<div id="' + match[1] + '">' + line + '</div>'
        }
      })
      return [lines.join('\n'), subtitles]
    }
  }
}
</script>
