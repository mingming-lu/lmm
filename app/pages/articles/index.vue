<template>
  <div class="container">
    <!-- Articles -->
    <div
      :class="{ 'mobile-left': isMobile }"
      class="posts">
      <div :class="{container: !isMobile}">
        <no-ssr>
          <table v-if="isArticleListLoaded">
            <tr
              v-for="article in articles"
              :key="article.id">
              <td>
                <p class="post-title">
                  <nuxt-link
                    :to="'/articles/' + article.id"
                    class="link">{{ article.title }}</nuxt-link>
                </p>
                <p class="post-info">
                  <i class="fa fa-fw fa-calendar-o"/>
                  {{ formatted(article.post_at) }}
                </p>
              </td>
            </tr>
          </table>
          <div
            v-else
            class="center">
            <LdsEllipsis class="fade-in" />
          </div>
        </no-ssr>
      </div>
    </div>

    <div
      v-if="!isMobile"
      class="nav">
      <!-- Tags -->
      <div class="tags container">
        <h3><i class="fa fa-fw fa-tags"/>Tags</h3>
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
import { formattedUTCString } from '~/assets/js/utils'
export default {
  components: {
    LdsEllipsis
  },
  asyncData({$axios}) {
    return axios.all([
      $axios.get(`/v1/articles`),
      $axios.get(`/v1/articleTags`)
    ]).then(([articlesRes, tagsRes]) => {
      return {
        isArticleListLoaded: true,
        isMobile: false,
        articles: articlesRes.data.articles,
        tags:     tagsRes.data
      }
    })
  },
  created() {
    if (process.browser) {
      window.addEventListener('resize', this.calcIsMobile)
    }
  },
  mounted() {
    this.calcIsMobile()
  },
  beforeDestroy() {
    if (process.browser) {
      window.removeEventListener('resize', this.calcIsMobile)
    }
  },
  methods: {
    formatted(dtString) {
      return formattedUTCString(dtString)
    },
    calcIsMobile() {
      if (process.browser) {
        this.isMobile = window.innerWidth <= 768
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
    width: 66.6666%;
    table {
      width: 100%;
      td {
        border-bottom: 1px solid rgba(0, 0, 0, 0.1);
      }
    }
    .post-title {
      @media screen and (min-width: $max_width_device + 1) {
        font-size: 1.8em;
      }
      @media screen and(max-width: $max_width_device) {
        font-size: 1.5em;
      }
    }
    .post-info {
      color: #777;
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
    .categories {
      .category {
        font-size: 1.1em;
        i {
          opacity: 0;
        }
      }
    }
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
.fade-in {
  @include fade_in($opacity: 0.2, $duration: 2s);
}
.mobile-left {
  width: 100% !important;
}
i {
  margin-right: 8px;
}
</style>
