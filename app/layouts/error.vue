<template>
  <div class="center container">
    <p>{{ errorMessage }}</p>
    <br>
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
    error: {
      type: Object,
      default: () => {
        return {
          statusCode: 404,
          message: 'Page Not Found'
        }
      }
    }
  },
  computed: {
    errorMessage() {
      return `${this.error.statusCode} ${this.error.message}`
    }
  }
}
</script>

<style lang="scss" scoped>
@import '~/assets/scss/styles.scss';
.center {
  position: relative;
}
.link {
  text-decoration: underline;
}
</style>
