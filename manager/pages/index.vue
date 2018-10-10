<template>
  <v-layout column>
    <v-list>
      <v-subheader class="headline">Index</v-subheader>
      <template v-for="(item, index) in indices">
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
        <v-divider v-if="index + 1 < indices.length" :key="`divider-${index}`"></v-divider>
      </template>
    </v-list>
  </v-layout>
</template>

<script>
export default {
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
