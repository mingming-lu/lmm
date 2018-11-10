<template>
  <div id="app">
    <Header/>
    <nuxt id="content" :style="{marginBottom: footerHeight + 'px'}"/>
    <Footer ref="footer"/>
  </div>
</template>

<script>
import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'

export default {
  components: {
    Header,
    Footer
  },
  data() {
    return {
      footerHeight: 0
    }
  },
  mounted() {
    this.calcFooterHeight()
    window.addEventListener('resize', this.calcFooterHeight)
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.calcFooterHeight)
  },
  methods: {
    calcFooterHeight() {
      this.footerHeight = this.$refs.footer.$el.clientHeight
    }
  }
}
</script>

<style lang="scss" scoped>
@import '~/assets/scss/styles.scss';
#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: $color_text;
  background-color: $color_primary;
}
#content {
  @media screen and (min-width: $max_width_device + 1) {
    margin-top: 128px;
  }
  @media screen and (max-width: 960px) {
    margin-top: 64px;
  }
  @media screen and (max-width: $max_width_device) {
    margin-top: 0;
  }
}
</style>
