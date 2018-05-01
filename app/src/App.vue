<template>
  <div id="app">
    <Header />
    <Index />
    <Footer />
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
  }
}
</script>

<style lang="scss" scoped>
@import "@/assets/scss/styles.scss";
#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: $font_color;
  background-color: white;
}
#index {
  margin-bottom: 86px;
}
</style>
