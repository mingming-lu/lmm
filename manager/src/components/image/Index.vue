<template>
  <div class="container">
    <router-link to="/image/upload" class="title">Upload</router-link>
    <div align="left">Tip: 开关按钮用来控制是否显示在画廊</div>
    <div v-for="image in images" :key="image.name">
      <div class="gallery">
        <img :src="imageURL(image.name)" align="top">
        <label class="switch">
          <input type="checkbox" :ref="image.name" :checked="true === image.isPhoto" @click="togglePhotoSwitch(image.name)">
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
      images: []
    }
  },
  created () {
    this.fetchImages()
  },
  methods: {
    fetchImages () {
      axios.all([
        axios.get(process.env.API_URL_BASE + '/v1/users/1/images'),
        axios.get(process.env.API_URL_BASE + '/v1/users/1/images/photos')
      ]).then(axios.spread((resImages, resPhotos) => {
        let images = resImages.data
        let photos = resPhotos.data
        images.forEach((image, index) => {
          let isPhoto = () => {
            return photos.filter(photo => {
              return photo.name === image.name
            }).length !== 0
          }
          images[index].isPhoto = isPhoto()
        })
        this.images = images
      })).catch(e => {
        console.log(e.response.data)
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
      axios.put(process.env.API_URL_BASE + '/v1/images/putPhoto', {
        name: name
      }, {
        headers: {
          Authorization: localStorage.getItem('token')
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
      axios.put(process.env.API_URL_BASE + '/v1/images/removePhoto', {
        name: name
      }, {
        headers: {
          Authorization: localStorage.getItem('token')
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
