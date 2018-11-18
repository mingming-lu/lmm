<template>
  <v-layout>
    <input
      accept="image/*"
      ref="imagePicker"
      style="display: none"
      type="file"
      @change="onPhotoPicked"
    >
    <v-flex xs12>
      <v-card v-if="images.length">
        <v-container grid-list-sm fluid>
          <v-layout row wrap>
            <v-flex
              v-for="image in images"
              :key="image.name"
              xs4
            >
              <v-img
                class="img"
                @click="copyURLToClipboard(image.name)"
                :src="wrapAssetURL(image.name)"
              >
                <v-layout
                  slot="placeholder"
                >
                  <v-progress-circular
                    indeterminate
                    color="grey lighten-5"
                  />
                </v-layout>
              </v-img>
            </v-flex>
          </v-layout>
        </v-container>
      </v-card>
    </v-flex>
    <v-btn
      absolute
      color="accent"
      dark
      fab
      right
      top
      @click="pickOnePhoto"
    >
      <v-icon>add_photo_alternate</v-icon>
    </v-btn>
    <v-snackbar
      v-model="copied"
      :timeout="2000"
      bottom
      color="success"
    >
      The URL has been copied to clipboard.
    </v-snackbar>
  </v-layout>
</template>

<script>

const fetcher= axios => {
  return {
    fetch: page => {
      return axios.get('/v1/assets/images?perPage=100')
    }
  }
}

export default {
  head() {
    return {
      title: 'Photos'
    }
  },
  asyncData({ $axios }) {
    return fetcher($axios).fetch(1)
      .then(res => {
        return {
          images:      res.data.images,
          hasNextPage: res.data.hasNextPage
        }
      })
  },
  data() {
    return {
      copied: false,
    }
  },
  methods: {
    fallbackCopyURLToClipboard(name) {
      const url = this.wrapAssetURL(name)
      const textArea = document.createElement("textarea");
      textArea.value = url
      textArea.style = 'display: none'
      document.body.appendChild(textArea)
      textArea.select()
      if (document.execCommand('copy') === true) {
        this.copied = true
      }
      document.body.removeChild(textArea);
    },
    copyURLToClipboard(name) {
      if (!navigator.clipboard) {
        return this.fallbackCopyURLToClipboard(name)
      }
      const url = this.wrapAssetURL(name)
      navigator.clipboard.writeText(url).then(() => {
        this.copied = true
      }, err => {
        console.log(err)
      })
    },
    pickOnePhoto() {
      this.$refs.imagePicker.click()
    },
    wrapAssetURL(name) {
      return `${process.env.ASSET_URL}/images/${name}`
    },
    onPhotoPicked({ target }) {
      const image = target.files[0]
      if (!image) {
        return
      }

      let formData = new FormData();
      formData.append("image", image);
      this.$axios
        .post('/v1/assets/images', formData, {
          headers: {
            'Authorization': `Bearer ${window.localStorage.getItem('accessToken')}`,
            'Content-Type':  'multipart/form-data',
          }
        })
        .then(res => {
          alert(`Uploaded\nmessage: ${res.data}`)
          fetcher(this.$axios).fetch(1)
            .then(res => {
              this.images      = res.data.images
              this.hasNextPage = res.data.hasNextPage
            })
        })
    }
  }
}
</script>

<style scoped>
.img:hover {
  cursor: pointer;
}
</style>
