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
      categoryID: 0,
      categories: [],
      tags: ''
    }
  },
  created () {
    let pattern = /^\/posts\/(\d)\/edit$/g
    let match = pattern.exec(this.$route.path)
    let id = match[1]
    let urlArticle = 'http://api.lmm.local' + this.$route.path.replace(pattern, '/article?id=' + id)

    axios.all([
      axios.get(urlArticle),
      axios.get('http://api.lmm.local/articles/categories?user_id=1')
    ]).then(axios.spread((article, categories) => {
      if (id !== article.data.id.toString()) {
        throw new Error('id not equal! expected: ' + this.id + ', got: ' + article.data.id)
      }
      this.id = article.data.id
      this.title = article.data.title
      this.text = article.data.text
      this.textOriginal = article.data.text
      this.textPreview = this.marked(article.data.text)
      this.categoryID = article.data.category_id
      this.categories = categories.data
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
      if (this.categoryID === 0) {
        alert('must select one category')
        return
      }
      axios.put('http://api.lmm.local/article', {
        id: this.id,
        title: this.title,
        text: this.text,
        category_id: this.categoryID,
        tags: this.tags
      }).then((res) => {
        this.$router.push('/')
      }).catch((e) => {
        console.log(e)
      })
    },
    marked: (text) => {
      return md.render(text)
    }
  }
}
</script>
