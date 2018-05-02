<template>
  <footer class="container center">
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
.container {
  padding: 32px;
}
footer {
  i {
    margin: 4px;
  }
  background-color: $color_primary;
  color: $color_text;
  bottom: 0 !important;
  position: fixed !important;
  width: 100%;
  z-index: -99 !important;
}
</style>
