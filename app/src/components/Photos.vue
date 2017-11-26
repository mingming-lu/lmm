<template>
  <ul>
    <li class="lmm-half">
      <img v-for="(photo, index) in photos.slice(0, photos.length/2)" :src="photo.url" :key="photo.url" class="lmm-box">
    </li>
    <li class="lmm-half">
      <img v-for="(photo, index) in photos.slice(photos.length/2)" :src="photo.url" :key="photo.url" class="lmm-box">
    </li>
  </ul>
</template>

<script>
import * as request from '@/request'

export default {
  data () {
    request.get('http://localhost:8081/photos', (response) => {
      response.images.forEach((photo) => {
        this.photos.push(photo)
      })
    })
    return {
      photos: []
    }
  }
}
</script>
