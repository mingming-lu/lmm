<template>
  <div class="container">
    <!-- Posts -->
    <div class="left" style="width:75%;">
      <div v-for="(article, index) in articles" :key="article.id">
        <div class="container">
          <h2>
            <router-link :to="'/articles/' + article.id" class="link white">{{ article.title }}</router-link>
          </h2>
          <p class="text-right opacity">{{ article.created_date }}</p>
        </div>
        <hr v-if="index !== articles.length - 1" class="opacity-plus">
      </div>
    </div>

    <div class="right nav" style="width:25%;">
      <!-- Categories -->
      <div class="container">
        <h4>Categories</h4>
        <router-link to="" v-for="category in categories" :key="category.id" class="white link">
          <p>{{ category.name }}</p>
        </router-link>
      </div>

      <!-- Tags -->
      <div class="container">
        <h4>Tags</h4>
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
      articles: [],
      categories: [],
      tags: []
    }
  },
  created () {
    this.fetchArticles()
    this.fetchCategories()
    this.fetchTags()
  },
  methods: {
    fetchArticles () {
      axios.get('http://api.lmm.im/users/1/articles').then(res => {
        this.articles = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    fetchCategories () {
      axios.get('http://api.lmm.im/users/1/categories').then(res => {
        this.categories = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    fetchTags () {
      axios.get('http://api.lmm.im/users/1/tags').then(res => {
        this.tags = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
    }
  }
}
</script>