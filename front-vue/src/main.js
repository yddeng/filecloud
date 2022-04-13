import Vue from 'vue'
import App from './App.vue'
import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/antd.css';

Vue.config.productionTip = false
Vue.use(Antd);


import VueClipboard from 'vue-clipboard2'
Vue.use(VueClipboard)

import moment from 'moment'
// http://momentjs.cn/docs/#/use-it/
moment.locale('zh-cn')

new Vue({
  render: h => h(App)
}).$mount('#app')
