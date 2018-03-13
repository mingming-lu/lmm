<template>
  <div class="container">
    <router-link to="/image/upload">Upload</router-link>
    <img v-for="image in images" :src="imageURL(image.name)" :key="image.name">
  </div>
</template>

<script>
import axios from 'axios'

export default {
  data () {
    return {
      images: []
    }
  },
  created () {
    this.fetchImages()
  },
  methods: {
    fetchImages () {
      axios.get('https://api.lmm.im/v1/users/1/images').then(res => {
        this.images = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    imageURL: (name) => {
      return 'https://image.lmm.im/' + name
    }
  }
}
</script>
<style scoped>
img {
  width: 100%;
}
</style>
