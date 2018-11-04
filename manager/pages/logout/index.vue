<template>
  <p class="text-xs-center"><nuxt-link to="/">Go to home page</nuxt-link></p>
</template>

<script>
export default {
  created() {
    if (confirm('Are you really going to logout ?')) {
      window.localStorage.removeItem('accessToken')
      this.$store.commit('setAccessToken', undefined)
      this.$router.push('/')
    } else {
      this.$router.back()
    }
  },
}
</script>
