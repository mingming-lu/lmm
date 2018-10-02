<template>
  <div class="container">
    <form enctype="multipart/form-data">
      <div>
        <input type="file" multiple accept="image/*" @change="select($event.target.files)">
      </div>
    </form>

    <p class="error" v-if="errMsg !== ''">{{ errMsg }}</p>

    <div class="container">
      <img v-for="image in images" :src="image" :key="image.name">
    </div>

    <button @click="upload()" :disabled="files.length === 0">Upload</button>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  data () {
    return {
      src: '',
      onSaving: false,
      errMsg: '',
      files: [],
      images: []
    }
  },
  methods: {
    select (files) {
      Array.from(files).forEach(file => {
        this.files.push(file)
        let reader = new FileReader()
        reader.onloadend = () => {
          this.images.push(reader.result)
        }
        reader.readAsDataURL(file)
      })
    },
    upload () {
      this.onSaving = true
      let formData = new FormData()
      this.files.forEach(file => {
        formData.append('image', file)
      })
      axios.post(process.env.API_URL_BASE + '/v1/images', formData, {
        headers: {
          Authorization: 'Bearer ' + localStorage.getItem('token')
        }
      }).then(res => {
        alert(res.data)
        this.files = []
        this.images = []
        this.onSaving = false
      }).catch(e => {
        this.errMsg = e.response.data
        this.files = []
        this.images = []
        this.onSaving = false
      })
    }
  }
}
</script>
<style scoped>
img {
  width: 100%;
}
</style>
