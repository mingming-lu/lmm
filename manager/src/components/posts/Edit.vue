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
    const pattern = /^\/posts\/(\d+)\/edit$/g
    const match = pattern.exec(this.$route.path)
    this.articleID = match[1]
    this.fetchData()
  },
  methods: {
    onSubmit () {
      if (!this.canSubmit()) {
        alert('no change')
        return
      }
      axios.put('http://api.lmm.local/articles', {
        id: Number(this.articleID),
        title: this.title,
        text: this.text
      }, {
        headers: {
          Authorization: localStorage.getItem('token')
        }
      }).then(res => {
        alert(res.data)
        this.$router.push('/posts')
      }).catch(e => {
        alert(e.response.data)
      })
    },
    canSubmit () {
      this.title = this.title.trim()
      this.text = this.text.trim()
      const isTitleOK = this.title !== this.articleOriginal.title
      const isTextOK = this.text !== this.articleOriginal.text
      return isTitleOK || isTextOK
    },
    marked: (text) => {
      return md.render(text)
    },
    onAddTag (name) {
    },
    onRemoveTag (tag) {
    },
    fetchData () {
      this.fetchBlog()
      this.fetchCategories()
      this.fetchTags()
    },
    fetchBlog () {
      axios.get('http://api.lmm.local/articles?user=1&id=' + this.articleID).then(res => {
        this.articleOriginal = res.data[0]

        this.id = this.articleOriginal.id
        this.title = this.articleOriginal.title
        this.text = this.articleOriginal.text
        this.textPreview = this.marked(this.articleOriginal.text)
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    fetchCategories () {
      axios.all([
        axios.get('http://api.lmm.local/categories?user=1'),
        axios.get('http://api.lmm.local/categories?user=1&article=' + this.articleID)
      ]).then(axios.spread((categories, category) => {
        this.categories = categories.data
        this.categoryID = category.data[0]
      })).catch(e => {
        console.log(e.response.data)
      })
    },
    fetchTags () {
      axios.get('http://api.lmm.local/tags?user=1&article=' + this.articleID).then(res => {
        this.tags = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
    }
  }
}
</script>
