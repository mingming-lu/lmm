<template>
  <div>
    <router-link to="/posts/new">post new article here</router-link>
    <hr class="lmm-level-opacity">
    <h3>Articles List</h3>
    <div v-for="article in articles" :key="article.id"> 
      <router-link :to="'/posts/' + article.id + '/edit'">{{ article.title }}</router-link>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  data () {
    return {
      articles: []
    }
  },
  created () {
    axios.get('http://api.lmm.local/articles?user_id=1').then((res) => {
      this.articles = res.data
    }).catch((e) => {
      console.log(e)
    })
  }
}
</script>
