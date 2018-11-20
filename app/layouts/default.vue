<template>
  <div id="app">
    <Header ref="header"/>
    <nuxt
      id="content"
      :style="{top: `${headerHeight}px`, marginBottom: `${footerHeight}px`}"
    />
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
      headerHeight: 0,
      footerHeight: 0,
    }
  },
  mounted() {
    window.addEventListener('resize', this.calcEdge)
    this.calcEdge()
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.calcEdge)
  },
  methods: {
    calcEdge() {
      this.headerHeight = this.$refs.header.$el.clientHeight
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
</style>
