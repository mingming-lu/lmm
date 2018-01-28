<template>
  <div id="app">
    <header class="white" :class="{'text-center': !isMobile, 'text-left': isMobile}">
      <nav>
        <router-link v-if="!isMobile" v-for="item in items" :key="item.name" :to="item.link" active-class="nav-item-active" class="nav-item">
          {{ item.name }}
        </router-link>
        <button v-if="isMobile" class="nav-item" @click="drawerToggle">&#9776;</button>
      </nav>
    </header>

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
    }
  }
}
</script>

<style>
  @import './assets/styles/common.css';
</style>
