<template>
  <header class="shadow">

    <nav v-if="wideMode" class="top-nav">
      <router-link v-for="item in items" :key="item.name" :to="item.link" class="nav-item">
        {{ item.name }}
      </router-link>
    </nav>

    <nav v-if="!wideMode" class="drawer-nav">
      <div class="text-left">
        <router-link to="" class="nav-item" @click.native="toggleDrawer">&#9776;</router-link>
        <div class="drawer animate-left" :class="[drawerShown && !wideMode ? 'drawer-show' : 'drawer-hide']">
          <router-link to="" class="nav-item" @click.native="toggleDrawer">&#x2715;</router-link>
          <div class="container">
            <router-link v-for="item in items" :key="item.name" :to="item.link" active-class="drawer-item-active" class="link" @click.native="toggleDrawer">
              <p><i class="fa fa-fw" :class="item.icon"></i>{{ item.name }}</p>
            </router-link>
          </div>
        </div>
      </div>
    </nav>

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
@import '@/assets/scss/styles.scss';
header {
  width: 100%;
  background-color: white;
  border: none;
  position: -webkit-sticky;
  position: -moz-sticky;
  position: -ms-sticky;
  position: -o-sticky;
  position: sticky !important;
  z-index: 99 !important;
  top: 0 !important;
  &.shadow {
    box-shadow: 0 2px 5px 0 rgba(0,0,0,0.16), 0 2px 10px 0 rgba(0,0,0,0.12);
  }
}
i {
  margin-right: 8px;
}
.drawer-nav {
  .drawer {
    height: 100%;
    width: 100%;
    top: 0;
    left: 0;
    background-color: #fff;
    position: fixed !important;
    z-index: 99 !important;
    overflow:auto;
    &.drawer-show {
      display: block;
    }
    &.drawer-hide {
      display: none;
    }
    .container {
      padding: 0 8px;
      border-top: 1px solid #f1f1f1;
    }
    .drawer-item-active {
      color: $secondary_color;
    }
  }
}
nav {
  .nav-item {
    border: none;
    display: inline-block;
    outline: 0;
    padding: 16px;
    vertical-align: top;
    overflow: hidden;
    text-decoration: none;
    color: inherit;
    font-weight: 600;
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
    &:hover {
      opacity: 0.8;
      color: $secondary_color;
      background-color: #f1f1f1;
    }
  }
}
.animate-left {
  position: relative;
  animation:animateleft 0.4s
}
@keyframes animateleft {
  from {
    left: -300px;
    opacity:0;
  } to {
    left:0;
    opacity:1;
  }
}
</style>
