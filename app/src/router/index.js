import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/Home'
import Index from '@/components/Index'
import Article from '@/components/articles/Article'
import ArticlesList from '@/components/articles/List'
import Photos from '@/components/Photos'
import Projects from '@/components/Project'
import Profile from '@/components/Profile'

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
      path: '/articles',
      component: Index,
      children: [
        {
          path: '',
          component: ArticlesList
        },
        {
          path: ':id',
          component: Article
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
      path: '/profile',
      component: Profile
    }
  ]
})
