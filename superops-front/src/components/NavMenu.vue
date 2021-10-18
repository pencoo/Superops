<template>
<div>
<el-row class="tac">
  <el-col :span="24">
    <el-menu
      default-active="2"
      class="el-menu-vertical-demo"
      unique-opened
      router
      @open="handleOpen"
      @close="handleClose"
      background-color="#545c64"
      text-color="#fff"
      active-text-color="#ffd04b">
      <el-menu-item index="/"><i class="el-icon-s-home"></i>首页</el-menu-item>
      <el-menu-item index="/user/usercenter"><i class="el-icon-user"></i>个人中心</el-menu-item>
    </el-menu>
  </el-col>
</el-row>
</div>
</template>

<script>
// import menu from '@/config/menu-config'
export default {
    data () {
      return {
        menu: ''
      }
    },
    methods: {
      handleOpen(key, keyPath) {
        console.log(key, keyPath);
      },
      handleClose(key, keyPath) {
        console.log(key, keyPath);
      }
    }
}
</script>

<style scoped>
.el-container{
    height:100%
}
</style>