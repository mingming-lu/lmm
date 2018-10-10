<template>
  <v-layout column>
    <v-list>
      <template v-for="(article, index) in articles">
        <v-list-tile
          :to="`/articles/${article.id}`"
          :key="article.id"
          nuxt
          exact
        >
          <v-list-tile-content>
            <v-list-tile-title v-text="article.title" />
            <v-list-tile-sub-title v-text="new Date(article.post_at).toLocaleString()" />
          </v-list-tile-content>
        </v-list-tile>
        <v-divider v-if="index + 1 < articles.length" :key="`divider-${index}`"></v-divider>
      </template>
    </v-list>
  </v-layout>
</template>

<script>
export default {
  asyncData({$axios}) {
    return $axios.get('/v1/articles')
    .then(res => {
      console.log(res)
      return {
        articles:      res.data.articles,
        has_next_page: res.data.has_next_page
      }
    })
  }
}
</script>
