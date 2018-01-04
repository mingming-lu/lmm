<template src="@/templates/posts/edit.html">
</template>

<script>
import Markdownit from 'markdown-it'
import axios from 'axios'

let md = new Markdownit({
  html: true,
  typographer: true
})

export default {
  data () {
    return {
      id: '',
      title: '',
      text: '',
      textOriginal: '',
      textPreview: '',
      category: 'Default',
      categories: [],
      tags: ''
    }
  },
  created () {
    let pattern = /^\/posts\/(\d)\/edit$/g
    let match = pattern.exec(this.$route.path)
    let urlArticle = 'http://api.lmm.local' + this.$route.path.replace(pattern, '/article?id=' + match[1])
    this.id = match[1]

    axios.all([
      axios.get(urlArticle),
      axios.get('http://api.lmm.local/articles/categories?user_id=1')
    ]).then(axios.spread((article, categories) => {
      if (this.id !== article.data.id.toString()) {
        throw new Error('id not equal! expected: ' + this.id + ', got: ' + article.data.id)
      }
      this.title = article.data.title
      this.text = article.data.text
      this.textOriginal = article.data.text
      this.textPreview = md.render(article.data.text)
      this.categories = categories.data
      axios.get('http://api.lmm.local/articles/category?id=' + article.data.category_id).then((res) => {
        this.category = res.data.name
      }).catch((e) => {
        console.log(e)
      })
    })).catch((e) => {
      console.log(e)
    })
  },
  methods: {
    submit () {
      this.text = this.text.trim()
      // update article
      if (this.text === this.textOriginal) {
        alert('no change')
        return
      }
      axios.put('http://api.lmm.local/article?id=' + this.id, {
        title: this.title,
        text: this.text,
        category: this.category,
        tags: this.tags
      }).then((res) => {
        this.$router.push('/')
      }).catch((e) => {
        console.log(e)
      })
    },
    marked () {
      this.textPreview = md.render(this.text)
    }
  }
}
</script>
