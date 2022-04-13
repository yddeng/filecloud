<template>
  <div id="content">
    <h2 style="text-align: center">个人云盘</h2>
    <br/>
    <a-input size="large" v-model="userInfo.username" placeholder="用户名">
      <a-icon slot="prefix" type="user" />
    </a-input><br/><br/>
    <a-input-password size="large" v-model="userInfo.password">
      <a-icon slot="prefix" type="lock" />
    </a-input-password><br/><br/><br/>
    <a-button size="large" type="primary" style="width:100%" @click="handleLogin" :loading="loadingLogin">登陆</a-button>
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
      loadingLogin:false,
      
    }
  },
  methods: {
    handleLogin(){
      if (this.userInfo.username === "" || this.userInfo.password === ""){
        return
      }
      //console.log(this.userInfo,this.$store);

      this.loadingLogin = true
      this.$store.dispatch("Login",this.userInfo)
      .then(()=>{
        setTimeout(() => {
          this.loadingLogin = false
          this.$router.push({ path: '/' })
        }, 500)
      })
      .catch(() => {
        this.loadingLogin = false
      })
    }
  }
}
</script>
<style>
#content{
  margin:100px auto;
  width:500px;
  height:320px;
  box-shadow: 0 1px 20px rgba(0,0,0,.2);
  border:2px;
  border-radius:4px;
  padding: 20px 60px;
}
</style>