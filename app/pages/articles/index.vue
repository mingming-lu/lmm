<template>
  <div 
    v-if="isMounted" 
    class="container">
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
                  <i class="far fa-clock"/>
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
        <nuxt-link
          v-if="!isMobile"
          :to="prevPage"
          :class="{disabled: !Boolean(prevPage)}"
          class="navigation link prev"
        >
          &#10094; Prev
        </nuxt-link>
        <nuxt-link
          v-for="item in pager.items"
          :key="item"
          :to="item === page ? $route.fullPath : `${$route.path}?page=${item}&perPage=${perPage}`"
          :class="{active: item === page }"
          class="link page"
        >
          {{ item }}
        </nuxt-link>
        <nuxt-link
          v-if="!isMobile"
          :to="nextPage"
          :class="{disabled: !Boolean(nextPage) }"
          class="navigation link next"
        >
          Next &#10095;
        </nuxt-link>
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
        <h3><i class="fas fa-hashtag"/>Tags</h3>
        <p>
          <nuxt-link
            v-for="tag in tags"
            :class="{active: tag.name === $route.query.tag}"
            :key="tag.name"
            :to="buildLinkWithTagQuery(tag.name)"
            class="link tag"
          >
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
import {
  buildPageNumbers,
  buildURLEncodedString,
  formattedDateFromTimeStamp
} from '~/assets/js/utils'

const apiPath = '/v1/articles'

const articleFetcher = axiosClient => {
  return {
    fetch: uri => {
      return axiosClient.get(uri)
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
      title: 'Articles',
      link: this.links
    }
  },
  asyncData({ $axios, query, route }) {
    const q = buildURLEncodedString({
      flavor: true,
      page: Boolean(query.page) ? query.page : 1,
      perPage: Boolean(query.perPage) ? query.perPage : 5,
      tag: Boolean(query.tag) ? query.tag : undefined
    })
    const uri = `${apiPath}?${q}`
    return axios
      .all([$axios.get(uri), $axios.get(`/v1/articleTags`)])
      .then(([articlesRes, tagsRes]) => {
        const articles = articlesRes.data
        const page = articles.page
        const perPage = articles.perPage
        const total = articles.total
        return {
          isMobile: false,
          isPageLoaded: true,
          currentURI: uri,
          articles: articles.articles,
          tags: tagsRes.data,
          page: page,
          perPage: perPage,
          pager: {
            items: buildPageNumbers(page, Math.ceil(total / perPage), 5)
          },
          total: articles.total,
          prevPage: articles.prevPage
            ? articles.prevPage.replace(apiPath, route.path)
            : '',
          nextPage: articles.nextPage
            ? articles.nextPage.replace(apiPath, route.path)
            : '',
          links: buildLinks(
            {
              first: articles.firstPage,
              prev: articles.prevPage,
              next: articles.nextPage,
              last: articles.lastPage
            },
            route.path
          )
        }
      })
  },
  data() {
    return {
      isMounted: false
    }
  },
  watchQuery: ['page', 'perPage', 'tag'],
  mounted() {
    window.addEventListener('resize', this.calcIsMobile)
    this.calcIsMobile()
    this.isMounted = true
  },
  updated() {
    this.calcIsMobile()
    this.isMounted = true
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.calcIsMobile)
  },
  methods: {
    buildLinkWithTagQuery(tagName) {
      if (tagName === this.$route.query.tag) {
        return '/articles'
      }
      return `/articles?tag=${tagName}`
    },
    formatted(dtString) {
      return formattedDateFromTimeStamp(dtString)
    },
    calcIsMobile() {
      const newState = window.innerWidth <= 768
      if (this.isMobile !== newState) {
        this.isMobile = newState
      }
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
        font-size: 1.5em;
      }
      @media screen and(max-width: $max_width_device) {
        font-size: 1.2em;
      }
    }
    .post-at {
      color: #777;
    }
    .pagination {
      align-items: center;
      display: flex;
      font-size: 1.15em;
      justify-content: center;
      text-align: center;
      padding: 8px;
      .navigation {
        &.disabled {
          opacity: 0.2;
          pointer-events: none !important;
        }
        &.prev {
          margin-right: 16px;
        }
        &.next {
          margin-left: 16px;
        }
      }
      .page {
        color: $color_text;
        height: 2em;
        line-height: 2em;
        width: 2em;
        margin: 4px;
        &.active {
          background-color: white !important;
          color: $color_accent !important;
          pointer-events: none !important;
        }
        &:hover {
          background-color: $color_accent;
          border-radius: 2px;
          color: white;
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
        &.active {
          background-color: $color_accent;
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
