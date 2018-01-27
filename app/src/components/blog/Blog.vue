<template>
  <div class="container">
    <!-- Blog text -->
    <div class="left">
      <div class="container">
        <h2 class="center">{{ title }}</h2>
        <div v-html="text" v-hljs class="text"></div>
        <p v-if="createdAt === updatedAt" class="text-right opacity">Created at {{ createdAt }}</p>
        <p v-else class="text-right opacity">Updated at {{ updatedAt }}</p>
      </div>
    </div>

    <div class="right">
      <div class="container">
        <h4><i class="fa fa-fw fa-tags"></i>Tags</h4>
        <router-link to="" v-for="tag in tags" :key="tag.id" class="link tag">
          {{ tag.name }}
        </router-link>
      </div>
    </div>

    <!-- Blog chapters navigation -->
    <div class="right nav">
      <div class="container">
        <h4><i class="fa fa-fw fa-bookmark-o"></i>Chapters</h4>
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
      createdAt: '',
      updatedAt: '',
      category: '',
      tags: []
    }
  },
  created () {
    const pattern = /^\/blog\/(\d+)$/g
    const match = pattern.exec(this.$route.path)
    const id = match[1]
    this.blogID = id
    this.fetchBlog()
    this.fetchCategory()
    this.fetchTags()
  },
  methods: {
    fetchBlog: function () {
      const md = new Markdownit({
        html: true,
        typographer: true
      })
      axios.get('http://api.lmm.im/blog?user=1&id=' + this.blogID).then(res => {
        const blog = res.data[0]
        this.title = blog.title
        this.text = md.render(blog.text)
        this.createdAt = blog.created_at
        this.updatedAt = blog.updated_at

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
    fetchCategory: function () {
      axios.get('http://api.lmm.im/categories?user=1&blog=' + this.blogID).then(res => {
        this.category = res.data[0]
      }).catch(e => {
        console.log(e)
      })
    },
    fetchTags: function () {
      axios.get('http://api.lmm.im/tags?user=1&blog=' + this.blogID).then(res => {
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
