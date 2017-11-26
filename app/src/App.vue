<template>
  <div id="app">
    <header class="lmm-panel lmm-center" style="padding:128px 16px">
      <h1 class="lmm-xlarge">一个一点也不萌的网站</h1>
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

    <footer class="lmm-container lmm-padding-64 lmm-light-grey lmm-center">
      <i class="fa fa-wechat lmm-hover-opacity"></i>
      <i class="fa fa-qq lmm-hover-opacity"></i>
      <i class="fa fa-weibo lmm-hover-opacity"></i>
      <i class="fa fa-facebook-official lmm-hover-opacity"></i>
      <i class="fa fa-twitter lmm-hover-opacity"></i>
      <p>Powered by <router-link to="/profile" class="lmm-light-grey lmm-hover-light-grey">Lu Mingming</router-link></p>
    </footer>
  </div>
</template>

<script>
export default {
  name: 'app'
}
</script>

<style>
.lmm-large {
  font-size: 18px !important;
}
.lmm-xlarge {
  font-size: 24px !important;
}
.lmm-padding-32 {
  padding-top: 32px !important;
  padding-bottom: 32px !important;
}
.lmm-padding-64 {
  padding-top: 64px !important;
  padding-bottom: 64px !important;
}
.lmm-container, .lmm-penal {
  padding: 0.01em 16px;
  margin: 16px 0;
}
.lmm-bar {
  width: 100%;
  overflow: hidden;
}
.lmm-border {
  border: 1px solid #ccc !important;
}
.lmm-card-4 {
  box-shadow: 0 4px 10px 0 rgba(0,0,0,0.2),0 4px 20px 0 rgba(0,0,0,0.19)
}
.lmm-center .lmm-bar {
  display: inline-block;
  width: auto;
}
.lmm-bar:before, .lmm-bar:after, .lmm-container:after, .lmm-container:before, .lmm-row:after, .lmm-row:before {
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
  background-color: #ccc !important;
}
.lmm-hover-opacity:hover {
  opacity: 0.60;
}
.lmm-bar .lmm-button {
  white-space:normal;
}
.lmm-center {
  text-align: center !important;
}
.lmm-display-container {
  position: relative;
}
.lmm-grey {
  color: #424949;
  background-color:#9e9e9e!important}
.lmm-left {
  float: left;
}
.lmm-light-grey {
  color: #424949;
  background-color: #f1f1f1 !important;
}
.lmm-link {
  text-decoration:none;
}
.lmm-margin {
  margin: 16px !important;
}
.lmm-margin-right {
  margin-right: 16px !important;
}
.lmm-hover-light-grey:hover {
  color: #757575 !important;
}
.lmm-half {
  display: inline-block;
  width: 49%;
  margin: 0 4px;
  vertical-align: top;
}
.lmm-half>.lmm-box {
  width: 100%;
  margin: 4px 0;
}
.lmm-padding-large {
  padding: 12px 24px !important;
}
.lmm-right {
  float: right !important;
}
.lmm-round-xlarge {
  border-radius: 16px;
}
.lmm-tag {
  display: inline-block;
  text-align: center;
  margin-bottom: 8px !important;
}
.lmm-white {
  color: #424949;
  background-color: #fff !important;
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
  color: #424949;
  background-color: #fff;
}
</style>
