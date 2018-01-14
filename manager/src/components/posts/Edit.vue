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
      newTagName: '',
      tags: [],
      articleOriginal: null
    }
  },
  created () {
    this.fetchData()
  },
  methods: {
    onSubmit () {
      this.text = this.text.trim()

      if (!this.canSubmit()) {
        alert('no change')
        return
      }

      if (!confirm('Are you sure you want to submit?')) {
        return
      }

      // update article
      axios.put('http://api.lmm.local/users/1/articles/' + this.id, {
        title: this.title,
        text: this.text
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
      return titleOK || textOK
    },
    marked: (text) => {
      return md.render(text)
    },
    onAddTag (name) {
      if (!name.trim()) {
        return
      }
      axios.post('http://api.lmm.local/article/tags', [{
        user_id: 1,
        article_id: this.id,
        name: this.newTagName
      }]).then((res) => {
        alert(res.data)
        this.fetchData()
        this.newTagName = ''
      }).catch((e) => {
        console.log(e)
        if (e.response) {
          console.log(e.response)
        }
      })
    },
    onRemoveTag (tag) {
      axios.delete('http://api.lmm.local/article/tags/' + tag.id).then((res) => {
        alert(res.data)
        this.fetchData()
      }).catch((e) => {
        console.log(e)
      })
    },
    fetchData () {
      const pattern = /^\/posts\/(\d+)\/edit$/g
      const match = pattern.exec(this.$route.path)
      const id = match[1]
      const baseURL = 'http://api.lmm.local/users/1/articles/' + id

      axios.all([
        axios.get(baseURL)
        // axios.get('http://api.lmm.local/articles/1/categories'),
        // axios.get('http://api.lmm.local/article/' + id + '/tags')
      ]).then(axios.spread((article) => {
        if (id !== article.data.id.toString()) {
          throw new Error('id not equal! expected: ' + this.id + ', got: ' + article.data.id)
        }
        this.articleOriginal = article.data

        this.id = this.articleOriginal.id
        this.title = this.articleOriginal.title
        this.text = this.articleOriginal.text
        this.textPreview = this.marked(this.articleOriginal.text)
        // this.categoryID = this.articleOriginal.category_id

        // this.categories = categories.data
        // this.tags = tags.data
      })).catch((e) => {
        console.log(e)
      })
    }
  }
}
</script>
