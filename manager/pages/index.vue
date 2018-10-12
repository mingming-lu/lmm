<template>
  <v-layout column>
    <v-list>
      <template v-for="item in indices">
        <v-list-tile
          :to="item.to"
          :key="item.to"
          nuxt
          exact
        >
          <v-list-tile-content>
            <v-list-tile-title v-text="item.name" />
            <v-list-tile-sub-title v-text="item.description" />
          </v-list-tile-content>
        </v-list-tile>
      </template>
    </v-list>
  </v-layout>
</template>

<script>
export default {
  head () {
    return {
      title: 'Index'
    }
  },
  data () {
    return {
      indices: [
        {
          to: '/articles',
          name: 'Articles',
          description: 'Post new articles or edit existing articles'
        },
        {
          to: '/assets/images',
          name: 'Assets (Images)',
          description: 'Upload images or view uploaded images'
        },
        {
          to: '/assets/photos',
          name: 'Assets (Photos)',
          description: 'Upload photos or view uploaded photos'
        }
      ]
    }
  }
}
</script>
