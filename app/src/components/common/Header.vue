<template>
   <header class="white sticky header-shadow">
      <nav>

        <!-- wide mode-->
        <div>
          <router-link v-if="wideMode" v-for="item in items" :key="item.name" :to="item.link" class="nav-item">
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
</template>

<script>
export default {
  created () {
    this.calcIsWideMode()
    window.addEventListener('resize', this.calcIsWideMode)
  },
  beforeDestroy () {
    window.removeEventListener('resize', this.calcIsWideMode)
  },
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
          name: 'Photos',
          icon: 'fa-camera-retro'
        },
        {
          link: '/projects',
          name: 'Projects',
          icon: 'fa-archive'
        },
        {
          link: '/reviews',
          name: 'Reviews',
          icon: 'fa-star-half-o'
        }
      ]
    }
  },
  methods: {
    calcIsWideMode () {
      this.wideMode = window.innerWidth > 600
    },
    toggleDrawer () {
      this.drawerShown = !this.drawerShown
    }
  }
}
</script>

<style lang="scss" scoped>
.drawer-show {
  display: block;
}
.drawer-hide {
  display: none;
}
.drawer {
  .container {
    border-top: 1px solid #f1f1f1;
  }
}
.drawer-item-active {
  color: deepskyblue;
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
