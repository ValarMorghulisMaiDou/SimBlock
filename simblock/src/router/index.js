import Vue from 'vue'
import Router from 'vue-router'
import Main from '@/components/Main'
import Network from '@/components/Network'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      component: Main,
    },
    {
      path: '/network/:name',
      component: Network,
    },
  ]
})
