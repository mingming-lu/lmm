<template>
  <div id="app">
    <header class="white" :class="{'text-center': wideMode}">
      <nav>
        <router-link v-if="wideMode" v-for="item in items" :key="item.name" :to="item.link" active-class="nav-item-active" class="nav-item">
          {{ item.name }}
        </router-link>
        <button v-if="!wideMode" class="nav-item button" @click="toggleDrawer">&#9776;</button>

        <div class="drawer animate-left text-left" :class="[drawerShown ? 'drawer-show' : 'drawer-hide']">
          <button class="nav-item button close-button" @click="toggleDrawer">&times;</button>
          <div class="container" style="margin-top: 44px; /* height of header */">
            <router-link v-for="item in items" :key="item.name" :to="item.link" class="link white" @click.native="toggleDrawer">
              <p>{{ item.name }}</p>
            </router-link>
          </div>
        </div>
      </nav>
      <div class="overlay animate-opacity" :class="[drawerShown ? 'drawer-show' : 'drawer-hide']" @click="toggleDrawer"></div>
    </header>

    <div class="overlay animate-opacity" :class="[drawerShown ? 'drawer-show' : 'drawer-hide']" @click="toggleDrawer"></div>

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
      wideMode: false,
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
    this.calcIsWideMode()
    window.addEventListener('resize', this.calcIsWideMode)
  },
  beforeDestroy () {
    window.removeEventListener('resize', this.calcIsWideMode)
  },
  methods: {
    calcIsWideMode () {
      this.wideMode = window.innerWidth > 450
    },
    toggleDrawer () {
      this.drawerShown = !this.drawerShown
      console.log(this.drawerShown)
    }
  }
}
</script>

<style>
  @import './assets/styles/common.css';
</style>
<style scoped>
.drawer-show {
  display: block;
}
.drawer-hide {
  display: none;
}
.animate-left {
  position: relative;
  animation:animateleft 0.4s
}
@keyframes animateleft {
  from {
    left: -300px;
    opacity:0
  } to {
    left:0;opacity:1
  }
}
.close-button {
  position: absolute;
  right: 0;
  padding: 14px;
}
.button {
  font-size: 13px;
}
</style>
