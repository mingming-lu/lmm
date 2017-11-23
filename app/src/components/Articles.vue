<template>
  <div>
    <div v-for="(article, index) in articles" :key="article.title">
      <div class="lmm-card-4 lmm-margin lmm-white" style="width:66.6666%">
        <div class="lmm-container" style="text-align: left">
          <h2><b>{{ article.title }}</b></h2>
          <p class="lmm-opacity">Posted time: {{ article.posted_time }}</p>
          <p>{{ article.text }}</p>
          <p style="float: right"><button class="lmm-button lmm-padding-large lmm-white lmm-border"><b>READ MORE »</b></button></p>
        </div>
      </div>
    </div>
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