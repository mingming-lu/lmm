<template>
  <div>
    <p></p>
    <form>
      <label for="username">User Name</label>
      <input type="text" size="32" id="username" v-model="username"><br>
      <label for="password">Password</label>
      <input type="password" size="32" v-model="password"><br>
      <input type="submit" value="Signin" @click.prevent="signin">
    </form>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  data () {
    return {
      username: '',
      password: ''
    }
  },
  methods: {
    signin () {
      axios.post(process.env.API_URL_BASE + '/v1/auth/login', null, {
        headers: {
          Authorization: 'Basic ' + btoa(JSON.stringify({
            username: this.username,
            password: this.password
          }))
        }
      }).then(res => {
        if (res.status !== 200) {
          alert('signin failed')
          return
        }
        localStorage.setItem('token', res.data.accessToken)
        this.$router.push('/')
        location.reload()
      }).catch(e => {
        console.log(e.response.data)
      })
    }
  }
}
</script>
