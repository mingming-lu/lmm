<template>
  <div class="lmm-row">
    <!-- Posts -->
    <div class="lmm-left" style="width:66.6666%; display:inline-block">
      <div v-for="article in articles" :key="article.title">
        <div class="lmm-card-4 lmm-container lmm-margin" style="text-align:left">
          <h2>{{ article.title }}</h2>
          <p class="lmm-opacity">{{ article.created_date }}</p>
          <p>{{ article.text }}</p>
          <p class="lmm-right">
            <router-link :to="'/articles/' + article.id">
              <button class="lmm-button lmm-padding-large lmm-white lmm-border">
                <b>Read More >></b>
              </button>
            </router-link>
          </p>
        </div>
      </div>
    </div>

    <!-- Categories -->
    <div class="lmm-right" style="width:33.3333%; display:inline-block; text-align:left">
      <div class="lmm-container lmm-margin lmm-card-4">
        <p><b>Categories</b></p>
        <hr>
        <div v-for="category in categories" :key="category.id">
          <p>
            <router-link to="" class="lmm-white lmm-link lmm-hover-light-grey">{{ category.name }}</router-link>
          </p>
        </div>
      </div>

    <!-- Tags -->
      <div class="lmm-container lmm-margin lmm-card-4">
        <p><b>Tags</b></p>
        <hr>
        <p>
          <span v-for="tag in tags" :key=tag.id class="lmm-tag">
            <router-link to="" class="lmm-white lmm-hover-light-grey lmm-link">{{ tag.name }}</router-link>
            <br>
          </span>
        </p>
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
      console.log(e)
    })
  }
}
</script>