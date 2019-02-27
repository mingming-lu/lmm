<template>
  <v-layout>
    <input
      ref="photoPicker"
      accept="image/*"
      style="display: none"
      type="file"
      @change="onPhotoPicked"
    >
    <v-flex xs12>
      <v-card v-if="photos.length">
        <v-container 
          grid-list-sm 
          fluid>
          <v-layout 
            row 
            wrap 
            align-center>
            <v-flex
              v-for="photo in photos"
              :key="photo.name"
              xs4
            >
              <nuxt-link :to="`/assets/photos/${photo.name}`">
                <v-img
                  :src="wrapAssetURL(photo.name)"
                  class="img"
                />
              </nuxt-link>
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
      <v-icon>add_a_photo</v-icon>
    </v-btn>
  </v-layout>
</template>

<script>
const fetcher = axios => {
  return {
    fetch: page => {
      return axios.get('/v1/assets/photos?perPage=100')
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
    return fetcher($axios)
      .fetch(1)
      .then(res => {
        return {
          photos: res.data.photos,
          hasNextPage: res.data.hasNextPage
        }
      })
  },
  methods: {
    pickOnePhoto() {
      this.$refs.photoPicker.click()
    },
    wrapAssetURL(name) {
      return `${process.env.ASSET_URL}/photos/${name}`
    },
    onPhotoPicked({ target }) {
      const photo = target.files[0]
      if (!photo) {
        return
      }

      let formData = new FormData()
      formData.append('photo', photo)
      this.$axios
        .post('/v1/assets/photos', formData, {
          headers: {
            Authorization: `Bearer ${window.localStorage.getItem(
              'accessToken'
            )}`,
            'Content-Type': 'multipart/form-data'
          }
        })
        .then(res => {
          alert(`Uploaded\nmessage: ${res.data}`)
          fetcher(this.$axios)
            .fetch(1)
            .then(res => {
              this.photos = res.data.photos
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
