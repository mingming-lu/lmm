<template>
  <div class="row">
    <!-- Article text -->
    <div class="left" style="width:75%;">
      <div class="container">
        <h2 class="center">{{ title }}</h2>
        <div v-html="text" v-hljs class="text"></div>
        <p v-if="createdDate === updatedDate" class="text-right opacity">Created at {{ createdDate }}</p>
        <p v-else class="text-right opacity">Updated at {{ updatedDate }}</p>
      </div>
    </div>

    <div class="right" style="width:25%;">
      <div class="container">
        <h4>Category</h4>
        <p>{{ category }}</p>
      </div>
    </div>

    <div class="right" style="width:25%;">
      <div class="container">
        <h4>Tags</h4>
          <router-link to="" v-for="tag in tags" :key="tag.id" class="white link">
            {{ tag.name }}
          </router-link> 
      </div>
    </div>

    <!-- Article chapters navigation -->
    <div class="right nav" style="width:25%;">
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
      tags: [
        {
          id: 1,
          name: '莫哈'
        },
        {
          id: 2,
          name: '闷声发大财'
        }
      ],
      category: '论如何考虑到历史的行程'
    }
  },
  created () {
    let match = /^\/articles\/(\d)$/g.exec(this.$route.path)
    let url = 'http://api.lmm.local' + this.$route.path.replace(/^\/articles\/\d$/, '/article/' + match[1])
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
