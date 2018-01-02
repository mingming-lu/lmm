<template>
  <div id="app">
    <header class="lmm-panel lmm-center" style="padding:128px 16px 32px 16px">
      <h1 class="lmm-xxxlarge">一个一点也不萌的网站</h1>
      <h1>卢明鸣的博客</h1>
  
      <div class="lmm-padding-32">
        <div class="lmm-bar lmm-border">
          <router-link to="/home" class="lmm-bar-item lmm-button">Home</router-link>
          <router-link to="/articles" class="lmm-bar-item lmm-button">Articles</router-link>
          <router-link to="/photos" class="lmm-bar-item lmm-button">Photographs</router-link>
          <router-link to="/projects" class="lmm-bar-item lmm-button">Projects</router-link>
          <router-link to="/profile" class="lmm-bar-item lmm-button">Profile</router-link>
        </div>
      </div>
    </header>

    <router-view/>

    <footer>
      <hr class="lmm-level-opacity">
      <div class="lmm-container lmm-padding-16 lmm-center">
        <a href="https://github.com/akinaru-lu" target="_blank"><i class="fa fa-github lmm-hover-opacity lmm-xxxlarge lmm-white"></i></a>
        <p>Powered by <router-link to="/profile" class="mm-light-grey lmm-white lmm-hover-light-grey">Lu Mingming</router-link></p>
      </div>
    </footer>
  </div>
</template>

<script>
export default {
  name: 'app'
}
</script>

<style>
  @import './assets/styles/common.css'
</style>
