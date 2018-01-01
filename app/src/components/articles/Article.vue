<template>
  <div class="lmm-container lmm-margin" style="margin-left:auto !important; margin-right:auto !important; width:720px">
    <h2 class="lmm-center">{{ title }}</h2>
    <br>
    <div v-html="text" style="text-align:justify"></div>
    <br>
    <p v-if="createdDate === editedDate" class="lmm-right lmm-opacity">Created {{ editedDate }}</p>
    <p v-else class="lmm-right lmm-opacity">Edited {{ createdDate }}</p>
  </div>
</template>

<script>
import axios from 'axios'
import Markdownit from 'markdown-it'
export default {
  data () {
    return {
      title: '',
      text: '',
      createdDate: '',
      editedDate: ''
    }
  },
  created () {
    let pattern = /^\/articles\/(\d)$/g
    let match = pattern.exec(this.$route.path)
    let url = 'http://api.lmm.local' + this.$route.path.replace(/^\/articles\/\d$/, '/article?id=' + match[1])
    let md = new Markdownit({
      html: true,
      typographer: true
    })
    axios.get(url).then((res) => {
      let article = res.data
      this.title = article.title
      this.text = md.render(article.text)
      this.createdDate = article.created_date
      this.editedDate = article.edited_date
    }).catch((e) => {
      console.log(e)
    })
  }
}
</script>
