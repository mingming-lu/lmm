<template>
  <header class="shadow">
    <nav v-if="wideMode" class="top-nav">
      <router-link to="/">
        <div class="logo">
          <img class="icon" src="/static/img/logo.png">
          明鳴的树洞
        </div>
      </router-link>

      <div :class="{narrowTopNav: moderateWideMode}">
        <router-link v-for="item in items" :key="item.name" :to="item.link" class="nav-item">
          {{ item.name }}
        </router-link>
      </div>
    </nav>

    <nav v-if="!wideMode" class="drawer-nav" >
      <div class="text-right" ref="drawerNavBar">
        <router-link to="" class="toggler" @click.native="toggleDrawer">
          <i v-if="drawerShown" class="fa fa-fw fa-times"></i>
          <i v-else class="fa fa-fw fa-bars"></i>
        </router-link>
      </div>
      <div class="drawer animate-left" :class="[drawerShown && !wideMode ? 'drawer-show' : 'drawer-hide']" :style="{marginTop: drawerNavBarHeight - 1 + 'px'}">
        <div class="container">
          <router-link v-for="item in items" :key="item.name" :to="item.link" active-class="drawer-item-active" class="link" @click.native="toggleDrawer">
            <p><i class="fa fa-fw" :class="item.icon"></i>{{ item.name }}</p>
          </router-link>
        </div>
      </div>
    </nav>
  </header>
</template>

<script>
export default {
  created () {
    window.addEventListener('resize', this.calcDrawerNavBarHeight)
    window.addEventListener('resize', this.calcIsWideMode)
    window.addEventListener('resize', this.calcIsModerateWideMode)
  },
  mounted () {
    this.calcDrawerNavBarHeight()
    this.calcIsWideMode()
  },
  beforeDestroy () {
    window.removeEventListener('resize', this.calcDrawerNavBarHeight)
    window.removeEventListener('resize', this.calcIsWideMode)
    window.removeEventListener('resize', this.calcIsModerateWideMode)
  },
  data () {
    return {
      drawerNavBarHeight: 0,
      drawerShown: false,
      wideMode: false,
      moderateWideMode: false,
      items: [
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
    calcDrawerNavBarHeight () {
      this.drawerNavBarHeight = this.$refs.drawerNavBar.clientHeight
    },
    calcIsWideMode () {
      this.wideMode = window.innerWidth >= 680 // $width_max_drawer_to_view + 1
    },
    calcIsModerateWideMode () {
      this.moderateWideMode = window.innerWidth <= 960
    },
    toggleDrawer () {
      this.drawerShown = !this.drawerShown
    }
  }
}
</script>

<style lang="scss" scoped>
@import '@/assets/scss/styles.scss';
.narrowTopNav {
  display: inline-block;
  width: 310px;
}
.logo {
  cursor: pointer;
  float: left;
  outline:none;
  user-select: none;
  color: $color_text;
  padding: 0 16px;
  vertical-align: middle;
  .icon {
    width: 64px;
    height: 64px;
  }
}
header {
  background-color: $color_primary_dark;
  border: none;
  font-size: 1.5em;
  @media screen and (min-width: $width_max_drawer_to_view + 1) {
    padding: 48px;
  }
  @media screen and (max-width: $width_max_drawer_to_view) {
    position: sticky;
    top: 0;
  }
}
.drawer-nav {
  .toggler {
    color: $color_text;
    display: inline-block;
    padding: 16px;
  }
  .drawer {
    height: 100%;
    width: 100%;
    top: 0;
    left: 0;
    background-color: $color_primary_dark;
    position: fixed !important;
    overflow:auto;
    &.drawer-show {
      display: block;
    }
    &.drawer-hide {
      display: none;
    }
    .container {
      padding: 0 1em;
      i {
        margin-right: 8px;
      }
    }
    .drawer-item-active {
      color: $color_accent;
    }
  }
}
nav {
  &.top-nav {
    text-align: right;
    .nav-item {
      margin: 0 16px;
    }
  }
  .nav-item {
    border: none;
    display: inline-block;
    outline: 0;
    padding: 16px;
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
    &:hover {
      opacity: 0.8;
      color: $color_accent;
      transition: all 0.3s ease-out;
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
