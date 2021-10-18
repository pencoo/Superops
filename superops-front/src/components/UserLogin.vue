<template>
    <div>
    <div style="display: flex;justify-content: center;margin-top: 150px">
      <el-card style="width: 400px">
        <div slot="header" class="clearfix">
          <span>用户登录</span>
        </div>
        <el-form ref="user" :model="user" :rules="rules" label-width="80px">
            <el-form-item label="用户名" prop="username">
                <el-input v-model="user.username" placeholder="请输入用户名"></el-input>
            </el-form-item>
            <el-form-item label="密码" prop="password">
                <el-input v-model="user.password" placeholder="请输入密码" show-password></el-input>
            </el-form-item>
            <el-form-item>
                <el-button type="success" @click="Login">立即登录</el-button>
                <el-button type="warning" @click="CleanAll">重置</el-button>
                <el-button type="primary" @click="Register">注册</el-button>
            </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script>
// import {mapState,mapMutations} from 'vuex'
  export default {
    //单页面中不支持前面的data:{}方式
    data() {
      return{
        user:{
          username: 'pengbilong',
          password: '309745197'
        },
        rules: {
            username: [
                {required: true, message: '用户名不能为空', trigger: 'blur'},
                {min: 5, max: 20, message: '用户名在5-20位之间', trigger: 'blur'}
            ],
            password: [
                {required: true, message: '密码不能为空', trigger: 'blur'},
                {min: 6, max: 30, message: '密码必须是6-30位', trigger: 'blur'}
            ]
        }
      }
    },
    // computed:{
        // ...mapState(['Token','UserInfo'])
    // },
    methods:{
      // ...mapMutations(['setToken','setUserInfo']),
      Login(){
        this.$axios({
            method: 'post',
            url: '/login',
            data: JSON.stringify(this.user)
        }).then(
            res => {
                if (res.data.code === 200) {
                  // this.setToken(res.data.data['Access-Token'])
                  // this.setUserInfo(res.data.data['userinfo'])
                  window.sessionStorage.setItem("Access-Token", res.data.data['Access-Token'])
                  window.sessionStorage.setItem("userinfo",JSON.stringify(res.data.data['userinfo']))
                  if (window.sessionStorage.getItem('Access-Token')) {
                    this.$router.push({name:'main'})
                  }
                } else {
                  // 登录失败
                  this.$message({message: '登录失败：'+ res.data.message, type: 'error'});
                }
            }
        )
      },
      CleanAll(){
          this.$refs.user.ResetFields()
      },
      Register(){
          this.$router.push({name:'register'})
      }
    },
    created(){
      if(window.sessionStorage.getItem("Access-Token")) {
        this.$router.push("/")
      }
    }
  }
</script>

<style scoped>
</style>