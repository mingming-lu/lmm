<template>
  <div class="container">
    <!-- Posts -->
    <div class="left" style="width:75%;">
      <div v-for="(blog, index) in blog" :key="blog.id">
        <div class="container">
          <h2>
            <router-link :to="'/blog/' + blog.id" class="link white">{{ blog.title }}</router-link>
          </h2>
          <p class="opacity"><i class="fa fa-fw fa-calendar"></i>{{ blog.created_at }}</p>
        </div>
        <hr v-if="index !== blog.length - 1" class="opacity-plus">
      </div>
    </div>

    <div class="right nav" style="width:25%;">
      <!-- Categories -->
      <div class="container">
        <h4><i class="fa fa-fw fa-folder-o"></i>Categories</h4>
        <router-link to="" v-for="category in categories" :key="category.id" class="white link">
          <p>{{ category.name }}</p>
        </router-link>
      </div>

      <!-- Tags -->
      <div class="container">
        <h4><i class="fa fa-fw fa-tags"></i>Tags</h4>
          <router-link to="" v-for="tag in tags" :key="tag.name" class="link tag">
            {{ tag.name }}
          </router-link>
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
      blog: [],
      categories: [],
      tags: []
    }
  },
  created () {
    this.fetchBlog()
    this.fetchCategories()
    this.fetchTags()
  },
  methods: {
    fetchBlog () {
      axios.get('http://api.lmm.im/blog?user=1').then(res => {
        this.blog = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    fetchCategories () {
      axios.get('http://api.lmm.im/categories?user=1').then(res => {
        this.categories = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    fetchTags () {
      axios.get('http://api.lmm.im/tags?user=1').then(res => {
        this.tags = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
    }
  }
}
</script>