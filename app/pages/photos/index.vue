<template>
  <div class="container">
    <div v-if="wideMode" class="content">
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
    <div v-else>
      <div v-for="photo in photos" :key="photo.name">
        <img :src="url(photo.name)">
      </div>
    </div>

    <div v-if="!isPageLoaded" class="center">
      <LdsEllipsis class="fade-in" />
    </div>

    <div v-if="hasNext && isPageLoaded" class="center">
      <br>
      <button class="more" @click.prevent="fetchPhotos()">See more&hellip;</button>
    </div>

    <div v-if="!hasNext && isPageLoaded" class="center">
      <p class="hint">No more photos.</p>
    </div>

  </div>
</template>

<script>
import LdsEllipsis from '~/components/loadings/LdsEllipsis'
export default {
  components: {
    LdsEllipsis
  },
  data () {
    return {
      isPageLoaded: false,
      wideMode: false,
      page: 0,
      hasNext: true,
      left: [],
      right: [],
      photos: []
    }
  },
  created () {
    if (process.browser) {
      window.addEventListener('resize', this.calcIsWideMode)
    }
  },
  mounted () {
    this.calcIsWideMode()
    this.fetchPhotos()
  },
  beforeDestroy () {
    if (process.browser) {
      window.removeEventListener('resize', this.calcIsWideMode)
    }
  },
  methods: {
    fetchPhotos () {
      this.page += 1
      this.isPageLoaded = false
      this.$axios.get(`/v1/assets/photos?perPage=10&page=${this.page}`).then((res) => {
        this.photos.push(...res.data.photos)
        res.data.photos.forEach((photo, index) => {
          if (index % 2 === 0) {
            this.left.push(photo)
          } else {
            this.right.push(photo)
          }
        })
        this.hasNext = res.data.hasNextPage
        this.isPageLoaded = true
      }).catch((e) => {
        console.log(e)
      })
    },
    url: (name) => {
      // TODO, create a plugin to convert name to imageURL
      return `${process.env.ASSET_URL}/photos/${name}`
    },
    calcIsWideMode () {
      if (process.browser) {
        this.wideMode = window.innerWidth >= 800
      }
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
.content {
  display: inline-block;
}
.more {
  border: 1px solid rgba(1, 1, 1, 0.1);
  border-radius: 2px;
  padding: 8px 32px;
  color: $color_text;
  background-color: transparent;
  cursor: pointer;
  font-size: 1.12em;
  &:hover {
    background: transparent;
    border: 1px solid rgba(30, 144, 255, 0.1);
    color: $color_accent;
    outline: none;
  }
}
.hint {
  color: rgba(1, 1, 1, 0.1);
  font-size: 1.12em;
}
.fade-in {
  @include fade_in($opacity: 0.2, $duration: 2s);
}
</style>
