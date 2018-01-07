<template>
  <div class="container">
    <div class="left" style="width:50%">
        <img v-for="photo in left" :key="photo.url" :src="photo.url" class="picture">
    </div>
    <div class="right" style="width:50%">
        <img v-for="photo in right" :key="photo.url" :src="photo.url" class="picture">
    </div>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  data () {
    return {
      left: [],
      right: []
    }
  },
  created () {
    axios.get('http://api.lmm.local/photos').then((res) => {
      for (let i = 0; i < res.data.images.length; i++) {
        if (i % 2 === 0) {
          this.left.push(res.data.images[i])
        } else {
          this.right.push(res.data.images[i])
        }
      }
    }).catch((e) => {
      console.log(e)
    })
  }
}
</script>

<style scoped>
div {
  padding-left: 0 !important;
  padding-right: 0 !important;
}
</style>

