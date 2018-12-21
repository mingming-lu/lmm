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
      <button class="more" @click.prevent="fetchMorePhotos()">See more&hellip;</button>
    </div>

    <div v-if="!hasNext && isPageLoaded" class="center">
      <p class="hint">No more photos.</p>
    </div>

  </div>
</template>

<script>
import LdsEllipsis from '~/components/loadings/LdsEllipsis'
import buildURLEncodedString from '~/assets/js/utils'

const apiPath = '/v1/assets/photos'
const photoFetcher = httpClient => {
  return {
    fetch: page => {
      return httpClient.get(`${apiPath}?perPage=10&page=${page}`)
    }
  }
}
const buildLinks = (obj, path) => {
  return Object.entries(obj)
    .filter(kv => {
      return Boolean(kv[1])
    })
    .map(kv => {
      return {
        rel:  kv[0],
        href: kv[1].replace(apiPath, path)
      }
    })
}

export default {
  components: {
    LdsEllipsis
  },
  head() {
    return {
      title: 'Photos',
      link:  this.links,
    }
  },
  asyncData({$axios, query, route}) {
    const page = Boolean(query.page) ? query.page : 1

    return photoFetcher($axios)
      .fetch(page)
      .then(res => {
        const next = res.data.hasNextPage ? `${apiPath}?page=${Number(page) + 1}` : undefined
        return {
          photos:       res.data.photos,
          left:         res.data.photos.filter((item, index) => index % 2 === 0),
          right:        res.data.photos.filter((item, index) => index % 2 === 1),
          page:         page,
          hasNext:      res.data.hasNextPage,
          isPageLoaded: true,
          wideMode:     false,
          links:        buildLinks({
            next: next,
          }, route.path),
        }
      })
      .catch(e => {
        console.log(e)
      })
  },
  watchQuery: ['page'],
  mounted() {
    window.addEventListener('resize', this.calcIsWideMode)
    this.calcIsWideMode()
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.calcIsWideMode)
  },
  methods: {
    fetchMorePhotos() {
      this.isPageLoaded = false
      const nextPage = Number(this.page) + 1
      photoFetcher(this.$axios)
        .fetch(nextPage)
        .then(res => {
          const photos = res.data.photos
          const next = res.data.hasNextPage ? `?page=${nextPage}` : undefined

          this.photos.push(...photos)
          this.left.push(...photos.filter((item, index) => index % 2 === 0))
          this.right.push(...photos.filter((item, index) => index % 2 === 1))
          this.hasNext = res.data.hasNextPage
          this.isPageLoaded = true
          this.links = buildLinks({next: next}, this.$route.path)
        })
        .catch(e => {
          console.log(e)
        })
    },
    url: name => {
      // TODO, create a plugin to convert name to imageURL
      return `${process.env.ASSET_URL}/photos/${name}`
    },
    calcIsWideMode() {
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
