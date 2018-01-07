<template>
  <div id="app">
    <header class="white center">
      <nav>
        <router-link v-for="item in items" :key="item.name" :to="item.link" active-class="nav-item-active" class="nav-item">
          {{ item.name }}
        </router-link>
      </nav>
    </header>

    <router-view style="margin-bottom:86px;" />

    <footer class="center">
      <hr class="opacity-plus">
      <a href="https://github.com/akinaru-lu" target="_blank">
        <i class="fa fa-github link white"></i>
      </a>
      <p>&copy; <router-link to="/profile" class="link white">Lu Mingming</router-link></p>
    </footer>
  </div>
</template>

<script>
export default {
  name: 'app',
  data () {
    return {
      items: [
        {
          link: '/home',
          name: 'Home'
        },
        {
          link: '/articles',
          name: 'Articles'
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
  }
}
</script>

<style>
  @import './assets/styles/common.css'
</style>
