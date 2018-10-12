<template>
  <v-layout column>
    <v-list>
      <template v-for="article in articles">
        <v-list-tile
          :to="`/articles/edit?articleID=${article.id}`"
          :key="article.id"
          nuxt
          exact
        >
          <v-list-tile-content>
            <v-list-tile-title v-text="article.title" />
            <v-list-tile-sub-title v-text="new Date(article.post_at).toLocaleString()" />
          </v-list-tile-content>
        </v-list-tile>
      </template>
    </v-list>
    <v-btn
      absolute
      color="accent"
      dark
      fab
      right
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
      title: 'Articles'
    }
  },
  asyncData({$axios}) {
    return $axios.get('/v1/articles')
    .then(res => {
      return {
        articles:      res.data.articles,
        has_next_page: res.data.has_next_page
      }
    })
  }
}
</script>
