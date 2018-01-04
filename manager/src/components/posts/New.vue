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
      title: '',
      text: '',
      textPreview: '',
      categoryID: 0,
      categories: [],
      tags: ''
    }
  },
  created () {
    axios.get('http://api.lmm.local/articles/categories?user_id=1').then((res) => {
      this.categories = res.data
    }).catch((e) => {
      console.log(e)
    })
  },
  methods: {
    submit () {
      this.text = this.text.trim()
      if (!confirm('Are you sure you want to submit?')) {
        return
      }
      if (this.categoryID === 0) {
        alert('must select one category')
        return
      }

      axios.post('http://api.lmm.local/article', {
        user_id: 1,
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
