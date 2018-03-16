<template>
  <div class="container">
    <div class="left">
      <img v-for="photo in left" :key="photo.name" :src="url(photo.name)" class="picture">
    </div>
    <div class="right">
      <img v-for="photo in right" :key="photo.name" :src="url(photo.name)" class="picture">
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
  },
  methods: {
    url: (name) => {
      return 'https://image.lmm.im/' + name
    }
  }
}
</script>

<style scoped>
div {
  padding: 0 !important;
}
</style>

