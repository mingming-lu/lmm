<template>
  <nuxt-link to="/login">{{ `Redirect to login page in ${countDown}s` }}</nuxt-link>
</template>

<script>
export default {
  data() {
    return {
      countDown: 3,
    }
  },
  mounted() {
    window.localStorage.removeItem('accessToken')
    const i = setInterval(() => {
      this.countDown--;
      if (this.countDown === 0) {
        clearInterval(i)
        this.$router.push('/login')
      }
    }, 1000);
  }
}
</script>
