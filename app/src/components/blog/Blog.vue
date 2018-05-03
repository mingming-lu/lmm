<template>
  <div class="container">
    <!-- Blog text -->
    <div class="blog" :class="{ 'left': !isMobile, 'mobile': isMobile }">
      <div :class="{container: !isMobile}">
        <p class="title">{{ title }}</p>
        <div class="info">
          <span><i class="fa fa-fw fa-folder-open-o"></i>{{ category.name }}</span>
          <span style="white-space: pre;">  |  </span>
          <span><i class="fa fa-fw fa-calendar-o"></i><span>{{ createdAt }}</span></span>
        </div>
        <div ref="text" class="marked text" v-html="text" v-hljs></div>
        <p v-if="createdAt !== updatedAt" class="info text-right">Updated at {{ updatedAt }}</p>
      </div>
    </div>

    <!-- Blog tags -->
    <div v-if="!isMobile" class="tags">
      <div :class="{container: !isMobile}">
        <h3><i class="fa fa-fw fa-tags"></i>Tags</h3>
        <p>
          <router-link to="" v-for="tag in tags" :key="tag.id" class="link tag">
            {{ tag.name }}
          </router-link>
        </p>
      </div>
    </div>

    <!-- Blog chapters -->
    <div v-if="!isMobile" class="chapters">
      <div :class="{container: !isMobile}">
        <h3><i class="fa fa-fw fa-bookmark-o"></i>Chapters</h3>
        <div ref="progress" class="progress-bar"/>
        <p v-for="subtitle in subtitles" :key="subtitle.name">
          <router-link :to="subtitle.link" @click.native="jumpToHash(subtitle.link)" class="link chapter-item">
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
          document.getElementById(id).className = 'highlight-dispear-trans'
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

<style lang="scss" scoped>
@import '@/assets/scss/styles.scss';
i {
  margin-right: 8px;
}
.container {
  @media screen and (min-width: $width_max_drawer_to_view + 1) {
    padding: 0 32px;
  }
  @media screen and (max-width: $width_max_drawer_to_view) {
    padding: 0 16px;
  }
  .blog {
    float: left;
    width: 66.666%;
    .title {
      color: $color_accent;
      font-weight: 600;
      font-size: 2em;
    }
    .text {
      font-size: 1.1em;
      line-height: 1.8;
      text-align: justify;
    }
    .info {
      // opacity: 0.6;
      color: #777 !important;
    }
  }
  .tags {
    float: right;
    width: 33.3333%;
    .tag {
      display: inline-block;
      background-color: #777;
      padding: 1px 8px;
      margin: 2px;
      border-radius: 2px;
      font-weight: bold;
      font-size: 0.98em;
      color: white !important;
      &:hover {
        background-color: $color_accent;
        opacity: 0.8;
      }
    }
  }
  .chapters {
    float: right;
    position: sticky !important;
    top: 44px;
    width: 33.3333%;
    .chapter-item {
      font-size: 1.1em;
    }
    .chapter-item /deep/ .h3 {
      padding-left: 1em;
    }
    .chapter-item /deep/ .h4 {
      padding-left: 2em;
    }
    .chapter-item /deep/ .h5 {
      padding-left: 3em;
    }
    .chapter-item /deep/ .h6 {
      padding-left: 4em;
    }
  }
}
.mobile {
  width: 100% !important;
}
.progress-bar {
  border-top: 1px solid $color_accent;
  width: 0;
}
.marked /deep/ h2 {
  font-weight: 400;
  color: $color_accent;
  border-bottom: 1px solid #eee;
}
.marked /deep/ h3 {
  font-weight: 400;
  color: $color_accent;
  &:before {
    white-space: pre-wrap;
    border-left: 5px solid $color_accent;
    opacity: 0.6;
    content: '  '; 
  }
}
.marked /deep/ h4 {
  font-size: 400;
  color: $color_text;
}
.marked /deep/ h5 {
  font-size: 400;
  color: $color_text;
}
.marked /deep/ h6 {
  font-size: 400;
  color: $color_text;
}
.marked /deep/ a {
  color: $color_text;
  &:hover {
    opacity: 0.8;
  }
}
.marked /deep/ code {
  background-color: #f1f1f1 !important;
  font-family: Monaco, "Courier", monospace;
}
.marked /deep/ s {
  opacity: 0.5;
}
.marked /deep/ img {
  width: 100%;
}
.marked /deep/ blockquote {
  background: #f9f9f9;
  border-left: 8px solid #ccc;
  margin: 1.5em 0;
  padding: 0.5em 16px;
}
.marked /deep/ table {
  border-bottom: 1px solid #ddd;
  border-top: 1px solid #ddd;
  border-collapse: collapse;
  width: 100%;
}
.marked /deep/ th {
  border-left: 1px solid #ddd;
  border-right: 1px solid #ddd;
  padding: 8px;
  text-align: center;
}
.marked /deep/ tr:nth-child(odd) {
  background-color: #eee;
}
.marked /deep/ td {
  border-left: 1px solid #ddd;
  border-right: 1px solid #ddd;
  padding: 4px 8px;
}
.marked /deep/ .highlighted {
  background-color: lemonchiffon;
}
.marked /deep/ .highlight-dispear-trans {
  color: $color_text;
  background-color: #fff;
  transition: all 0.5s ease-in-out;
}
</style>
