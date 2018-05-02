<template>
  <div class="container">

    <div v-if="wideMode">
      <div class="left">
        <img v-for="photo in left" :key="photo.name" :src="url(photo.name)">
      </div>
      <div class="right">
        <img v-for="photo in right" :key="photo.name" :src="url(photo.name)">
      </div>
    </div>

    <img v-if="!wideMode" v-for="photo in photos" :key="photo.name" :src="url(photo.name)">
  </div>
</template>

<script>
import axios from 'axios'
export default {
  data () {
    return {
      wideMode: false,
      left: [],
      right: [],
      photos: []
    }
  },
  created () {
    this.calcIsWideMode()
    window.addEventListener('resize', this.calcIsWideMode)
    axios.get('https://api.lmm.im/v1/users/1/images/photos').then((res) => {
      this.photos.push(...res.data)
      res.data.forEach((photo, index) => {
        if (index % 2 === 0) {
          this.left.push(photo)
        } else {
          this.right.push(photo)
        }
      })
    }).catch((e) => {
      console.log(e)
    })
  },
  beforeDestroy () {
    window.removeEventListener('resize', this.calcIsWideMode)
  },
  methods: {
    url: (name) => {
      return 'https://image.lmm.im/' + name
    },
    calcIsWideMode () {
      this.wideMode = window.innerWidth > 800
    }
  }
}
</script>

<style lang="scss" scoped>
img {
  display: block;
  width: 100%;
}
.container {
  .left {
    float: left;
    width: 50%;
  }
  .right {
    float: right;
    width: 50%;
  }
}
</style>

