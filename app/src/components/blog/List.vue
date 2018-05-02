<template>
  <div class="container">
    <!-- Posts -->
    <div class="posts" :class="{ 'mobile-left': isMobile }">
      <div class="container">
        <div v-for="(blog, index) in blogList" :key="blog.id">
          <h3>
            <router-link :to="'/blog/' + blog.id" class="link">{{ blog.title }}</router-link>
          </h3>
          <p class="post-info"><i class="fa fa-fw fa-calendar-o"></i>{{ blog.created_at }}</p>
          <hr v-if="index !== blogList.length - 1" class="post-separator">
        </div>
      </div>
    </div>

    <div v-if="!isMobile" class="nav">
      <!-- Categories -->
      <div class="categories container">
        <h4><i class="fa fa-fw fa-folder"></i>Categories</h4>
        <router-link to="" v-for="category in categories" :key="category.id" class="link">
          <p class="category">{{ category.name }}</p>
        </router-link>
      </div>

      <!-- Tags -->
      <div class="tags container">
        <h4><i class="fa fa-fw fa-tags"></i>Tags</h4>
          <p>
            <router-link to="" v-for="tag in tags" :key="tag.name" class="link tag">
              {{ tag.name }}
            </router-link>
          </p>
      </div>
    </div>

    <router-view/>

  </div>
</template>

<script>
import axios from 'axios'
export default {
  data () {
    return {
      isMobile: false,
      blogList: [],
      categories: [],
      tags: []
    }
  },
  created () {
    this.fetchBlog()
    this.fetchCategories()
    this.fetchTags()
    this.calcIsMobile()
    window.addEventListener('resize', this.calcIsMobile)
  },
  beforeDestroy () {
    window.removeEventListener('resize', this.calcIsMobile)
  },
  methods: {
    fetchBlog () {
      axios.get('https://api.lmm.im/v1/users/1/blog').then(res => {
        this.blogList = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    fetchCategories () {
      axios.get('https://api.lmm.im/v1/users/1/categories').then(res => {
        this.categories = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    fetchTags () {
      axios.get('https://api.lmm.im/v1/users/1/tags').then(res => {
        this.tags = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    calcIsMobile () {
      this.isMobile = window.innerWidth <= 768
    }
  }
}
</script>

<style lang="scss" scoped>
@import '@/assets/scss/styles.scss';
.container {
  padding: 0 16px;
  .posts {
    float: left;
    width: 66.6666%;
    .post-info {
      opacity: 0.6;
    }
    .post-separator {
      opacity: 0.15;
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
        font-size: 0.88em;
        color: white !important;
        &:hover {
          background-color: $color_primary;
          opacity: 0.8;
        }
      }
    }
  }
}
.mobile-left {
  width: 100% !important;
}
i {
  margin-right: 8px;
}
</style>
