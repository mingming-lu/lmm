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
      title: '',
      text: '',
      textPreview: '',
      categoryID: 0,
      categories: [],
      newTagName: '',
      tags: []
    }
  },
  created () {
  },
  methods: {
    onSubmit () {
      axios.post('https://api.lmm.im/v1/blog', {
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
    marked: (text) => {
      return md.render(text)
    },
    onAddTag (name) {
    },
    onRemoveTag (tag) {
    }
  }
}
</script>
