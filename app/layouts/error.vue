<template>
  <div class="center container">
    <p>{{statusCode}} {{message}}</p>
    <nuxt-link
      to="/"
      class="link"
      >
      Home page
    </nuxt-link>
  </div>
</template>

<script>
export default {
  props: {
    statusCode: {
      type: Number,
      default: 404,
    },
    message: {
      type: String,
      default: 'Page Not Found',
    },
  },
}
</script>

<style lang="scss" scoped>
@import '~/assets/scss/styles.scss';
.center {
  position: relative;
}
.link{
  text-decoration: underline;
}
</style>