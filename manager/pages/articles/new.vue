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
        <v-textarea label="body" auto-grow required v-model="articleBody"/>
      </v-flex>
      <v-flex xs6>
        <v-subheader>Article Body Preview</v-subheader>
        <div class="mx-3 preview" v-html="marked(articleBody)"></div>
      </v-flex>
    </v-layout>
    <v-btn
      absolute
      color="accent"
      dark
      fab
      right
      top
      @click="postArticle"
    >
      <v-icon>save</v-icon>
    </v-btn>
  </v-container>
</template>

<script>
import Markdownit from 'markdown-it'

const marker = new Markdownit({
  html:        true,
  typographer: true
})

export default {
  head() {
    return {
      title: 'Post an article',
    }
  },
  asyncData({ $axios }) {
    return $axios
      .get('/v1/articleTags')
      .then(res => {
        return {
          articleTitle:   '',
          articleBody:    '',
          articleTags:    [],
          tags:           res.data.map(tag => tag.name),
          row:            false,
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
      return marker.render(text)
    },
    removeTag(item) {
      this.articleTags.splice(this.articleTags.indexOf(item), 1)
      this.articleTags = [...this.articleTags]
    },
    onResize() {
      this.row = window.innerWidth > 960
    },
    postArticle() {
      this.$axios
        .post('/v1/articles', {
          title: this.articleTitle,
          body:  this.articleBody,
          tags:  this.articleTags,
        }, {
          headers: {
            Authorization: `Bearer ${window.localStorage.getItem('accessToken')}`,
          }
        })
        .then(res => {
          alert(`Article posted\nmessage: ${JSON.stringify(res.data)}`)
          this.$router.push('/articles')
        })
    },
  }
}
</script>

<style scoped>
.preview {
  font-size: 16px; /* adjust to v-textarea */
}
.preview /deep/ code {
  font-family: Monaco, "Courier", monospace;
  width: 100%;
}
</style>
