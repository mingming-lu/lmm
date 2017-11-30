<template>
  <div class="lmm-row">
    <!-- Posts -->
    <div class="lmm-left" style="width:66.6666%; display:inline-block">
      <div v-for="(article, index) in articles" :key="article.title">
        <div class="lmm-card-4 lmm-container lmm-margin" style="text-align:left">
          <h2><b>{{ article.title }}</b></h2>
          <p>{{ format(article.posted_time) }}</p>
          <p>{{ article.text }}</p>
          <router-link :to="'/articles/' + article.id" class="lmm-right lmm-button lmm-border" tag="button">READ MORE >></router-link>
        </div>
      </div>
    </div>

    <!-- Category -->
    <div class="lmm-right" style="width:33.3333%; display:inline-block; text-align:left">
      <div class="lmm-card-4 lmm-container lmm-margin">
        <h4>Category</h4>
        <hr>
        <p><router-link to="" class="lmm-white lmm-hover-light-grey" style="text-decoration:none">哈三篇</router-link></p>
        <p><router-link to="" class="lmm-white lmm-hover-light-grey" style="text-decoration:none">番外篇</router-link></p>
        <p><router-link to="" class="lmm-white lmm-hover-light-grey" style="text-decoration:none">哈轶事</router-link></p>
      </div>
    </div>

    <!-- Tags -->
    <div class="lmm-right" style="width:33.3333%; display:inline-block; text-align:left">
      <div class="lmm-card-4 lmm-container lmm-margin">
        <h4>Tags</h4>
        <hr>
        <p>
          <span class="lmm-tag"><router-link to="" class="lmm-white lmm-hover-light-grey lmm-link">上海交通大学</router-link></span>
          <span class="lmm-tag"><router-link to="" class="lmm-white lmm-hover-light-grey lmm-link">夏威夷吉他</router-link></span>
          <span class="lmm-tag"><router-link to="" class="lmm-white lmm-hover-light-grey lmm-link">张宝华</router-link></span>
          <span class="lmm-tag"><router-link to="" class="lmm-white lmm-hover-light-grey lmm-link">Excited</router-link></span>
          <span class="lmm-tag"><router-link to="" class="lmm-white lmm-hover-light-grey lmm-link">高腰裤</router-link></span>
        </p>
      </div>
    </div>
    <router-view/>

  </div>
</template>

<script>
import * as request from '@/request'
import * as utils from '@/utils'
export default {
  data () {
    request.get('http://localhost:8081/articles', (response) => {
      response.articles.forEach((article) => {
        article.id = 1
        this.articles.push(article)
      })
    })
    return {
      articles: [],
      format: utils.formattedTime
    }
  }
}
</script>