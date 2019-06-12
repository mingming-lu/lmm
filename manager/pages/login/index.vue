<template>
  <v-layout justify-center>
    <v-form class="form">
      <v-text-field
        v-model="username"
        :rules="[rules.username.required, rules.username.min, rules.username.max]"
        label="User Name"
        counter
      />

      <v-text-field
        v-model="password"
        :append-icon="showPassword ? 'visibility_off' : 'visibility'"
        :rules="[rules.password.required, rules.password.min, rules.password.max]"
        :type="showPassword ? 'text' : 'password'"
        label="Password"
        counter
        @click:append="showPassword = !showPassword"
      />
      <v-btn 
        color="accent" 
        type="submit" 
        @click.prevent="login">Login</v-btn>
    </v-form>
  </v-layout>
</template>

<script>
export default {
  asyncData({ query, redirect }) {
    return {
      callback: () => {
        const location = query.redirect ? query.redirect : '/'
        redirect(location)
      },
      username: '',
      password: '',
      showPassword: false,
      rules: {
        username: {
          min: v => v.length >= 3 || 'Min 3 characters',
          max: v => v.length <= 18 || 'Max 18 characters',
          required: v => !!v || 'Required.'
        },
        password: {
          min: v => v.length >= 8 || 'Min 8 characters',
          max: v => v.length <= 250 || 'Max 250 characters',
          required: v => !!v || 'Required.'
        }
      }
    }
  },
  watchQuery: ['redirect'],
  methods: {
    login() {
      this.$axios
        .post(
          '/v1/auth/token',
          {
            grantType: 'basicAuth'
          },
          {
            headers: {
              Authorization:
                'Basic ' +
                btoa(
                  JSON.stringify({
                    username: this.username,
                    password: this.password
                  })
                )
            }
          }
        )
        .then(res => {
          window.localStorage.setItem('accessToken', res.data.access_token)
          this.$store.commit(
            'setAccessToken',
            window.localStorage.getItem('accessToken')
          )
          this.callback()
        })
        .catch(e => {
          alert(e.response.data.error)
        })
    }
  }
}
</script>

<style scoped>
.form {
  width: 512px;
  min-width: 256px;
  max-width: 512px;
}
</style>
