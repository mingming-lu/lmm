<template>
  <v-layout column>
    <v-list v-if="articles.length">
      <template v-for="article in articles">
        <v-list-tile
          :to="`/articles/edit?articleID=${article.id}`"
          :key="article.id"
          nuxt
          exact
        >
          <v-list-tile-content>
            <v-list-tile-title v-text="article.title" />
            <v-list-tile-sub-title v-text="new Date(article.post_at * 1e3).toLocaleString()" />
          </v-list-tile-content>
        </v-list-tile>
      </template>
    </v-list>
    <v-btn
      absolute
      color="accent"
      dark
      fab
      nuxt
      right
      to="/articles/new"
      top
    >
      <v-icon>create</v-icon>
    </v-btn>
  </v-layout>
</template>

<script>
export default {
  head () {
    return {
      title: 'Articles list'
    }
  },
  asyncData({$axios}) {
    return $axios.get('/v1/articles?count=100')
    .then(res => {
      return {
        articles:      res.data.articles,
        has_next_page: res.data.has_next_page
      }
    })
  }
}
</script>
