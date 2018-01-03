<template>
  <div>
    <div class="lmm-row">

      <!-- left (edit) -->
      <div class="lmm-left" style="width:50%">
        <div class="lmm-container">
        <input v-model="title" placeholder="title">
        <div align="right">
          <select v-model="category">
            <option disabled value="">Select one category</option>
            <option>Default</option>
            <option v-for="ctg in categories" :key="ctg.id">{{ ctg.name }}</option>
            <option>New category</option>
          </select>
        </div>
        <br>
        <textarea v-model="text" v-on:input="marked" placeholder="text"></textarea>
        </div>
      </div>

      <!-- right (preview) -->
      <div class="lmm-right" style="width:50%">
        <div class="lmm-container">
        <span id="title-preview">{{ title }}</span>
        <br>
        <br>
        <div id="text-preview" v-html="textPreview" v-hljs style="text-align:left"></div>
        </div>
      </div>

    </div>
    <br>
    <button v-on:click="submit">Submit</button>
  </div>
</template>

<script>
import Markdownit from 'markdown-it'
import axios from 'axios'

let md = new Markdownit({
  html: true,
  typographer: true
})

export default {
  data () {
    return {
      title: '',
      text: '',
      textPreview: '',
      category: '',
      categories: [],
      tags: ''
    }
  },
  created () {
    axios.get('http://api.lmm.local/articles/categories?user_id=1').then((res) => {
      this.categories = res.data
    }).catch((e) => {
      console.log(e)
    })
  },
  methods: {
    submit () {
      if (!confirm('Are you sure you want to submit?')) {
        return
      }
      axios.post('http://api.lmm.local/article?user_id=1', {
        title: this.title,
        text: this.text,
        category: this.category,
        tags: this.tags
      }).then((res) => {
        this.$router.push('/')
      }).catch((e) => {
        console.log(e)
      })
    },
    marked () {
      this.textPreview = md.render(this.text)
    }
  }
}
</script>

<style scoped>
textarea {
  width: 100%;
  max-width: 100%;
  resize: vertical;
  overflow: hidden;
}
#title-preview:before {
  content: 'title: ';
  opacity: 0.4;
}
#text-preview:before {
  content: 'text: ';
  opacity: 0.4;
}
</style>