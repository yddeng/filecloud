import axios from 'axios'
import { VueAxios } from './axios'
import message from 'ant-design-vue/es/message'
import { target } from '@/config/config'

// 创建 axios 实例
const request = axios.create({
  // API 请求的默认前缀
  baseURL: target,
  timeout: 6000 // 请求超时时间
})

// 异常拦截处理器
const errorHandler = (error) => {
  if (error.response){
    //if (error.response.status === 400){
      message.error(error.response.statusText);
    //}
  }
  return Promise.reject(error)
}

// request interceptor
request.interceptors.request.use(config => {
  return config
}, errorHandler)

// response interceptor
request.interceptors.response.use((response) => {
  if ('success' in response.data){
    if (!response.data.success) { // 成功
      message.error(response.data.message);
      return Promise.reject(response)
    } else {
      return response.data.data
    }
  }else{
    return response
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
