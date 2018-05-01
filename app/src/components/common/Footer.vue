<template>
  <footer class="center">
    <hr>
    <a href="https://github.com/akinaru-lu" target="_blank">
      <i class="fa fa-github link"></i>
    </a>
    <a href="https://www.linkedin.com/in/lumingming/" target="_blank">
      <i class="fa fa-linkedin link"></i>
    </a>
    <p>&copy; 2018 <router-link to="/home#author" class="link">
      <u>Lu Mingming</u>
    </router-link></p>
  </footer>
</template>

<style lang="scss" scoped>
@import '@/assets/scss/styles.scss';
footer {
  hr {
    opacity: 0.15;
  }
  i {
    margin: 4px;
  }
  bottom: 0 !important;
  position: fixed !important;
  width: 100%;
  z-index: -99 !important;
}
</style>
