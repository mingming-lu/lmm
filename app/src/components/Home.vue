<template>
  <div class="center">

    <div class="bg bgimg-1">
      <div class="container">
        <h1 class="center title">一个一点也不萌的网站</h1>
        <h1 class="center">卢明鸣的博客</h1>
      </div>
    </div>

  </div>
</template>

<script>
export default {
  created () {
    document.getElementsByTagName('header')[0].className = 'fixed'
  },
  destroyed () {
    document.getElementsByTagName('header')[0].className = 'sticky white header-shadow'
  }
}
</script>

<style scoped>
h1 {
  font-weight: 500;
  letter-spacing: 6px;
}
.bg {
  background-attachment: fixed;
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
}
.bgimg-1 {
  background-image: url('/static/img/bg.png');
  min-height: 100vh;
}
.title {
  margin-top: 128px;
}
</style>
