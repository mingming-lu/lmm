<template>
  <div class="container">
    <router-link to="/image/upload" class="title">Upload</router-link>
    <div align="left">Tip: 开关按钮用来控制是否显示在画廊</div>
    <div v-for="image in normals" :key="image.name">
      <div class="gallery">
        <img :src="imageURL(image.name)" align="top">
        <label class="switch"> <input type="checkbox" :ref="image.name" :checked="true === image.isPhoto" @click="togglePhotoSwitch(image.name)">
          <span class="slider"></span>
        </label>
      </div>
    </div>

    <div v-for="image in photos" :key="image.name">
      <div class="gallery">
        <img :src="imageURL(image.name)" align="top">
        <label class="switch"> <input type="checkbox" :ref="image.name" :checked="true === image.isPhoto" @click="togglePhotoSwitch(image.name)">
          <span class="slider"></span>
        </label>
      </div>
    </div>
  </div>
</template>
<script>
import axios from 'axios'

export default {
  data () {
    return {
      normals: [],
      photos: []
    }
  },
  created () {
    this.fetchNormalImages()
    this.fetchPhotos()
  },
  methods: {
    fetchNormalImages () {
      axios.get(process.env.API_URL_BASE + '/v1/images?type=normal').then(res => {
        let normals = res.data.images
        normals.forEach((image, index) => {
          normals[index].isPhoto = false
        })
        this.normals = normals
      }).catch(err => {
        console.log(err)
      })
    },
    fetchPhotos () {
      axios.get(process.env.API_URL_BASE + '/v1/images?type=photo').then(res => {
        let photos = res.data.images
        photos.forEach((image, index) => {
          photos[index].isPhoto = true
        })
        this.photos = photos
      }).catch(err => {
        console.log(err)
      })
    },
    imageURL: (name) => {
      return process.env.IMAGE_URL_BASE + '/' + name
    },
    togglePhotoSwitch (name) {
      let sw = this.$refs[name][0]
      sw.disabled = true
      if (sw.checked === true) {
        this.turnOnPhotoSwitch(name)
      } else {
        this.turnOffPhotoSwitch(name)
      }
    },
    turnOnPhotoSwitch (name) {
      var sw = this.$refs[name][0]
      axios.put(process.env.API_URL_BASE + '/v1/images/' + name + '?type=photo', {
      }, {
        headers: {
          Authorization: 'Bearer ' + localStorage.getItem('token')
        }
      }).then(res => {
        console.log(res.data)
        sw.disabled = false
      }).catch(e => {
        console.log(e.response.data)
        sw.disabled = false
      })
    },
    turnOffPhotoSwitch (name) {
      var sw = this.$refs[name][0]
      axios.put(process.env.API_URL_BASE + '/v1/images/' + name + '?type=normal', {
      }, {
        headers: {
          Authorization: 'Bearer ' + localStorage.getItem('token')
        }
      }).then(res => {
        console.log(res.data)
        sw.disabled = false
      }).catch(e => {
        console.log(e.response.data)
        sw.disabled = false
      })
    }
  }
}
</script>
<style scoped>
.title {
  display: block;
}
div .gallery {
  float: left;
  width: 240px;
  margin: 8px;
  border: 1px solid #ccc;
}
div .gallery img {
  width: 100%;
  height: auto;
}
.switch {
  position: relative;
  display: inline-block;
  width: 60px;
  height: 34px;
}

.switch input {display:none;}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
}

.slider:before {
  position: absolute;
  content: "";
  height: 26px;
  width: 26px;
  left: 4px;
  bottom: 4px;
  background-color: white;
}

input:checked + .slider {
  background-color: dodgerblue;
}

input:checked + .slider:before {
  -webkit-transform: translateX(26px);
  -ms-transform: translateX(26px);
  transform: translateX(26px);
}

input:disabled + .slider {
  opacity: 0.2;
}

.slider {
  border-radius: 34px;
}

.slider:before {
  border-radius: 50%;
}
</style>
