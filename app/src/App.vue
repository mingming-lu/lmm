<template>
  <div id="app">
    <Header />
    <router-view style="margin-bottom:86px;" />
    <Footer />
  </div>
</template>

<script>
import Header from '@/components/common/Header'
import Footer from '@/components/common/Footer'
export default {
  name: 'app',
  components: {
    Header,
    Footer
  }
}
</script>

<style>
  @import './assets/styles/common.css';
</style>
