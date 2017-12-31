<template>
  <div class="lmm-row">
    <div class="lmm-left" style="width:50%">
      <div style="margin:8px">
        <div v-for="photo in photos.slice(0, photos.length/2)" :key="photo.url" class="lmm-box">
          <img :src="photo.url" class="lmm-box">
        </div>
      </div>
    </div>
    <div class="lmm-right" style="width:50%">
      <div style="margin:8px">
        <div v-for="photo in photos.slice(photos.length/2)" :key="photo.url" class="lmm-box">
          <img :src="photo.url" class="lmm-box">
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  data () {
    return {
      photos: []
    }
  },
  created () {
    axios.get('http://api.lmm.local/photos').then((res) => {
      this.photos = res.data.images
    }).catch((e) => {
      console.log(e)
    })
  }
}
</script>
