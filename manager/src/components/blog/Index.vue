<template>
  <div>

    <!-- blog list -->
    <h3>Blog List</h3>
    <router-link to="/blog/new" tag="button">post new blog here</router-link>
    <div v-for="blog in blogList" :key="blog.id"> 
      <router-link :to="'/blog/' + blog.id + '/edit'">{{ blog.title }}</router-link>
    </div>
    <hr class="opacity">

    <!-- categories list -->
    <h3>Categories List</h3>
    <form>
      <input size="32" v-model="newCategoryName" placeholder="input new category here">
      <input type="submit" value="Add" @click.prevent="onSubmitCategory()">
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
      blogList: [],
      categories: [],
      newCategoryName: ''
    }
  },
  created () {
    this.fetchData()
  },
  methods: {
    fetchData () {
      this.fetchBlogList()
      this.fetchCategories()
    },
    fetchBlogList () {
      axios.get('https://api.lmm.im/users/1/blogs').then(res => {
        this.blogList = res.data
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
    onSubmitCategory () {
    },
    onUpdateCategory (category) {
    },
    onDeleteCategory (category) {
    }
  }
}
</script>
