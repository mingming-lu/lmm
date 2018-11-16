<template>
  <div class="container">
    <!-- Articles -->
    <div
      :class="{ 'desktop': !isMobile, 'mobile': isMobile }"
      class="posts">
      <div :class="{ 'container': !isMobile }">
        <table v-if="isPageLoaded">
          <tbody>
            <tr
              v-for="article in articles"
              :key="article.id">
              <td>
                <p class="title">
                  <nuxt-link
                    :to="'/articles/' + article.id"
                    class="link">{{ article.title }}</nuxt-link>
                </p>
                <p class="post-at">
                  <i class="far fa-clock"></i>
                  {{ formatted(article.post_at) }}
                </p>
              </td>
            </tr>
          </tbody>
        </table>
        <div
          v-else
          class="center">
          <LdsEllipsis class="fade-in" />
        </div>
      </div>
      <div
        v-if="articles.length !== 0"
        class="container pagination"
        >
        <button
          v-on:click="fethcArticles(prevPage)"
          :class="{enable: Boolean(prevPage)}"
          class="button prev">
          &lt;
        </button>
        <span class="page">{{ page }}</span>
        <button
          v-on:click="fethcArticles(nextPage)"
          :class="{ 'enable': Boolean(nextPage) }"
          class="button next"
          >
          &gt;
        </button>
      </div>
      <div
        v-else
        class="center"
        >
        <div class="center">
          <p class="hint">No more articles.</p>
        </div>
        <nuxt-link
          v-if="page > 1"
          class="link hint"
          to="/articles"
          >
          Go to first page
        </nuxt-link>
      </div>
    </div>

    <div
      v-if="!isMobile"
      class="nav">
      <!-- Tags -->
      <div class="tags container">
        <h3><i class="fas fa-hashtag"></i>Tags</h3>
        <p>
          <nuxt-link
            v-for="tag in tags"
            :key="tag.name"
            to=""
            class="link tag">
            {{ tag.name }}
          </nuxt-link>
        </p>
      </div>
    </div>

  </div>
</template>

<script>
import axios from 'axios'
import LdsEllipsis from '~/components/loadings/LdsEllipsis'
import { buildURLEncodedString, formattedUTCString } from '~/assets/js/utils'

const apiPath = '/v2/articles'

const articleFetcher = axiosClient => {
  return {
    fetch: uri => {
      return axiosClient.get(uri)
    },
  }
}

const buildLinks = (obj, path) => {
  let links = []
  if (obj.prev && typeof obj.prev === typeof '') {
    links.push({
      rel: 'prev', href: obj.prev.replace(apiPath, path),
    })
  }
  if (obj.next && typeof obj.next === typeof '') {
    links.push({
      rel: 'next', href: obj.next.replace(apiPath, path),
    })
  }
  return links
}

export default {
  components: {
    LdsEllipsis
  },
  head () {
    return {
      title: 'Articles',
      link:  this.links,
    }
  },
  asyncData({$axios, query, route}) {
    const q = buildURLEncodedString({
      page:    Boolean(query.page)    ? query.page    : 1,
      perPage: Boolean(query.perPage) ? query.perPage : 5,
    })
    const uri = `${apiPath}?${q}`
    return axios.all([
      $axios.get(uri),
      $axios.get(`/v1/articleTags`)
    ]).then(([articlesRes, tagsRes]) => {
      return {
        isMobile:     false,
        isPageLoaded: true,
        currentURI:   uri,
        articles:     articlesRes.data.articles,
        tags:         tagsRes.data,
        page:         articlesRes.data.page,
        perPage:      articlesRes.data.perPage,
        total:        articlesRes.data.total,
        prevPage:     articlesRes.data.prevPage,
        nextPage:     articlesRes.data.nextPage,
        links:        buildLinks({
          prev: articlesRes.data.prevPage,
          next: articlesRes.data.nextPage,
        }, route.path),
      }
    })
  },
  watchQuery: ['page', 'perPage'],
  mounted() {
    window.addEventListener('resize', this.calcIsMobile)
    this.calcIsMobile()
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.calcIsMobile)
  },
  methods: {
    formatted(dtString) {
      return formattedUTCString(dtString)
    },
    calcIsMobile() {
      this.isMobile = window.innerWidth <= 768
    },
    fethcArticles(uri) {
      if (!uri || uri === this.currentURI) {
        return
      }
      this.isPageLoaded = false
      articleFetcher(this.$axios)
        .fetch(uri)
        .then(res => {
          this.currentURI = uri
          this.articles   = res.data.articles
          this.page       = res.data.page
          this.prevPage   = res.data.prevPage
          this.nextPage   = res.data.nextPage
          this.links      = buildLinks({
            prev: res.data.prevPage,
            next: res.data.nextPage,
          }, this.$route.path)
          this.isPageLoaded = true
        })
    }
  }
}
</script>

<style lang="scss" scoped>
@import '~/assets/scss/styles.scss';
.container {
  @media screen and (min-width: $max_width_device + 1) {
    margin: 0 32px;
  }
  @media screen and (max-width: $max_width_device) {
    margin: 0 16px;
  }
  .posts {
    float: left;
    table {
      width: 100%;
      td {
        border-bottom: 1px solid rgba(0, 0, 0, 0.1);
      }
    }
    .title {
      @media screen and (min-width: $max_width_device + 1) {
        font-size: 1.8em;
      }
      @media screen and(max-width: $max_width_device) {
        font-size: 1.5em;
      }
    }
    .post-at {
      color: #777;
    }
    .pagination {
      align-items: center;
      display: flex;
      font-size: 1.1em;
      justify-content: center;
      text-align: center;
      .page {
        color: $color_text;
        cursor: default;
        height: 2em;
        line-height: 2em;
        width: 2em;
      }
      .button {
        background-color: transparent;
        border: none;
        border-radius: 50%;
        display:inline-block;
        font-size: 1.1em;
        opacity: 0.2;
        margin: 8px;
        &.prev {
          margin-right: 32px;
        }
        &.next {
          margin-left: 32px;
        }
        &.enable {
          opacity: 1;
          &:hover {
            color: $color_accent;
            cursor: pointer;
          }
        }
        &:focus {
          outline: 0;
          box-shadow: none;
        }
      }
    }
  }
  .nav {
    float: right;
    position: -webkit-sticky;
    position: -moz-sticky;
    position: -ms-sticky;
    position: -o-sticky;
    position: sticky !important;
    top: 44px !important; /* height of header */
    width: 33.3333%;
    .tags {
      .tag {
        display: inline-block;
        background-color: #777;
        padding: 1px 8px;
        margin: 2px;
        border-radius: 2px;
        font-weight: bold;
        font-size: 0.98em;
        color: white !important;
        &:hover {
          background-color: $color_accent;
          opacity: 0.8;
        }
      }
    }
  }
}
.hint {
  color: rgba(1, 1, 1, 0.1);
  font-size: 1.1em;
}
.fade-in {
  @include fade_in($opacity: 0.2, $duration: 2s);
}
.desktop {
  width: 66.6666%;
}
.mobile {
  width: 100%;
}
i {
  margin-right: 8px;
}
</style>
