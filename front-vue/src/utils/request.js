import axios from 'axios'
import { VueAxios } from './axios'
import notification from 'ant-design-vue/es/notification'

// 创建 axios 实例
const request = axios.create({
  // API 请求的默认前缀
  baseURL: "http://127.0.0.1:9987",
  timeout: 6000 // 请求超时时间
})

// 异常拦截处理器
const errorHandler = (error) => {
  return Promise.reject(error)
}

// request interceptor
request.interceptors.request.use(config => {
  return config
}, errorHandler)

// response interceptor
request.interceptors.response.use((response) => {
  if (!response.data.success) { // 成功
    notification.error({
      message: '请求出错',
      description: response.data.message
    })
    return Promise.reject(response)
  } else {
    return response.data.data
  }
}, errorHandler)

const installer = {
  vm: {},
  install (Vue) {
    Vue.use(VueAxios, request)
  }
}

export default request

export {
  installer as VueAxios,
  request as axios
}
