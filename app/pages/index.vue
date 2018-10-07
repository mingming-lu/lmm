<template>
  <div class="background-image-1" />
</template>

<style lang="scss" scoped>
@import '~/assets/scss/styles.scss';
.background-image {
  background-attachment: fixed;
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
}
.background-image-1 {
  @extend .background-image;
  background-image: url('/img/home-bg1.png');
  min-height: 100vh;
}
</style>
