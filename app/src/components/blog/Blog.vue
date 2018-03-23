<template>
  <div class="container">
    <!-- Blog text -->
    <div class="marked" :class="{ 'left': !isMobile, 'mobile': isMobile }">
      <div class="container">
        <h1>{{ title }}</h1>
        <span><i class="fa fa-fw fa-folder-open-o"></i><router-link to="" class="white link">{{ category.name }}</router-link></span>
        <span style="white-space: pre;">  |  </span>
        <span><i class="fa fa-fw fa-calendar-o"></i><span class="opacity">{{ createdAt }}</span></span>
        <div ref="text" v-html="text" v-hljs class="text"></div>
        <p v-if="createdAt !== updatedAt" class="text-right opacity">Updated at {{ updatedAt }}</p>
      </div>
    </div>

    <div v-if="!isMobile" class="right">
      <div class="container">
        <h4><i class="fa fa-fw fa-tags"></i>Tags</h4>
        <p>
          <router-link to="" v-for="tag in tags" :key="tag.id" class="link tag">
            {{ tag.name }}
          </router-link>
        </p>
      </div>
    </div>

    <!-- Blog chapters navigation -->
    <div v-if="!isMobile" class="right nav">
      <div class="container chapter">
        <h4><i class="fa fa-fw fa-bookmark-o"></i>Chapters</h4>
        <div ref="progress" class="progress-bar"/>
        <p v-for="subtitle in subtitles" :key="subtitle.name">
          <router-link :to="subtitle.link" @click.native="jumpToHash(subtitle.link)" class="white link chapter-item">
            <div v-html="subtitle.name"></div>
          </router-link>
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
      category: {},
      tags: []
    }
  },
  created () {
    const pattern = /^\/blog\/(\d+)$/g
    const match = pattern.exec(this.$route.path)
    const id = match[1]
    this.blogID = id
    this.fetchBlog()
    this.calcIsMobile()
    window.addEventListener('resize', this.calcIsMobile)
    window.addEventListener('scroll', this.calcProgress)
  },
  watch: {
    text: function () {
      this.$nextTick(() => {
        this.calcProgress()
      })
    }
  },
  beforeDestroy () {
    window.removeEventListener('resize', this.calcIsMobile)
    window.removeEventListener('scroll', this.calcProgress)
  },
  methods: {
    fetchBlog: function () {
      const md = new Markdownit({
        html: true,
        typographer: true
      })
      axios.get('https://api.lmm.im/v1/blog/' + this.blogID).then(res => {
        const blog = res.data

        this.title = blog.title
        this.createdAt = blog.created_at
        this.updatedAt = blog.updated_at
        this.category = blog.category
        this.tags = blog.tags

        // prepare subtitles and their links
        const text = md.render(blog.text)
        const results = this.extractSubtitles(text, this.$route.path)
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

      lines.forEach((line, index) => {
        const h = /^<h(\d)>(.+)<\/h(\d)>$/g.exec(line)
        if (!h || h.length !== 4) {
          return
        }
        let className = ''
        if (h[1] === h[3]) {
          className = 'h' + h[1]
        }
        let subtitle = {
          name: '<div class="' + className + '">' + h[2] + '</div>',
          link: '#' + h[2]
        }
        subtitles.push(subtitle)
        lines[index] = '<div id="' + h[2] + '">' + line + '</div>'
      })
      return [lines.join('\n'), subtitles]
    },
    calcProgress () {
      let el = this.$refs.text
      let progress = ((window.scrollY + window.innerHeight) - el.offsetTop) / (el.offsetHeight)
      progress = progress > 1 ? 100 : progress * 100
      this.$refs.progress.style.width = progress + '%'
    },
    calcIsMobile () {
      this.isMobile = window.innerWidth <= 768
    }
  }
}
</script>

<style scoped>
@import '../../assets/styles/blog.css';
.container .left {
  width: 75%;
}
.container .right {
  width: 25%;
}
.mobile {
  width: 100% !important;
}
.progress-bar {
  border-top: 1px solid deepskyblue;
  width: 0;
}
.marked >>> h3:before {
  white-space: pre-wrap;
  border-left: 5px solid deepskyblue;
  opacity: 0.6;
  content: '  ';
}
.marked >>> h2 {
  font-weight: 400;
  color: deepskyblue;
  border-bottom: 1px solid #eee;
}
.marked >>> h3 {
  font-weight: 400;
  color: deepskyblue;
}
.marked >>> h4 {
  font-size: 400;
  color: deepskyblue;
}
.marked >>> h5 {
  font-size: 400;
  color: deepskyblue;
}
.marked >>> h6 {
  font-size: 400;
  color: deepskyblue;
}
.marked >>> s {
  opacity: 0.5;
}
.marked >>> img {
  width: 100%;
}
.marked >>> blockquote {
  background: #f9f9f9;
  border-left: 8px solid #ccc;
  margin: 1.5em 0;
  padding: 0.5em 16px;
}
.marked >>> table {
  border-bottom: 1px solid #ddd;
  border-top: 1px solid #ddd;
  border-collapse: collapse;
  width: 100%;
}
.marked >>> th {
  border-left: 1px solid #ddd;
  border-right: 1px solid #ddd;
  padding: 8px;
  text-align: center;
}
.marked >>> tr:nth-child(odd) {
  background-color: #eee;
}
.marked >>> td {
  border-left: 1px solid #ddd;
  border-right: 1px solid #ddd;
  padding: 4px 8px;
}
.chapter-item >>> .h3 {
  padding-left: 1em;
}
.chapter-item >>> .h4 {
  padding-left: 2em;
}
.chapter-item >>> .h5 {
  padding-left: 3em;
}
.chapter-item >>> .h6 {
  padding-left: 4em;
}
</style>
