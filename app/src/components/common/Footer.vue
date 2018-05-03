<template>
  <footer class="center">
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
  i {
    font-size: 1.23em;
    margin: 4px;
  }
  background-color: $color_primary;
  bottom: 0 !important;
  color: $color_text;
  padding-top: 32px;
  padding-bottom: 8px;
  position: fixed !important;
  width: 100%;
  z-index: -99 !important;
}
</style>
