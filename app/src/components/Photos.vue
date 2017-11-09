<template>
  <div class="w3-row-padding" style="margin-bottom:128px">
    <div class="w3-half">
      <img v-for="(photo, index) in photos.slice(0, Math.floor(photos.length/2))" :src="photo.url" :key="photo.url" style="width:100%">
    </div>
    <div class="w3-half">
      <img v-for="(photo, index) in photos.slice(Math.floor(photos.length/2))" :src="photo.url" :key="photo.url" style="width:100%">
    </div>
  </div>
</template>

<script>
import * as request from '@/request'

export default {
  data () {
    request.get('http://localhost:8081/photos', (response) => {
      response.result.images.forEach((photo) => {
        this.photos.push(photo)
      })
    })
    return {
      photos: []
    }
  }
}
</script>
