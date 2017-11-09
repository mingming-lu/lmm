<template>
  <div id="app">
    <div class="w3-content" style="max-width:1500px">

      <header class="w3-panel w3-center w3-opacity" style="padding:128px 16px">
        <h1 class="w3-xlarge">一个一点也不萌的网站</h1>
        <h1>卢明鸣的博客</h1>
  
        <div class="w3-padding-32">
          <div class="w3-bar w3-border">
            <router-link to="/" class="w3-bar-item w3-button">Home</router-link>
            <router-link to="/articles" class="w3-bar-item w3-button">Articles</router-link>
            <router-link to="/photos" class="w3-bar-item w3-button">Photographs</router-link>
            <router-link to="/projects" class="w3-bar-item w3-button">Projects</router-link>
            <router-link to="/profile" class="w3-bar-item w3-button">Profile</router-link>
          </div>
        </div>
      </header>

      <router-view/>

      <footer class="w3-container w3-padding-64 w3-light-grey w3-center w3-large"> 
        <i class="fa fa-facebook-official w3-hover-opacity"></i>
        <i class="fa fa-instagram w3-hover-opacity"></i>
        <i class="fa fa-snapchat w3-hover-opacity"></i>
        <i class="fa fa-pinterest-p w3-hover-opacity"></i>
        <i class="fa fa-twitter w3-hover-opacity"></i>
        <i class="fa fa-linkedin w3-hover-opacity"></i>
        <p>Powered by <router-link to="/profile" class="w3-hover-text-grey">Lu Mingming</router-link></p>
      </footer>

    </div>
  </div>
</template>

<script>
export default {
  name: 'app'
}
</script>

<style>
body {
  font-family: 'Lucida Grande', Helvetica, Arial,
               'Microsoft Yahei', '微软雅黑', STXihei, '华文细黑',
               'Yu Gothic', YuGothic,
               sans-serif;
}
h1 {
  letter-spacing: 6px
}
.w3-row-padding img {margin-bottom: 12px}
#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
