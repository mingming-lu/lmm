<template>
  <div id="app">
    <header class="white" :class="{'text-center': !isMobile, 'text-left': isMobile}">
      <nav>
        <router-link v-if="!isMobile" v-for="item in items" :key="item.name" :to="item.link" active-class="nav-item-active" class="nav-item">
          {{ item.name }}
        </router-link>
        <button v-if="isMobile" class="nav-item" @click="toggleDrawer">&#9776;</button>
      </nav>

      <div class="w3-sidebar w3-bar-block w3-animate-left container" style="display:none;z-index:5" id="mySidebar">
        <router-link v-for="item in items" :key="item.name" :to="item.link" class="nav-item link" @click.native="toggleDrawer">
          <p>{{ item.name }}</p>
        </router-link>
      </div>
    </header>

    <div class="w3-overlay w3-animate-opacity" @click="toggleDrawer" style="cursor:pointer" id="myOverlay"></div>
    <router-view style="margin-bottom:86px;" />

    <footer class="center">
      <hr class="opacity-plus">
      <a href="https://github.com/akinaru-lu" target="_blank">
        <i class="fa fa-github link white"></i>
      </a>
      <a href="https://www.linkedin.com/in/lumingming/" target="_blank">
        <i class="fa fa-linkedin link white"></i>
      </a>
      <p>&copy; 2018 <router-link to="/profile" class="link white"><u>Lu Mingming</u></router-link></p>
    </footer>
  </div>
</template>

<script>
export default {
  name: 'app',
  data () {
    return {
      drawerShown: false,
      isMobile: false,
      items: [
        {
          link: '/home',
          name: 'Home'
        },
        {
          link: '/blog',
          name: 'Blog'
        },
        {
          link: '/photos',
          name: 'Photographs'
        },
        {
          link: '/projects',
          name: 'Projects'
        },
        {
          link: '/profile',
          name: 'Profile'
        }
      ]
    }
  },
  created () {
    this.calcIsMobile()
    window.addEventListener('resize', this.calcIsMobile)
  },
  beforeDestroy () {
    window.removeEventListener('resize', this.calcIsMobile)
  },
  methods: {
    calcIsMobile () {
      this.isMobile = window.innerWidth <= 450
    },
    toggleDrawer () {
      this.drawerShown = !this.drawerShown

      if (this.drawerShown) {
        document.getElementById('mySidebar').style.display = 'block'
        document.getElementById('myOverlay').style.display = 'block'
      } else {
        document.getElementById('mySidebar').style.display = 'none'
        document.getElementById('myOverlay').style.display = 'none'
      }
    }
  }
}
</script>

<style>
  @import './assets/styles/common.css';
</style>
