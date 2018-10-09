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
      axios.put(process.env.API_URL_BASE + '/v1/blog/' + this.blogID, {
        title: this.title,
        text: this.text
      }, {
        headers: {
          Authorization: 'Bearer ' + localStorage.getItem('token')
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
    onSetCategory () {
      axios.put(process.env.API_URL_BASE + '/v1/blog/' + this.blogID + '/category', {
        id: this.categoryID
      }, {
        headers: {
          Authorization: 'Bearer ' + localStorage.getItem('token')
        }
      }).then(res => {
        alert(res.data)
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    onAddTag (name) {
      axios.post(process.env.API_URL_BASE + '/v1/blog/' + this.blogID + '/tags', {
        name: name
      }, {
        headers: {
          Authorization: 'Bearer ' + localStorage.getItem('token')
        }
      }).then(res => {
        alert(res.data)
        this.fetchTags()
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    onRemoveTag (tag) {
      axios.delete(process.env.API_URL_BASE + '/v1/blog/' + this.blogID + '/tags/' + tag.id, {
        headers: {
          Authorization: 'Bearer ' + localStorage.getItem('token')
        }
      }).then(res => {
        alert('deleted')
        this.fetchTags()
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    fetchData () {
      this.fetchBlog()
      this.fetchCategories()
      this.fetchTags()
    },
    fetchBlog () {
      axios.get(process.env.API_URL_BASE + '/v1/blog/' + this.blogID).then(blog => {
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
      axios.get(process.env.API_URL_BASE + '/v1/categories').then(res => {
        this.categories = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
      axios.get(process.env.API_URL_BASE + '/v1/blog/' + this.blogID + '/category').then(res => {
        this.categoryID = res.data.id
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    fetchTags () {
      axios.get(process.env.API_URL_BASE + '/v1/blog/' + this.blogID + '/tags').then(res => {
        this.tags = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
    }
  }
}
</script>
