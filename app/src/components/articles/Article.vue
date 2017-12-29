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
import * as request from '@/request'
export default {
  data () {
    let pattern = /^\/articles\/(\d)$/g
    let match = pattern.exec(this.$route.path)
    let url = 'http://api.lmm.local' + this.$route.path.replace(/^\/articles\/\d$/, '/article?id=' + match[1])
    request.get(url, (response) => {
      this.title = response.title
      this.text = response.text
      this.createdDate = response.created_date
      this.editedDate = response.edited_date
    })
    return {
      title: '',
      text: '',
      createdDate: '',
      editedDate: ''
    }
  }
}
</script>
