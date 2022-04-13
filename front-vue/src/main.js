import Vue from 'vue'
import App from './App.vue'
import Antd from 'ant-design-vue';
import router from './router'
import store from './store'
import  storage ,{ ACCESS_TOKEN} from 'store'
import 'ant-design-vue/dist/antd.css';

Vue.config.productionTip = false
Vue.use(Antd);


import VueClipboard from 'vue-clipboard2'
Vue.use(VueClipboard)

import moment from 'moment'
// http://momentjs.cn/docs/#/use-it/
moment.locale('zh-cn')

const loginRoutePath = '/filecloud/login'
router.beforeEach((to, from, next) => {
  if (to.meta.title){
    document.title = to.meta.title
  }

  const token = storage.get(ACCESS_TOKEN)
  console.log(token);
  if (token) {
    next()
  } else {
    if (to.path === loginRoutePath) {
      next()
    } else {
      next({ path: loginRoutePath, query: { redirect: to.fullPath } })
    }
  }
})

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
