<template>
  <v-app dark>
    <v-navigation-drawer
      :mini-variant="miniVariant"
      :clipped="clipped"
      v-model="drawer"
      fixed
      app
    >
      <v-list>
        <v-list-tile
          v-for="(item, i) in items"
          :to="item.to"
          :key="i"
          router
          exact
        >
          <v-list-tile-action>
            <v-icon v-html="item.icon" />
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title v-text="item.title" />
          </v-list-tile-content>
        </v-list-tile>
      </v-list>
    </v-navigation-drawer>
    <v-toolbar
      :clipped-left="clipped"
      fixed
      app
    >
      <v-toolbar-side-icon @click="drawer = !drawer" />
      <v-toolbar-title v-text="title"/>
    </v-toolbar>
    <v-content>
      <v-container class="container">
        <nuxt />
      </v-container>
    </v-content>
    <v-footer
      :fixed="fixed"
      app
    >
      <span>&copy; 2018 Lu Mingming</span>
    </v-footer>
  </v-app>
</template>

<script>
export default {
  data() {
    return {
      clipped: true,
      drawer: true,
      fixed: true,
      items: [
        { icon: 'home', title: 'Home', to: '/' },
        { icon: 'create', title: 'Articles', to: '/articles' },
        { icon: 'photo_library', title: 'Assets (Images)', to: '/assets/images' },
        { icon: 'photo_camera', title: 'Assets (Photos)', to: '/assets/photos' },
        { icon: 'people', title: 'Users', to: '/users' },
        { icon: 'exit_to_app', title: 'Logout', to: '/logout'},
      ],
      miniVariant: true,
      title: 'Manager',
    }
  },
  mounted () {
    window.addEventListener('resize', this.onResize, {passive: true})
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.onResize, {passive: true})
  },
  methods: {
    logout() {
      window.localStorage.removeItem('accessToken')
    },
    onResize() {
      // see https://vuetifyjs.com/en/layout/breakpoints
      this.miniVariant = window.innerWidth > 600
    }
  }
}
</script>
<style scoped>
.container {
  font-size: 1.5em;
}
</style>
