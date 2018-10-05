import Vue from 'vue'
import Router from 'vue-router'

const Home = () => import('@/components/Home')
const Index = () => import('@/components/Index')
const Article = () => import('@/components/article/Article')
const Articles = () => import('@/components/article/Articles')
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
      path: '/articles',
      component: Index,
      children: [
        {
          path: '',
          component: Articles
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
      path: '/reviews',
      component: Reviews
    }
  ]
})
