import Vue from 'vue'
import Router from 'vue-router'

const Home = () => import('@/components/Home')
const Index = () => import('@/components/Index')
const Blog = () => import('@/components/blog/Blog')
const BlogList = () => import('@/components/blog/List')
const Photos = () => import('@/components/Photos')
const Projects = () => import('@/components/Project')
const Reviews = () => import('@/components/Reviews')

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
      path: '/reviews',
      component: Reviews
    }
  ]
})
