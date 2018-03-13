<template>
  <div class="container">
    <div class="left">
      <img v-for="photo in left" :key="photo.url" :src="photo.url" class="picture">
    </div>
    <div class="right">
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
    axios.get('https://api.lmm.im/v1/users/1/images/photos').then((res) => {
      res.data.forEach((image, index) => {
        if (index % 2 === 0) {
          this.left.push(image)
        } else {
          this.right.push(image)
        }
      })
    }).catch((e) => {
      console.log(e)
    })
  }
}
</script>

<style scoped>
div {
  padding: 0 !important;
}
</style>

