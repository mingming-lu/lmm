<template>
  <div>

    <!-- articles list -->
    <h3>Articles List</h3>
    <router-link to="/posts/new" tag="button">post new article here</router-link>
    <div v-for="article in articles" :key="article.id"> 
      <router-link :to="'/posts/' + article.id + '/edit'">{{ article.title }}</router-link>
    </div>
    <hr class="lmm-opacity">

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
      axios.all([
        axios.get('http://api.lmm.local/articles/1'),
        axios.get('http://api.lmm.local/articles/1/categories')
      ]).then(axios.spread((articles, categories) => {
        this.articles = articles.data
        this.categories = categories.data
      })).catch((e) => {
        console.log(e)
      })
      this.newCategoryName = ''
    },
    onCreateCategory () {
      let name = this.newCategoryName.trim()
      if (!name) {
        return
      }
      if (!confirm('add category: `' + name + '`?')) {
        return
      }
      axios.post('http://api.lmm.local/articles/category', {
        user_id: 1,
        name: name
      }).then((res) => {
        this.fetchData()
      }).catch((e) => {
        console.log(e)
        alert(e.response.data)
      })
    },
    onUpdateCategory (category) {
      let name = document.getElementById(category.id).value.trim()
      if (!name) {
        return
      }
      if (name === category.name) {
        alert('no change')
        return
      }
      if (!confirm('change `' + category.name + '` to `' + name + '`?')) {
        return
      }

      // update category name by id
      axios.put('http://api.lmm.local/articles/category/' + category.id, {
        id: category.id,
        user_id: 1,
        name: name
      }).then((res) => {
        this.fetchData()
      }).catch((e) => {
        console.log(e)
        alert(e.response.data)
      })
    },
    onDeleteCategory (category) {
      if (!confirm('delete `' + category.name + '`?')) {
        return
      }
      axios.delete('http://api.lmm.local/articles/category/' + category.id).then((res) => {
        this.fetchData()
      }).catch((e) => {
        console.log(e)
        alert(e.response.data)
      })
    }
  }
}
</script>
