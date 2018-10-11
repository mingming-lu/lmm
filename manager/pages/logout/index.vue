<template>
  <div>Redirect to login page...</div>
</template>

<script>
export default {
  mounted() {
    window.localStorage.removeItem('accessToken')
    setTimeout(() => {
      this.$router.push('/login')
    }, 3000)
  }
}
</script>
