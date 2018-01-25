import Vue from 'vue'
import Router from 'vue-router'
import Index from '@/components/Index'
import Login from '@/components/Login'
import PageNotFound from '@/components/404'
import BlogIndex from '@/components/blog/Index'
import BlogNew from '@/components/blog/New'
import BlogEdit from '@/components/blog/Edit'
import Photographs from '@/components/Photographs'
import Profile from '@/components/Profile'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '*',
      component: PageNotFound
    },
    {
      path: '/',
      component: Index
    },
    {
      path: '/login',
      component: Login
    },
    {
      path: '/blog',
      component: {
        template: '<router-view/>'
      },
      children: [
        {
          path: '/',
          component: BlogIndex
        },
        {
          path: 'new',
          component: BlogNew
        },
        {
          path: ':id',
          component: {
            template: '<router-view/>'
          },
          redirect: {
            path: ':id/edit'
          },
          children: [
            {
              path: 'edit',
              component: BlogEdit
            }
          ]
        }
      ]
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
