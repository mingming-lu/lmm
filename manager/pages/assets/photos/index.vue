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
              :key="photo.url"
              xs4
            >
              <nuxt-link :to="`/assets/photos/${photo.name}`">
                <v-img
                  :src="photo.url"
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
    fetch: cursor => {
      return axios.get('/v1/photos?count=100')
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
      .fetch('TODO: paging')
      .then(res => {
        return {
          photos: res.data.items,
          cursor: res.data.next_cursor
        }
      })
  },
  methods: {
    pickOnePhoto() {
      this.$refs.photoPicker.click()
    },
    onPhotoPicked({ target }) {
      const photo = target.files[0]
      if (!photo) {
        return
      }

      let formData = new FormData()
      formData.append('photo', photo)
      this.$axios
        .post('/v1/photos', formData, {
          headers: {
            Authorization: `Bearer ${window.localStorage.getItem(
              'accessToken'
            )}`,
            'Content-Type': 'multipart/form-data'
          }
        })
        .then(res => {
          alert(`Uploaded\nmessage: ${res.data.message}`)
          fetcher(this.$axios)
            .fetch(1)
            .then(res => {
              this.photos = res.data.items
              this.cursor = res.data.next_cursor
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
