import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      redirect: '/filecloud',
    },
    {
      name: 'filecloud',
      path: '/filecloud',
      meta:{title:'个人云盘'},
      component: () => import('@/views/Filecloud')
    },
    {
      name: 'fileshare',
      path: '/shared/s/:key',
      meta:{title:'个人云盘分享'},
      component: () => import('@/views/FileShared')
    },
    {
      path: '*',
      redirect: '/404',
    },
    {
      path: '/404',
      component: () => import('@/views/404')
    }
  ]
})