import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/Home'
import ArticlesIndex from '@/components/articles/Index'
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
      component: ArticlesIndex,
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
