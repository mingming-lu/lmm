<template>
  <div class="container">
    <!-- Article text -->
    <div class="left">
      <div class="container">
        <h2 class="center">{{ title }}</h2>
        <div v-html="text" v-hljs class="text"></div>
        <p v-if="createdDate === updatedDate" class="text-right opacity">Created at {{ createdDate }}</p>
        <p v-else class="text-right opacity">Updated at {{ updatedDate }}</p>
      </div>
    </div>

    <div class="right">
      <div class="container">
        <h4>Tags</h4>
        <router-link to="" v-for="tag in tags" :key="tag.id" class="link tag">
          {{ tag.name }}
        </router-link>
      </div>
    </div>

    <!-- Article chapters navigation -->
    <div class="right nav">
      <div class="container">
        <h4>Chapters</h4>
        <p v-for="subtitle in subtitles" :key="subtitle.name">
          <router-link :to="subtitle.link" @click.native="jumpToHash(subtitle.link)" class="white link">{{ subtitle.name }}</router-link>
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
      updatedDate: '',
      category: '',
      tags: []
    }
  },
  created () {
    const pattern = /^\/articles\/(\d+)$/g
    const match = pattern.exec(this.$route.path)
    const id = match[1]
    this.baseURL = 'http://api.lmm.local/users/1/articles/' + id
    this.fetchArticle()
    this.fetchCategories()
    this.fetchTags()
  },
  methods: {
    fetchArticle: function () {
      const md = new Markdownit({
        html: true,
        typographer: true
      })
      axios.get(this.baseURL).then(res => {
        this.title = res.data.title
        this.text = md.render(res.data.text)
        this.createdDate = res.data.created_date
        this.updatedDate = res.data.updated_date

        // prepare subtitles and their links
        const results = this.extractSubtitles(this.text, this.$route.path)
        this.text = results[0]
        this.subtitles = results[1]
      }).catch(e => {
        if (e.response) {
          console.log(e.response.data)
        } else {
          console.log(e)
        }
      })
    },
    fetchCategories: function () {
    },
    fetchTags: function () {
      axios.get(this.baseURL + '/tags').then(res => {
        this.tags = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    jumpToHash: (hash) => {
      location.href = hash
      window.scrollTo(0, document.getElementById(hash.slice(1)).offsetTop - 64)

      // change background color of subtitle for 0.5s
      const match = /^#(.+)$/g.exec(hash)
      if (match !== null && match.length >= 2) {
        const id = match[1]
        document.getElementById(id).className = 'highlighted'
        setTimeout(() => {
          document.getElementById(id).className = 'white-trans'
        }, 500)
      }
    },
    extractSubtitles: (text, url) => {
      let lines = text.split('\n')
      let subtitles = []

      // regard all h3 as subtitle
      lines.forEach((line, index) => {
        const match = /^<h3>(.+)<\/h3>$/g.exec(line)
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

<style scoped>
.container .left {
  width: 75%;
}
.container .right {
  width: 25%;
}
</style>
