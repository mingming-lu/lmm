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
      {{ category.name }} <input size="32" :id="category.name">
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
      axios.get(process.env.API_URL_BASE + '/v1/blog').then(res => {
        this.blogList = res.data.blog
      }).catch(e => {
        console.log(e.response.data)
      })
      this.newCategoryName = ''
    },
    fetchCategories () {
      this.newCategoryName = ''
      axios.get(process.env.API_URL_BASE + '/v1/categories').then(res => {
        this.categories = res.data.categories
      }).catch(e => {
        console.log(e.response.data)
      })
    },
    onSubmitCategory () {
      axios.post(process.env.API_URL_BASE + '/v1/categories', {
        name: this.newCategoryName
      }, {
        headers: {
          Authorization: localStorage.getItem('token')
        }
      }).then(res => {
        alert(res.data)
        this.fetchCategories()
      }).catch(e => {
        console.log(e)
        console.log(e.response.data)
      })
    },
    onUpdateCategory (category) {
      axios.put(process.env.API_URL_BASE + '/v1/categories/' + category.id, {
        name: document.getElementById(category.name).value
      }, {
        headers: {
          Authorization: localStorage.getItem('token')
        }
      }).then(res => {
        alert(res.data)
        this.fetchCategories()
      }).catch(e => {
        console.log(e)
        console.log(e.response.data)
      })
    },
    onDeleteCategory (category) {
      axios.delete(process.env.API_URL_BASE + '/v1/categories/' + category.id, {
        headers: {
          Authorization: localStorage.getItem('token')
        }
      }).then(res => {
        alert('deleted')
        this.fetchCategories()
      }).catch(e => {
        console.log(e)
        console.log(e.response.data)
      })
    }
  }
}
</script>
