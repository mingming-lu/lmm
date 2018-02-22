import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/Home'
import Index from '@/components/Index'
import Blog from '@/components/blog/Blog'
import BlogList from '@/components/blog/List'
import Photos from '@/components/Photos'
import Projects from '@/components/Project'
import About from '@/components/About'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '*',
      component: {
        template: '<p>404 Not Found</p>'
      }
    },
    {
      path: '/',
      redirect: {
        path: '/home'
      }
    },
    {
      path: '/home',
      component: Home
    },
    {
      path: '/blog',
      component: Index,
      children: [
        {
          path: '',
          component: BlogList
        },
        {
          path: ':id',
          component: Blog
        }
      ]
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
      path: '/about',
      component: About
    }
  ]
})
