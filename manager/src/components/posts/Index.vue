<template>
  <div>

    <!-- articles list -->
    <h3>Articles List</h3>
    <router-link to="/posts/new" tag="button">post new article here</router-link>
    <div v-for="article in articles" :key="article.id"> 
      <router-link :to="'/posts/' + article.id + '/edit'">{{ article.title }}</router-link>
    </div>
    <hr class="opacity">

    <!-- categories list -->
    <h3>Categories List</h3>
    <form>
      <input size="32" v-model="newCategoryName" placeholder="input new category here">
      <input type="submit" value="Add" @click.prevent="onCreateCategory()">
    </form>
    <form v-for="category in categories" :key="category.name">
      {{ category.name }} <input size="32" :id="category.id">
      <input type="submit" value="Update" @click.prevent="onUpdateCategory(category)">
      <input type="submit" value="Delete" @click.prevent="onDeleteCategory(category)">
    </form>

  </div>
</template>

<script>
import axios from 'axios'

export default {
  data () {
    return {
      articles: [],
      categories: [],
      newCategoryName: ''
    }
  },
  created () {
    this.fetchData()
  },
  methods: {
    fetchData () {
      this.fetchBlogs()
      this.fetchCategories()
    },
    fetchBlogs () {
      axios.get('http://api.lmm.local/articles?user=1').then(res => {
        this.articles = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
      this.newCategoryName = ''
    },
    fetchCategories () {
      axios.get('http://api.lmm.local/categories?user=1').then(res => {
        this.categories = res.data
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    onCreateCategory () {
    },
    onUpdateCategory (category) {
    },
    onDeleteCategory (category) {
    }
  }
}
</script>
