import Vue from 'vue'
import Router from 'vue-router'

const originalPush = Router.prototype.push
Router.prototype.push = function push (location, onResolve, onReject) {
  if (onResolve || onReject) return originalPush.call(this, location, onResolve, onReject)
  return originalPush.call(this, location).catch(err => err)
}

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      redirect: '/filecloud',
    },
    {
      name:'filecloud',
      path: '/filecloud',
      meta:{title:'个人云盘'},
      component: () => import('@/views/Filecloud')
    },
    {
      name:"login",
      path: '/filecloud/login',
      meta:{title:'个人云盘登陆'},
      component: () => import('@/views/Login')
    },
    {
      name:'fileshared',
      path: '/shared/s/:key',
      meta:{title:'个人云盘分享'},
      component: () => import('@/views/FileShared')
    },
    {
      name:'sharedlist',
      path: '/shared/list',
      meta:{title:'我的分享'},
      component: () => import('@/views/SharedList')
    },
    {
      path: '*',
      redirect: '/404',
    },
    {
      path: '/404',
      meta:{title:'404'},
      component: () => import('@/views/404')
    }
  ]
})