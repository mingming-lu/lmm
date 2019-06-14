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
      v-model="alts"
      label="alternate text"
      chips
      multiple
    >
      <template 
        slot="selection" 
        slot-scope="data">
        <v-chip
          close
          @input="removeAlt(data.item)"
        >
          <strong>{{ data.item }}</strong>&nbsp;
        </v-chip>
      </template>
    </v-combobox>
    <v-img
      :src="url"
      :alt="alts.join(' ')"
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
    fetch: name => {
      return httpClient.get(`/v1/photos/${name}`, {
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
      .fetch(params.name)
      .then(res => {
        return {
          name: name,
          url: res.data.url,
          alts: res.data.tags
        }
      })
  },
  data() {
    return {
      committed: false
    }
  },
  methods: {
    removeAlt(name) {
      this.alts.splice(this.alts.indexOf(name), 1)
      this.alts = [...this.alts]
    },
    commit() {
      this.$axios
        .put(
          `/v1/assets/photos/${this.name}/alts`,
          {
            names: this.alts.map(alt => {
              return alt
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
