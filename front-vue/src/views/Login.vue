<template>
  <div id="content">
    <div style="height:36px;line-height:36px;text-align:center">
      <img src="@/assets/logo.png" alt="logo" style="width:28px;height:28px;">&nbsp;
      <span style="font-size:24px;font-weight: 600;">个人云盘</span>
    </div>
    <br/>
    <a-form-model  ref="loginFormModel" :model="userInfo" :rules="rules">
      <a-form-model-item  prop="username" >
        <a-input size="large" v-model="userInfo.username" placeholder="用户名" @keyup.enter="handleLogin">
          <a-icon slot="prefix" type="user" />
        </a-input>
      </a-form-model-item>
      <a-form-model-item prop="password">
        <a-input-password size="large" v-model="userInfo.password" @keyup.enter="handleLogin">
          <a-icon slot="prefix" type="lock" />
        </a-input-password>
      </a-form-model-item>
      <a-form-model-item>
        <a-button size="large" type="primary" style="width:100%" @click="handleLogin" :loading="loadingLogin">登陆</a-button>
      </a-form-model-item>
    </a-form-model>
  </div>
</template>

<script>
import  storage from 'store'
import { login } from '@/api/api'

export default {
  data(){
    return {
      userInfo:{
        username:"",
        password:"",
      },
      rules: {
        username: [
          { required: true, message: '用户名必填', trigger: 'change' },
        ],
        password: [
          { required: true, message: '密码必填', trigger: 'change' }
        ],
      },
      loadingLogin:false,
      
    }
  },
  methods: {
    handleLogin(){
      this.$refs.loginFormModel.validate(valid => {
        if (valid) {
          //console.log(this.userInfo,this.$store);

          this.loadingLogin = true
          login(this.userInfo).then(ret => {
            storage.set("Access-Token", ret.token, 8 * 60 * 60 * 1000)
            setTimeout(() => {
              this.loadingLogin = false
              this.$router.push({ name: 'filecloud' })
            }, 500)
          }).catch(() => {
            this.loadingLogin = false
          })
        }
      })
    }
  }
}
</script>
<style>
#content{
  position:absolute;
  top:50%;
  left:50%;
  width:500px;
  height:300px;
  margin-top:-150px;
  margin-left:-250px;
  box-shadow: 0 1px 20px rgba(0,0,0,.2);
  border:2px;
  border-radius:4px;
  padding: 20px 60px;
}
</style>