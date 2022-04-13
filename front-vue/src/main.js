import Vue from 'vue'
import App from './App.vue'
import Antd from 'ant-design-vue';
import router from './router'
import 'ant-design-vue/dist/antd.css';

Vue.config.productionTip = false
Vue.use(Antd);


import VueClipboard from 'vue-clipboard2'
Vue.use(VueClipboard)

import moment from 'moment'
// http://momentjs.cn/docs/#/use-it/
moment.locale('zh-cn')

router.beforeEach((to, from, next) => {
  if (to.meta.title){
    document.title = to.meta.title
  }
  next()
})

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
