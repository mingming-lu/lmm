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
        <div class="mx-3 preview" v-hljs v-html="marked(articleBody)"></div>
      </v-flex>
    </v-layout>
    <v-btn
      absolute
      color="accent"
      dark
      fab
      right
      top
      @click="updateArticle"
    >
      <v-icon>autorenew</v-icon>
    </v-btn>
    <v-snackbar
      v-model="updatedSnackbar"
      bottom
      color="success"
      :timeout="3000"
    >
      Updated
    </v-snackbar>
  </v-container>
</template>

<script>
import Markdownit from 'markdown-it'

const fetcher = axiosClient => {
  return {
    fetch: articleID => {
      return Promise.all([
        axiosClient.get(`/v1/articles/${articleID}`),
        axiosClient.get(`/v1/articleTags`)
      ])
      .then(([article, tags]) => {
        return {
          articleID:       article.data.id,
          articleTitle:    article.data.title,
          articleBody:     article.data.body,
          articleTags:     article.data.tags.map(tag => { return tag.name }),
          tags:            tags.data.map(tag => { return tag.name }),
          row:             false,
          updatedSnackbar: false,
        }
      })
    }
  }
}

const marker = new Markdownit({
  html:        true,
  typographer: true
})

export default {
  head() {
    return {
      title: 'Edit an article',
    }
  },
  validate({query}) {
    return /^[\d\w]{8}$/.test(query.articleID)
  },
  asyncData({$axios, query}) {
    return fetcher($axios).fetch(query.articleID)
  },
  mounted() {
    this.onResize()
    window.addEventListener('resize', this.onResize, { passive: true })
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.onResize, { passive: true })
  },
  watchQuery: ['articleID'],
  methods: {
    marked(text) {
      return marker.render(text)
    },
    removeTag(item) {
      this.articleTags.splice(this.articleTags.indexOf(item), 1)
      this.articleTags = [...this.articleTags]
    },
    updateArticle() {
      this.$axios
        .put(`/v1/articles/${this.articleID}`, {
          title: this.articleTitle,
          body:  this.articleBody,
          tags:  this.articleTags,
        }, {
          headers: {
            Authorization: `Bearer ${window.localStorage.getItem('accessToken')}`,
          },
        })
        .then(res => {
          this.updatedSnackbar = true
          fetcher(this.$axios).fetch(this.articleID).then(data => {
            this.articleID    = data.articleID
            this.articleTitle = data.articleTitle
            this.articleBody  = data.articleBody
            this.articleTags  = data.articleTags
            this.tags         = data.tags
          })
        })
        .catch(e => {
          alert(e)
        })
    },
    onResize() {
      this.row = window.innerWidth > 960
    }
  }
}
</script>

<style scoped>
.preview {
  font-size: 16px; /* adjust to v-textarea */
}
.preview /deep/ pre code {
  font-family: Monaco, "Courier", monospace;
}
</style>
