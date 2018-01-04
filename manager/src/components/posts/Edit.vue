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
      textPreview: '',
      categoryID: 0,
      categories: [],
      tags: '',
      articleOriginal: null
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
      this.articleOriginal = article.data

      this.id = this.articleOriginal.id
      this.title = this.articleOriginal.title
      this.text = this.articleOriginal.text
      this.textPreview = this.marked(this.articleOriginal.text)
      this.categoryID = this.articleOriginal.category_id

      this.categories = categories.data
    })).catch((e) => {
      console.log(e)
    })
  },
  methods: {
    submit () {
      this.text = this.text.trim()
      if (!this.canUpdate()) {
        alert('no change')
        return
      }

      // update article
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
    canUpdate () {
      let titleOK = this.articleOriginal.title !== this.title
      let textOK = this.articleOriginal.text !== this.text
      let categoryIDOK = this.articleOriginal.categoryID !== this.categoryID
      return titleOK || textOK || categoryIDOK
    },
    marked: (text) => {
      return md.render(text)
    }
  }
}
</script>
