<template>
    <div>
    <div style="display: flex;justify-content: center;margin-top: 150px">
      <el-card style="width: 450px">
        <div slot="header" class="clearfix">
          <span>用户注册</span>
        </div>
        <el-form ref="user" :model="user" :rules="rules" label-width="80px">
            <el-form-item label="用户名" prop="username">
                <el-input v-model="user.username" placeholder="请输入用户名"></el-input>
            </el-form-item>
            <el-form-item label="密码" prop="password">
                <el-input v-model="user.password" placeholder="请输入密码" show-password></el-input>
            </el-form-item>
            <el-form-item label="确认密码" prop="password2">
                <el-input v-model="user.password2" placeholder="请再次输入密码" show-password></el-input>
            </el-form-item>
            <el-form-item label="邮箱地址" prop="email">
                <el-input v-model="user.email" placeholder="请输入邮箱地址"></el-input>
            </el-form-item>
            <el-form-item label="手机号" prop="phone">
                <el-input v-model="user.phone" placeholder="请输入手机号"></el-input>
            </el-form-item>
            <el-form-item>
                <el-button type="success" @click="onSubmit">立即注册</el-button>
                <el-button type="primary" @click="Register">登录</el-button>
            </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script>
  export default {
    //单页面中不支持前面的data:{}方式
    data() {
      var checkphone = (rule, value, callback) => {
        setTimeout(() => {
          // if (!Number.isInteger(value)) {
          if (value.constructor === Number) {
            callback(new Error('手机号不为整形数字'));
          } else {
            if (value.length !== 11) {
              callback(new Error('手机号必须为11位'));
            } else {
              callback();
            }
          }
        }, 1000);
      };
      var validatePass = (rule, value, callback) => {
        if (value === '') {
          callback(new Error('请输入密码'));
        } else {
          if (value.length > 5 && value.length < 30) {
              callback()
          } else {
              callback('密码长度必须是5-30位字符串')
          }
        }
      };
      var validatePass2 = (rule, value, callback) => {
        if (value === '') {
          callback(new Error('请再次输入密码'));
        } else if (value !== this.user.password) {
          callback(new Error('两次输入密码不一致!'));
        } else {
          callback();
        }
      };
      return{
        user:{
          username: '',
          password: '',
          password2: '',
          email: '',
          phone: '',
        },
        rules: {
            username: [
                {required: true, message: '用户名不能为空', trigger: 'blur'},
                {min: 5, max: 20, message: '用户名在5-20位之间', trigger: 'blur'}
            ],
            password: [
                {required: true, message: '密码不能为空', trigger: 'blur'},
                {validator: validatePass, trigger: 'blur'}
            ],
            password2: [
                {required: true, message: '密码不能为空', trigger: 'blur'},
                {validator: validatePass2, trigger: 'blur'}
            ],
            email: [
                {required: true, message: '邮箱不能为空', trigger: 'blur'},
                {min: 5, max: 20, message: '邮箱在5-20位之间', trigger: 'blur'}
            ],
            phone: [
                {validator: checkphone , trigger: 'blur'}
            ]
        }
      }
    },
    methods:{
      onSubmit(){
        this.$refs.user.validate(valid => {
          if (valid) {
            this.$axios({
            method: 'post',
            url: '/register',
            data: JSON.stringify(this.user)
            }).then(
              res => {
                if (res.data.code === 200) {
                    this.$message({message: '恭喜你，注册成功。请登录', type: 'success'});
                    this.$router.push({name:'login'})
                } else {
                    this.$message({message: '抱歉注册失败: '+ res.data.message, type: 'error'});
                    this.$router.push({name:'register'})
                }
              }
            )
          } else {
            this.$message({message: '请检查输入内容是否符合要求', type: 'error'});
          }
        })
      },
      Register(){
          this.$router.push({name:'login'})
      }
    }
  }
</script>

<style scoped>
</style>