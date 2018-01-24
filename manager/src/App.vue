<template>
  <div id="app">
    <div class="text-right">
      {{ username }}
      <span v-if="username"><button @click="logout">Logout</button></span>
    </div>
    <router-view/> </div>
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
      axios.get('http://api.lmm.local/verify', {
        headers: {
          'Authorization': localStorage.getItem('token')
        }
      }).then(res => {
        this.username = localStorage.getItem('username')
      }).catch(e => {
        // go to /login if 401 or other error catched
        this.$router.push('/login')
      })
    },
    logout () {
      axios.get('http://api.lmm.local/logout', {
        headers: {
          Authorization: localStorage.getItem('token')
        }
      }).then(res => {
        alert('logout!')
        localStorage.removeItem('token')
        localStorage.removeItem('username')
        this.username = ''
        this.$router.push('/login')
      }).catch(e => {
        alert(e, e.response.data)
      })
    }
  }
}
</script>
