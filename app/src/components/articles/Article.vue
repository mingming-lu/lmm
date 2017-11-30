<template>
  <div class="lmm-container" style="margin: 0 128px">
    <h2 class="lmm-center">{{ title }}</h2>
    <br>
    <p style="text-align:justify">{{ text }}</p>
    <br>
    <p v-if="editedDate" class="lmm-right lmm-opacity">Edited {{ editedDate }} at {{ editedTime }}</p>
    <p v-else class="lmm-right lmm-opacity">Posted {{ postedDate }} at {{ postedTime }}</p>
  </div>
</template>

<script>
import * as request from '@/request'
export default {
  data () {
    let url = 'http://localhost:8081' + this.$route.path
    request.get(url, (response) => {
      this.title = response.title
      this.text = response.text
    })
    return {
      title: '',
      text: '',
      postedDate: '',
      postedTime: '',
      editedDate: '',
      editedTime: ''
    }
  }
}
</script>
