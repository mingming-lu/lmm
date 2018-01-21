<template>
  <div>
    <p></p>
    <form>
      <label for="username">User Name</label>
      <input type="text" size="32" id="username" v-model="username"><br>
      <label for="password">Password</label>
      <input type="password" size="32" v-model="password"><br>
      <input type="submit" value="Login" @click="login">
    </form>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  beforeCreate () {
    if (localStorage.getItem('token')) {
      this.$router.push('/')
    }
  },
  data () {
    return {
      username: '',
      password: ''
    }
  },
  methods: {
    login () {
      axios.post('http://api.lmm.local/login', {
        name: this.username,
        password: this.password
      }).then(res => {
        if (res.status !== 200) {
          alert('login failed')
          return
        }
        localStorage.setItem('username', res.data.name)
        localStorage.setItem('token', res.data.token)
      }).catch(e => {
        console.log(e.response.data)
      })
    }
  }
}
</script>
