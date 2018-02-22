<template>
  <div id="app">
    <header class="white" :class="{'text-center': wideMode}">
      <nav>

        <!-- wide mode-->
        <div>
          <router-link v-if="wideMode" v-for="item in items" :key="item.name" :to="item.link" active-class="nav-item-active" class="nav-item">
            {{ item.name }}
          </router-link>
        </div>

        <!-- narrow mode-->
        <div class="text-left">
          <router-link v-if="!wideMode" to="" class="nav-item" @click.native="toggleDrawer">&#9776;</router-link>
          <div class="drawer animate-left" :class="[drawerShown && !wideMode ? 'drawer-show' : 'drawer-hide']">
            <router-link to="" class="nav-item" @click.native="toggleDrawer">&#9776;</router-link>
            <div class="container">
              <router-link v-for="item in items" :key="item.name" :to="item.link" active-class="drawer-item-active" class="link white" @click.native="toggleDrawer">
                <p><i class="fa fa-fw" :class="item.icon"></i>{{ item.name }}</p>
              </router-link>
            </div>
          </div>
        </div>

      </nav>
      <div class="overlay animate-opacity" :class="[drawerShown && !wideMode ? 'drawer-show' : 'drawer-hide']" @click="toggleDrawer"></div>
    </header>

    <div class="overlay animate-opacity" :class="[drawerShown && !wideMode ? 'drawer-show' : 'drawer-hide']" @click="toggleDrawer"></div>

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
// TODO set drawer items margin-top equals to the height of header
export default {
  name: 'app',
  data () {
    return {
      drawerShown: false,
      wideMode: false,
      items: [
        {
          link: '/home',
          name: 'Home',
          icon: 'fa-home'
        },
        {
          link: '/blog',
          name: 'Blog',
          icon: 'fa-pencil'
        },
        {
          link: '/photos',
          name: 'Photographs',
          icon: 'fa-camera-retro'
        },
        {
          link: '/projects',
          name: 'Projects',
          icon: 'fa-star-half-o'
        },
        {
          link: '/about',
          name: 'About',
          icon: 'fa-info'
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
  watch: {
    content () {
      this.$nextTick(() => {
        console.log('exo me?')
      })
    }
  },
  methods: {
    calcIsWideMode () {
      this.wideMode = window.innerWidth > 600
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
.drawer>.container {
  border-top: 1px solid #f1f1f1;
}
.drawer-item-active {
  color: crimson;
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
</style>
