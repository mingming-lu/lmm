<template src="@/templates/blog/edit.html">
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
      blogOriginal: null
    }
  },
  created () {
    const pattern = /^\/blog\/(\d+)\/edit$/g
    const match = pattern.exec(this.$route.path)
    this.blogID = match[1]
    this.fetchData()
  },
  methods: {
    onSubmit () {
      if (!this.canSubmit()) {
        alert('no change')
        return
      }
      axios.put('http://api.lmm.local/blogs/' + this.blogID, {
        title: this.title,
        text: this.text
      }, {
        headers: {
          Authorization: localStorage.getItem('token')
        }
      }).then(res => {
        alert(res.data)
        this.$router.push('/blog')
      }).catch(e => {
        alert(e.response.data)
      })
    },
    canSubmit () {
      this.title = this.title.trim()
      this.text = this.text.trim()
      const isTitleOK = this.title !== this.blogOriginal.title
      const isTextOK = this.text !== this.blogOriginal.text
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
      axios.get('http://api.lmm.local/blogs/' + this.blogID).then(blog => {
        this.blogOriginal = blog.data

        this.id = this.blogOriginal.id
        this.title = this.blogOriginal.title
        this.text = this.blogOriginal.text
        this.textPreview = this.marked(this.blogOriginal.text)
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    fetchCategories () {
      axios.all([
        axios.get('http://api.lmm.local/categories?user=1'),
        axios.get('http://api.lmm.local/categories?user=1&blog=' + this.blogID)
      ]).then(axios.spread((categories, category) => {
        this.categories = categories.data
        this.categoryID = category.data[0].id
      })).catch(e => {
        console.log(e.response.data)
      })
    },
    fetchTags () {
      axios.get('http://api.lmm.local/tags?user=1&blog=' + this.blogID).then(res => {
        this.tags = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
    }
  }
}
</script>
