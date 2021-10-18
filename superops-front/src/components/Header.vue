<template>
  <div>
    <!-- 用户头像下拉菜单 -->
    <el-row >
      <el-dropdown @command="handleCommand">
        <span class="flex">
          <el-avatar class="iconclass" icon="el-icon-user-solid"></el-avatar>{{userinfos["username"]}}
        </span>
        <el-dropdown-menu slot="dropdown" style="text-align: center">
          <el-dropdown-item icon="el-icon-user" command="getuserinfo">用户信息</el-dropdown-item>
          <el-dropdown-item icon="el-icon-lock" command="changepassword">修改密码</el-dropdown-item>
          <el-dropdown-item icon="el-icon-s-home" command="setuserinfo">用户设置</el-dropdown-item>
          <el-dropdown-item icon="el-icon-key" command="getapitoken">获取Api Token</el-dropdown-item>
          <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </el-row>
    <!-- 查看用户信息 -->
    <el-dialog title="个人信息" :visible.sync="displayUserinfo" width="400px" style="text-align: center">
        <span>用户名：{{userinfos["username"]}}</span><br>
        <span>手机号：{{userinfos["phone"]}}</span><br>
        <span>邮 箱：{{userinfos["email"]}}</span><br>
        <span>描 述：{{userinfos["context"]}}</span><br>
        <span>个性签名：{{userinfos["autograph"]}}</span><br>
        <span slot="footer" class="dialog-footer">
          <el-button @click="displayUserinfo = false">取 消</el-button>
          <el-button type="primary" @click="displayUserinfo = false">确 定</el-button>
        </span>
    </el-dialog>
    <!-- 修改用户密码 -->
    <el-dialog title="修改密码" :visible.sync="updateUserpass" width="400px" style="text-align: center">
        <el-form :model="ChangepassForm" status-icon :rules="rules" ref="ChangepassForm" label-width="100px" class="demo-ChangepassForm">
          <el-form-item label="新密码" prop="password">
            <el-input type="password" v-model="ChangepassForm.password" autocomplete="off" placeholder="密码长度8-20位" show-password></el-input>
          </el-form-item>
          <el-form-item label="确认密码" prop="checkPass">
            <el-input type="password" v-model="ChangepassForm.checkPass" autocomplete="off" placeholder="请再次输入密码确认" show-password></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="UpdatePassDo">提交</el-button>
            <el-button @click="resetUpdatePassForm">重置</el-button>
          </el-form-item>
        </el-form>
    </el-dialog>
    <!-- 修改个人信息 -->
    <el-dialog title="修改个人信息" :visible.sync="updateUserinfo" width="400px" style="text-align: center">
        <el-form :model="ChangeInfoForm" status-icon :rules="rules" ref="ChangeInfoForm" label-width="100px" class="demo-ChangepassForm">
          <el-form-item label="手机号" prop="phone">
            <el-input type="text" v-model="ChangeInfoForm.phone" autocomplete="off"></el-input>
          </el-form-item>
          <el-form-item label="邮箱" prop="email">
            <el-input type="text" v-model="ChangeInfoForm.email" autocomplete="off"></el-input>
          </el-form-item>
          <el-form-item label="描述">
            <el-input type="text" v-model="ChangeInfoForm.context" autocomplete="off"></el-input>
          </el-form-item>
          <el-form-item label="个性签名">
            <el-input type="text" v-model="ChangeInfoForm.autograph" autocomplete="off"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="UpdateUserInfo">提交</el-button>
            <el-button @click="resetUserInfo">重置</el-button>
          </el-form-item>
        </el-form>
    </el-dialog>
  </div>
</template>

<script>
export default {
  data() {
    var checkphone = (rule, value, callback) => {
        setTimeout(() => {
          if (!Number.isInteger(value)) {
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
          if (this.ChangepassForm.checkPass !== '') {
            this.$refs.ChangepassForm.validateField('checkPass');
          }
          callback();
        }
    };
    var validatePass2 = (rule, value, callback) => {
        if (value === '') {
          callback(new Error('请再次输入密码'));
        } else if (value !== this.ChangepassForm.password) {
          callback(new Error('两次输入密码不一致!'));
        } else {
          callback();
        }
    };
    return {
      userinfos: {},
      displayUserinfo: false,
      updateUserpass: false,
      updateUserinfo: false,
      ChangepassForm: {
        password: '',
        checkPass: ''
      },
      ChangeInfoForm: {
        phone: '',
        email: '',
        context: '',
        autograph: ''
      },
      rules: {
          password: [
            { min: 3, max: 30, message: '长度在8到30个字符', trigger: 'blur' },
            { validator: validatePass, trigger: 'blur' }
          ],
          checkPass: [
            { validator: validatePass2, trigger: 'blur' }
          ],
          email: [
              {required: true, message: '邮箱不能为空', trigger: 'blur'},
              {min: 5, max: 20, message: '邮箱在5-20位之间', trigger: 'blur'}
          ],
          phone: [
              {required: true, message: '手机号不能为空', trigger: 'blur'},
              {validator: checkphone , trigger: 'blur'}
          ]
      }
    };
  },
  methods: {
    handleCommand(command) {
      if (command === "getuserinfo") {
        this.displayUserinfo = true
      }
      if (command === "changepassword") {
        this.updateUserpass = true
      }
      if (command === "setuserinfo") {
        this.updateUserinfo = true
      }
      if (command === "getapitoken") {
        console.log(this.$router)
      }
      // 退出登录
      if (command === "logout") {
        // 删除后端token
        this.$axios({ method: "get", url: "/users/logout" });
        // 删除浏览器token
        window.sessionStorage.removeItem("Access-Token");
        // 跳转到登录页
        this.$router.push({ name: "login" });
      }
    },
    // this.$refs.xxx.resetFields()
    resetUpdatePassForm(){
      this.$refs.ChangepassForm.resetFields()
    },
    UpdatePassDo(){
      this.$refs.ChangepassForm.validate(valid=>{
        if (valid) {
          this.$axios({
            method: "post",
            url: '/users/chagespass',
            params: {
              password: this.ChangepassForm.password
            }
          }).then(res => {
            // 密码修改成功会自动退出需要重新登录
            if (res.data.code === 200) {
              this.$message({message: '恭喜你，密码修改成功。请重新登录', type: 'success'});
              this.ChangepassForm.password = ''
              this.ChangepassForm.checkPass = ''
              // 删除后端token
              this.$axios({ method: "get", url: "/users/logout" });
              // 删除浏览器token
              window.sessionStorage.removeItem("Access-Token");
              window.sessionStorage.removeItem("userinfo");
              // 跳转到登录页
              this.$router.push({ name: "login" });
            } else {
              this.$message({message: '密码修改失败：' + res.data.message , type: 'error'});
            }
          })
        }
      })
    },
    UpdateUserInfo(){
      this.$refs.ChangeInfoForm.validate(valid=>{
        if (valid) {
          this.$axios({
            method: 'post',
            url: '/users/chageinfo',
            data: this.ChangeInfoForm
          }).then(res=>{
            if (res.data.code === 200) {
              this.$message({message: '用户信息更新成功', type: 'success'});
              this.updateUserinfo = false
            } else {
              this.$message({message: '用户信息更新失败:' + res.data.message, type: 'error'});
            }
          })
        }
      })
    },
    resetUserInfo(){
      this.$refs.ChangeInfoForm.resetFields()
    }
  },
  created() {
    this.userinfos = JSON.parse(window.sessionStorage.getItem("userinfo"));
  },
};
</script>

<style scoped>
.el-container{
    height:100%
}
.el-icon-arrow-down {
  font-size: 18px;
}
.flex {
  display: flex;
  align-items: center;
}
.iconclass {
  margin-right: 10px;
}
</style>