import Vue from 'vue'
import App from './App.vue'
import Antd from 'ant-design-vue';
import router from './router'
import  storage from 'store'
import 'ant-design-vue/dist/antd.css';

Vue.config.productionTip = false
Vue.use(Antd);


import VueClipboard from 'vue-clipboard2'
Vue.use(VueClipboard)

// http://momentjs.cn/docs/#/use-it/
import moment from 'moment'
moment.locale('zh-cn')

const allowList = ['login','fileshared'] // no redirect allowList
const loginRoutePath = '/filecloud/login'

router.beforeEach((to, from, next) => {
  if (to.meta.title){
    document.title = to.meta.title
  }

  const token = storage.get("Access-Token")
  if (token) {
    next()
  } else {
    if (allowList.includes(to.name)) {
      // 在免登录名单，直接进入
      next()
    } else {
      next({ path: loginRoutePath, query: { redirect: to.fullPath } })
    }
  }
})

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
