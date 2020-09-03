<template>
  <div>
    <v-btn
      absolute
      color="accent"
      dark
      fab
      right
      top
      @click="commit"
    >
      <v-icon>done</v-icon>
    </v-btn>
    <v-combobox
      v-model="tags"
      label="alternate text"
      chips
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
    <v-img
      :src="url"
      :alt="tags.join(' ')"
    />
    <v-snackbar
      v-model="committed"
      :timeout="2000"
      bottom
      color="success"
    >
      Committed
    </v-snackbar>
  </div>
</template>

<script>
const photoFetcher = httpClient => {
  return {
    fetch: id => {
      return httpClient.get(`/v1/photos/${id}`, {
        headers: {
          Authorization: `Bearer ${window.localStorage.getItem('accessToken')}`
        }
      })
    }
  }
}

export default {
  asyncData({ $axios, params }) {
    return photoFetcher($axios)
      .fetch(params.photoID)
      .then(res => {
        return {
          id: res.data.id,
          url: res.data.url,
          tags: res.data.tags
        }
      })
  },
  data() {
    return {
      committed: false
    }
  },
  methods: {
    removeTag(name) {
      this.tags.splice(this.tags.indexOf(name), 1)
      this.tags = [...this.tags]
    },
    commit() {
      this.$axios
        .put(
          `/v1/photos/${this.id}/tags`,
          {
            tags: this.tags.map(tag => {
              return tag
            })
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
          this.committed = true
        })
        .catch(err => {
          alert(err)
        })
    }
  }
}
</script>

<style scoped>
</style>
