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
    let urlArticle = 'http://api.lmm.local' + this.$route.path.replace(pattern, '/article/' + id)

    axios.all([
      axios.get(urlArticle),
      axios.get('http://api.lmm.local/articles/1/categories')
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
    onSubmit () {
      this.text = this.text.trim()

      if (!confirm('Are you sure you want to submit?')) {
        return
      }

      if (!this.canSubmit()) {
        alert('no change')
        return
      }

      // update article
      axios.put('http://api.lmm.local/article/' + this.id, {
        user_id: 1,
        title: this.title,
        text: this.text,
        category_id: this.categoryID,
        tags: this.tags
      }).then((res) => {
        this.$router.push('/posts')
      }).catch((e) => {
        console.log(e)
      })
    },
    canSubmit () {
      if (!this.articleOriginal) {
        return false
      }
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
