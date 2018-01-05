import Vue from 'vue'
import Router from 'vue-router'
import Index from '@/components/Index'
import PageNotFound from '@/components/404'
import PostsIndex from '@/components/posts/Index'
import PostsNew from '@/components/posts/New'
import PostsEdit from '@/components/posts/Edit'
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
      path: '/posts',
      component: {
        template: '<router-view/>'
      },
      children: [
        {
          path: '/',
          component: PostsIndex
        },
        {
          path: 'new',
          component: PostsNew
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
              component: PostsEdit
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
