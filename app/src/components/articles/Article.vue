<template>
  <div class="lmm-container lmm-margin" style="margin-left:auto !important; margin-right:auto !important; width:720px">
    <h2 class="lmm-center">{{ title }}</h2>
    <br>
    <div v-if="isHTML" v-html="text" style="text-align:justify"></div>
    <div v-else style="text-align:justify">{{ text }}</div>
    <br>
    <p v-if="createdDate === editedDate" class="lmm-right lmm-opacity">Created {{ createdDate }}</p>
    <p v-else class="lmm-right lmm-opacity">Edited {{ editedDate }}</p>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  data () {
    return {
      title: '',
      text: '',
      createdDate: '',
      editedDate: '',
      isHTML: 0
    }
  },
  created () {
    let pattern = /^\/articles\/(\d)$/g
    let match = pattern.exec(this.$route.path)
    let url = 'http://api.lmm.local' + this.$route.path.replace(/^\/articles\/\d$/, '/article?id=' + match[1])
    axios.get(url).then((res) => {
      let article = res.data
      this.title = article.title
      this.text = article.text
      this.createdDate = article.created_date
      this.editedDate = article.edited_date
      this.isHTML = article.is_html
    }).catch((e) => {
      console.log(e)
    })
  }
}
</script>
