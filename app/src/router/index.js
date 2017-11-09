import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/Home'
import Articles from '@/components/Articles'
import Photos from '@/components/Photos'
import Projects from '@/components/Project'
import Profile from '@/components/Profile'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      component: Home
    },
    {
      path: '/articles',
      component: Articles
    },
    {
      path: '/photos',
      component: Photos
    },
    {
      path: '/projects',
      component: Projects
    },
    {
      path: '/profile',
      component: Profile
    }
  ]
})
