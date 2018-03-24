<template>
  <div class="container">
    <!-- Posts -->
    <div class="left" :class="{ 'mobile-left': isMobile }">
      <div class="container">
        <div v-for="(blog, index) in blogList" :key="blog.id">
          <h3>
            <router-link :to="'/blog/' + blog.id" class="link white">{{ blog.title }}</router-link>
          </h3>
          <p class="opacity"><i class="fa fa-fw fa-calendar-o"></i>{{ blog.created_at }}</p>
          <hr v-if="index !== blogList.length - 1" class="opacity-plus">
        </div>
      </div>
    </div>

    <div v-if="!isMobile" class="right nav">
      <!-- Categories -->
      <div class="container">
        <h4><i class="fa fa-fw fa-folder"></i>Categories</h4>
        <router-link to="" v-for="category in categories" :key="category.id" class="white link">
          <p>{{ category.name }}</p>
        </router-link>
      </div>

      <!-- Tags -->
      <div id="tags" class="container">
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
<style scoped>
.container .left {
  width: 66.666%;
}
.container .right {
  width: 33.333%;
}
.mobile-left {
  width: 100% !important;
}
#tags {
  margin-top: 16px;
}
</style>
