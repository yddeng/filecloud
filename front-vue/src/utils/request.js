import axios from 'axios'
import { VueAxios } from './axios'
import message from 'ant-design-vue/es/message'
import storage from 'store'

// 创建 axios 实例
const request = axios.create({
  // API 请求的默认前缀
  baseURL: window.g.baseUrl,
  timeout: 6000 // 请求超时时间
})

// 异常拦截处理器
const errorHandler = (error) => {
  //console.log(error);
  if (error.response){
    if (error.response.status === 401){
      // token 验证失败
      message.error('授权验证失败');
      const token = storage.get("Access-Token")
      if(token){
        //storage.commit("SetToken",'')
        storage.remove("Access-Token")
      }
      setTimeout(() => {
        window.location.reload()
      }, 1000)
    }else{
      message.error(error.response.statusText);
    }
  }
  return Promise.reject(error)
}


// request interceptor
request.interceptors.request.use(config => {
  const token = storage.get("Access-Token")
  // 如果 token 存在
  // 让每个请求携带自定义 token
  //console.log(token);
  if (token) {
    config.headers["Access-Token"] = token
  }
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
