<template>
  <div id="body">
    <div id="content">
      <h2 style="text-align: center">个人云盘</h2>
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
  </div>
</template>

<script>
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
          this.$store.dispatch("Login",this.userInfo)
          .then(()=>{
            setTimeout(() => {
              this.loadingLogin = false
              this.$router.push({ name: 'filecloud' })
            }, 500)
          })
          .catch(() => {
            this.loadingLogin = false
          })
        }
      })
    }
  }
}
</script>
<style>
#body{
  padding-top:200px
}

#content{
  margin:0 auto;
  width:500px;
  height:300px;
  box-shadow: 0 1px 20px rgba(0,0,0,.2);
  border:2px;
  border-radius:4px;
  padding: 20px 60px;
}
</style>