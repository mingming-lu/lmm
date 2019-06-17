<template>
  <div
    v-if="isMounted" 
    class="container">
    <div
      v-if="wideMode" 
      class="content">
      <div class="left">
        <div :class="{container: wideMode}">
          <a
            v-for="photo in left"
            :key="photo.url"
            :href="photo.url"
          >
            <PhotoItem
              :url="photo.url"
              :tags="photo.tags"/>
          </a>
        </div>
      </div>
      <div class="right">
        <div :class="{container: wideMode}">
          <a
            v-for="photo in right"
            :key="photo.url"
            :href="photo.url"
          >
            <PhotoItem
              :url="photo.url"
              :tags="photo.tags"/>
          </a>
        </div>
      </div>
    </div>
    <div v-else>
      <a
        v-for="photo in photos"
        :key="photo.url"
        :href="photo.url"
      >
        <PhotoItem
          :url="photo.url"
          :tags="photo.tags"/>
      </a>
    </div>

    <div
      v-if="!isPageLoaded" 
      class="center">
      <LdsEllipsis class="fade-in" />
    </div>

    <div
      v-if="hasNext && isPageLoaded" 
      class="center">
      <br>
      <button
        class="more" 
        @click.prevent="fetchMorePhotos()">See more&hellip;</button>
    </div>

    <div
      v-if="!hasNext && isPageLoaded" 
      class="center">
      <p class="hint">No more photos.</p>
    </div>

  </div>
</template>

<script>
import LdsEllipsis from '~/components/loadings/LdsEllipsis'
import PhotoItem from '~/components/photos/photo-item'
import buildURLEncodedString from '~/assets/js/utils'

const apiPath = '/v1/photos'
const photoFetcher = httpClient => {
  return {
    fetch: cursor => {
      return httpClient.get(`${apiPath}?count=10&cursor=${cursor}`)
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
        rel: kv[0],
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
      title: 'Photos'
    }
  },
  asyncData({ $axios, error, query, route }) {
    const cursor = Boolean(query.cursor) ? query.cursor : ''

    return photoFetcher($axios)
      .fetch(cursor)
      .then(res => {
        return {
          photos: res.data.items,
          left: res.data.items.filter((item, index) => index % 2 === 0),
          right: res.data.items.filter((item, index) => index % 2 === 1),
          cursor: res.data.next_cursor,
          hasNext:
            Boolean(res.data.next_cursor) && res.data.next_cursor !== cursor,
          wideMode: false,
          isPageLoaded: true
        }
      })
      .catch(e => {
        if (e.response) {
          error({
            statusCode: e.response.status,
            message: e.response.data
          })
        } else {
          console.log(`failed to fetch photos ${e}`)
        }
      })
  },
  data() {
    return {
      isMounted: false
    }
  },
  watchQuery: ['page'],
  mounted() {
    window.addEventListener('resize', this.calcIsWideMode)
    this.calcIsWideMode()
    this.isMounted = true
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.calcIsWideMode)
  },
  methods: {
    fetchMorePhotos() {
      this.isPageLoaded = false
      photoFetcher(this.$axios)
        .fetch(this.cursor)
        .then(res => {
          this.hasNext =
            Boolean(res.data.next_cursor) &&
            res.data.next_cursor !== this.cursor

          this.cursor = res.data.next_cursor

          const photos = res.data.items
          this.photos.push(...photos)
          this.left.push(...photos.filter((item, index) => index % 2 === 0))
          this.right.push(...photos.filter((item, index) => index % 2 === 1))

          this.isPageLoaded = true
        })
        .catch(e => {
          console.log(e)
        })
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
