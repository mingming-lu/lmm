<template>
  <div id="app">
    <Header />
    <router-view style="margin-bottom:86px;" />
  </div>
</template>

<script>
import Header from '@/components/common/Header'
export default {
  name: 'app',
  components: {
    Header
  }
}
</script>

<style>
  @import './assets/styles/common.css';
</style>
