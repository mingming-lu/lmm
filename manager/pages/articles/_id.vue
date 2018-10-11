<template>
  <v-container grid-list-xl>
    <v-layout :row="row" :column="!row">
      <v-flex xs6>
        <v-text-field label="title" required v-model="articleTitle"/>
        <v-combobox
          v-model="articleTags"
          :items="tags"
          label="tags"
          chips
          clearable
          multiple
        >
          <template slot="selection" slot-scope="data">
            <v-chip
              close
              @input="removeTag(data.item)"
            >
              <strong>{{ data.item }}</strong>&nbsp;
            </v-chip>
          </template>
        </v-combobox>
        <v-textarea label="body" required v-model="articleBody"/>
      </v-flex>
      <v-flex xs6>
        <v-subheader>Article Body Preview</v-subheader>
        <v-textarea class="mx-3" v-html="marked(articleBody)"/>
      </v-flex>
    </v-layout>
    <v-btn
      absolute
      color="accent"
      dark
      fab
      right
      top
    >
      <v-icon>autorenew</v-icon>
    </v-btn>
  </v-container>
</template>

<script>
import Markdownit from 'markdown-it'
export default {
  validate({params}) {
    return /^[\d\w]{8}$/.test(params.id)
  },
  asyncData({$axios, params}) {
    return Promise.all([
      $axios.get(`/v1/articles/${params.id}`),
      $axios.get(`/v1/articleTags`)
    ]).then(([article, tags]) => {
      return {
        articleTitle: article.data.title,
        articleBody:  article.data.body,
        articleTags:  article.data.tags.map(tag => { return tag.name }),
        row:          false,
        tags:         tags.data.map(tag => { return tag.name })
      }
    })
  },
  mounted() {
    this.onResize()
    window.addEventListener('resize', this.onResize, { passive: true })
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.onResize, { passive: true })
  },
  methods: {
    marked(text) {
      return new Markdownit({
        html: true,
        typographer: true
      }).render(text)
    },
    removeTag(item) {
      this.articleTags.splice(this.articleTags.indexOf(item), 1)
      this.articleTags = [...this.articleTags]
    },
    onResize() {
      this.row = window.innerWidth > 960
    }
  }
}
</script>
