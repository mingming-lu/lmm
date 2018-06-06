<template>
  <div id="app">
    <div class="text-right" v-if="username">
      <router-link to="/" style="float: left">Home</router-link>
      {{ username }}
      <span><button @click="signout">Signout</button></span>
    </div>
    <router-view/>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'app',
  created () {
    this.verify()
  },
  data () {
    return {
      username: ''
    }
  },
  methods: {
    verify () {
      axios.get(process.env.API_URL_BASE + '/v1/verify', {
        headers: {
          'Authorization': 'Bearer ' + localStorage.getItem('token')
        }
      }).then(res => {
        this.username = localStorage.getItem('username')
      }).catch(e => {
        // go to /signin if 401 or other error catched
        this.clearLocalStorage()
        this.$router.push('/signin')
      })
    },
    signout () {
      this.clearLocalStorage()
      location.reload()
    },
    clearLocalStorage () {
      localStorage.removeItem('token')
      localStorage.removeItem('username')
    }
  }
}
</script>
