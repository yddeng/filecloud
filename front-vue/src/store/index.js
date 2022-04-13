import Vue from 'vue'
import Vuex from 'vuex'

import storage from 'store'
import { login } from '@/api/api'

Vue.use(Vuex)

export const  ACCESS_TOKEN = 'Access-Token'

export default new Vuex.Store({
  state: {
    token:''
  },

  mutations: {
    //这里是set方法
    setToken: (state, token) => {
      state.token = token
    },
  },

  getters:{
    //这里是get方法   
    getToken(state){
      return state.token
    }
  },        

  actions: {
    // 事件
    Login ({ commit }, userInfo) {
      return new Promise((resolve, reject) => {
        login(userInfo).then(response => {
          const result = response
          storage.set(ACCESS_TOKEN, result.token, 8 * 60 * 60 * 1000)
          commit('setToken', result.token)
          resolve()
        }).catch(error => {
          reject(error)
        })
      })
    },
  },

  modules: {}
})

