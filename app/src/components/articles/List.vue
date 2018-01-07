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
          <router-link to="" v-for="tag in tags" :key="tag.name" class="white link">
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
    axios.all([
      axios.get('http://api.lmm.local/articles/1'),
      axios.get('http://api.lmm.local/articles/1/categories'),
      axios.get('http://api.lmm.local/articles/1/tags')
    ]).then(axios.spread((articles, categories, tags) => {
      this.articles = articles.data
      this.categories = categories.data
      this.tags = tags.data
    })).catch((e) => {
      console.log(e.response.data)
    })
  }
}
</script>