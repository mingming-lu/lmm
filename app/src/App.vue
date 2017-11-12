<template>
  <div id="app">
    <header class="lmm-panel lmm-center lmm-opacity" style="padding:128px 16px">
      <h1 class="lmm-xlarge">一个一点也不萌的网站</h1>
      <h1>卢明鸣的博客</h1>
  
      <div class="lmm-padding-32">
        <div class="lmm-bar lmm-border">
          <router-link to="/" class="lmm-bar-item lmm-button">Home</router-link>
          <router-link to="/articles" class="lmm-bar-item lmm-button">Articles</router-link>
          <router-link to="/photos" class="lmm-bar-item lmm-button">Photographs</router-link>
          <router-link to="/projects" class="lmm-bar-item lmm-button">Projects</router-link>
          <router-link to="/profile" class="lmm-bar-item lmm-button">Profile</router-link>
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
</template>

<script>
export default {
  name: 'app'
}
</script>

<style>
.lmm-xlarge {
  font-size: 24px !important;
}
.lmm-padding-32 {
  padding-top: 32px !important;
  padding-bottom: 32px !important;
}
.lmm-penal {
  padding: 0.01em 16px;
  margin-top: 16px;
  margin-bottom: 16px;
}
.lmm-bar {
  width: 100%;
  overflow: hidden;
}
.lmm-border {
  border: 1px solid #ccc !important;
}
.lmm-center .lmm-bar {
  display: inline-block;
  width: auto;
}
.lmm-bar:before, .lmm-bar:after {
  content: "";
  display: table;
  clear: both;
}
.lmm-bar .lmm-bar-item {
  padding: 8px 16px;
  float: left;
  width: auto;
  border: none;
  outline: none;
  display: block;
}
.lmm-border {
  border: 1px solid #ccc !important;
}
.lmm-button {
  border: none;
  display: inline-block;
  outline: 0;
  padding: 8px 16px;
  vertical-align: middle;
  overflow: hidden;
  text-decoration: none;
  color: inherit;
  background-color: inherit;
  text-align: center;
  cursor: pointer;
  white-space: nowrap;
  -webkit-touch-callout: none;
  -webkit-user-select: none;
  -khtml-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
}
.lmm-button:hover {
  color: #000 !important;
  background-color: #ccc !important;
}
.lmm-bar .lmm-button {
  white-space:normal;
}
.lmm-center {
  text-align: center !important;
}
.lmm-opacity {
  opacity: 0.60;
}
body {
  font-family: 'Lucida Grande', Helvetica, Arial,
               'Microsoft Yahei', '微软雅黑', STXihei, '华文细黑',
               'Yu Gothic', YuGothic,
               sans-serif;
}
h1 {
  letter-spacing: 6px;
}
#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
