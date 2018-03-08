<template>
  <div class="center">

    <div class="bg bgimg-1">
      <div class="container display-middle">
        <h1 class="center">一个一点也不萌的网站</h1>
        <hr class="opacity-plus">
        <p class=""></p>
        <p class="center">没逼你看，爱看看不看滚</p>
        <p class="center">就是这么Dior，不服来咬我啊</p>
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
.display-middle{
  position:absolute;
  top:50%;
  left:50%;
  transform:translate(-50%,-50%);
  -ms-transform:translate(-50%,-50%);
}
</style>
