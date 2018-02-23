<template>
  <div class="container">
    <!-- Blog text -->
    <div class="left" :class="{ 'mobile-left': isMobile }">
      <div class="container">
        <h1 class="center">{{ title }}</h1>
        <div v-html="text" v-hljs class="text"></div>
        <p v-if="createdAt === updatedAt" class="text-right opacity">Created at {{ createdAt }}</p>
        <p v-else class="text-right opacity">Updated at {{ updatedAt }}</p>
      </div>
    </div>

    <div v-if="!isMobile" class="right">
      <div class="container">
        <h4><i class="fa fa-fw fa-tags"></i>Tags</h4>
        <router-link to="" v-for="tag in tags" :key="tag.id" class="link tag">
          {{ tag.name }}
        </router-link>
      </div>
    </div>

    <!-- Blog chapters navigation -->
    <div v-if="!isMobile" class="right nav">
      <div class="container chapter">
        <h4><i class="fa fa-fw fa-bookmark-o"></i>Chapters</h4>
        <p v-for="subtitle in subtitles" :key="subtitle.name">
          <router-link :to="subtitle.link" @click.native="jumpToHash(subtitle.link)" class="white link item">{{ subtitle.name }}</router-link>
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
      isMobile: false,
      title: '',
      subtitles: [],
      text: '',
      createdAt: '',
      updatedAt: '',
      category: null,
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
    this.calcIsMobile()
    window.addEventListener('resize', this.calcIsMobile)
  },
  beforeDestroy () {
    window.removeEventListener('resize', this.calcIsMobile)
  },
  methods: {
    fetchBlog: function () {
      const md = new Markdownit({
        html: true,
        typographer: true
      })
      axios.get('https://api.lmm.im/blogs/' + this.blogID).then(res => {
        const blog = res.data
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
      axios.get('https://api.lmm.im/blogs/' + this.blogID + '/category').then(res => {
        this.category = res.data
      }).catch(e => {
        console.log(e)
      })
    },
    fetchTags: function () {
      axios.get('https://api.lmm.im/blogs/' + this.blogID + '/tags').then(res => {
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
        const h = /^<h(\d)>(.+)<\/h(\d)>$/g.exec(line)
        if (!h || h.length !== 4) {
          return
        }
        let prefix = ''
        if (h[1] === h[3]) {
          prefix = '  '.repeat((Number(h[1]) - 2) * 2)
        }
        let subtitle = {
          name: prefix + h[2],
          link: '#' + h[2]
        }
        subtitles.push(subtitle)
        lines[index] = '<div id="' + h[2] + '">' + line + '</div>'
      })
      return [lines.join('\n'), subtitles]
    },
    calcIsMobile () {
      this.isMobile = window.innerWidth <= 768
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
.mobile-left {
  width: 100% !important;
}
.chapter .item{
  white-space: pre;
}
</style>
