<template>
  <v-container grid-list-xl>
    <v-layout
      :row="row"
      :column="!row"
    >
      <v-flex xs6>
        <v-text-field
          v-model="articleTitle"
          label="title"
          required
        />
        <v-combobox
          v-model="articleTags"
          :items="tags"
          label="tags"
          chips
          clearable
          multiple
        >
          <template
            slot="selection"
            slot-scope="data"
          >
            <v-chip
              close
              @input="removeTag(data.item)"
            >
              <strong>{{ data.item }}</strong>&nbsp;
            </v-chip>
          </template>
        </v-combobox>
        <v-textarea
          v-model="articleBody"
          label="body"
          auto-grow
          required
        />
      </v-flex>
      <v-flex xs6>
        <v-subheader>Article Body Preview</v-subheader>
        <div
          v-hljs
          class="mx-3 preview"
        >
          {{ marked(articleBody) }}
        </div>
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
      ]).then(([article, tags]) => {
        return {
          articleID: article.data.id,
          articleTitle: article.data.title,
          articleBody: article.data.body,
          articleTags: article.data.tags.map(tag => {
            return tag.name
          }),
          tags: tags.data.map(tag => {
            return tag.name
          }),
          row: false
        }
      })
    }
  }
}

const marker = new Markdownit({
  html: true,
  typographer: true
})

export default {
  asyncData({ $axios, query }) {
    return fetcher($axios).fetch(query.id)
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
    updateArticle() {
      this.$axios
        .put(
          `/v1/articles/${this.articleID}`,
          {
            title: this.articleTitle,
            body: this.articleBody,
            tags: this.articleTags
          },
          {
            headers: {
              Authorization: `Bearer ${window.localStorage.getItem(
                'accessToken'
              )}`
            }
          }
        )
        .then(res => {
          alert('success')
          this.$router.replace({ path: `/articles` })
        })
        .catch(e => {
          alert(e)
        })
    },
    onResize() {
      this.row = window.innerWidth > 960
    }
  },
  head() {
    return {
      title: 'Edit an article'
    }
  },
  watchQuery: ['articleID']
}
</script>

<style scoped>
.preview {
  font-size: 16px; /* adjust to v-textarea */
}
.preview /deep/ pre code {
  font-family: Monaco, 'Courier', monospace;
}
</style>
