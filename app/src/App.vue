<template>
  <div id="app">
    <Header />
    <Index id="content" :style="{marginBottom: footerHeight + 'px'}" />
    <Footer ref="footer" />
  </div>
</template>

<script>
import Header from '@/components/common/Header'
import Index from '@/components/Index'
import Footer from '@/components/common/Footer'
export default {
  name: 'app',
  components: {
    Header,
    Index,
    Footer
  },
  created () {
    window.addEventListener('resize', this.calcFooterHeight)
  },
  beforeDestroy () {
    window.removeEventListener('resize', this.calcFooterHeight)
  },
  mounted () {
    this.calcFooterHeight()
  },
  data () {
    return {
      footerHeight: 0
    }
  },
  methods: {
    calcFooterHeight () {
      this.footerHeight = this.$refs.footer.$el.clientHeight
    }
  }
}
</script>

<style lang="scss" scoped>
@import "@/assets/scss/styles.scss";
#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: $color_text;
  background-color: $color_primary;
}
#content {
  @media screen and (min-width: 680px) {
    margin-top: 128px;
  }
  @media screen and (max-width: 960px) {
    margin-top: 64px;
  }
  @media screen and (max-width: 659px) {
    margin-top: 0;
  }
}
</style>
