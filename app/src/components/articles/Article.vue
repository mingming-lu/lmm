<template>
  <div class="lmm-container lmm-margin" style="margin-left:auto !important; margin-right:auto !important; width:720px">
    <h2 class="lmm-center">{{ title }}</h2>
    <br>
    <p style="text-align:justify">{{ text }}</p>
    <br>
    <p v-if="createdDate === editedDate" class="lmm-right lmm-opacity">Created {{ editedDate }}</p>
    <p v-else class="lmm-right lmm-opacity">Edited {{ createdDate }}</p>
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
      editedDate: ''
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
      this.createdDate = article.createdDate
      this.editedDate = article.editedDate
    }).catch((e) => {
      console.log(e)
    })
  }
}
</script>
