<template>
  <div>
    <ol>
    <li v-for="(article, index) in articles" :key="article.title">
      {{ article.title }}: {{ article.text }}
    </li>
  </ol>
  </div>
</template>

<script>
import * as request from '@/request'
export default {
  data () {
    request.get('http://localhost:8081/articles', (response) => {
      response.result.articles.forEach((article) => {
        this.articles.push(article)
      })
    })
    return {
      articles: []
    }
  }
}
</script>