<template>
  <header class="shadow">

    <nav v-if="wideMode" class="top-nav">
      <a href="/">
        <div class="logo">
          <img class="icon" src="/static/img/logo.png">
          明鳴的树洞
        </div>
      </a>

      <div :class="{narrowTopNav: moderateWideMode}">
        <router-link v-for="item in items" :key="item.name" :to="item.link" class="nav-item">
          {{ item.name }}
        </router-link>
      </div>
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
    window.addEventListener('resize', this.calcIsModerateWideMode)
  },
  beforeDestroy () {
    window.removeEventListener('resize', this.calcIsWideMode)
    window.removeEventListener('resize', this.calcIsModerateWideMode)
  },
  data () {
    return {
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
    calcIsWideMode () {
      this.wideMode = window.innerWidth > 680
    },
    calcIsModerateWideMode () {
      this.moderateWideMode = window.innerWidth < 960
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
  padding: 0 32px;
  vertical-align: middle;
  .icon {
    width: 64px;
    height: 64px;
  }
}
header {
  background-color: $color_primary_dark;
  border: none;
  position: -webkit-sticky;
  position: -moz-sticky;
  position: -ms-sticky;
  position: -o-sticky;
  z-index: 99 !important;
  top: 0 !important;
  padding: 48px;
  font-size: 1.5em;
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
    background-color: $color_primary_dark;
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
    // font-weight: 600;
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
