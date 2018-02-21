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
    axios.get('https://api.lmm.im/users/1/photos').then((res) => {
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
  padding-left: 0 !important;
  padding-right: 0 !important;
}
</style>

