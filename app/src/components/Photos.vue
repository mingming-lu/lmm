<template>
  <div class="container">

    <div v-if="wideMode">
      <div class="left">
        <div :class="{container: wideMode}">
          <img v-for="photo in left" :key="photo.name" :src="url(photo.name)">
        </div>
      </div>
      <div class="right">
        <div :class="{container: wideMode}">
          <img v-for="photo in right" :key="photo.name" :src="url(photo.name)">
        </div>
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
      this.wideMode = window.innerWidth >= 800
    }
  }
}
</script>

<style lang="scss" scoped>
@import '@/assets/scss/styles.scss';
img {
  display: block;
  width: 100%;
  @media screen and (min-width: 800px) {
    margin-bottom: 64px;
  }
  @media screen and (max-width: 799px) {
    margin-bottom: 32px;
  }
  @media screen and (max-width: $max_width_device) {
    margin-bottom: 16px;
  }
}
.container {
  margin-bottom: -64px;
  padding: 0 32px;
  @media screen and (max-width: $max_width_device) {
    margin-bottom: -16px;
    padding: 0 16px;
  }
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

