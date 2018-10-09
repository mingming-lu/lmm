import Vue from 'vue'
import Router from 'vue-router'
import Index from '@/components/Index'
import Signin from '@/components/Signin'
import PageNotFound from '@/components/404'
import BlogIndex from '@/components/blog/Index'
import BlogNew from '@/components/blog/New'
import BlogEdit from '@/components/blog/Edit'
import ImageIndex from '@/components/image/Index'
import ImageUpload from '@/components/image/Upload'

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
      path: '/signin',
      component: Signin
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
      path: '/image',
      component: {
        template: '<router-view/>'
      },
      children: [
        {
          path: '/',
          component: ImageIndex
        },
        {
          path: 'upload',
          component: ImageUpload
        }
      ]
    }
  ]
})
