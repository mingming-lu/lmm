import Vue from 'vue'
import Router from 'vue-router'
import Index from '@/components/Index'
import Posts from '@/components/Posts'
import Photographs from '@/components/Photographs'
import Profile from '@/components/Profile'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      component: Index
    },
    {
      path: '/posts',
      component: Posts
    },
    {
      path: '/photographs',
      component: Photographs
    },
    {
      path: '/profile',
      component: Profile
    }
  ]
})
