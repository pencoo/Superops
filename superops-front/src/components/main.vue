<template>
  <el-container style="height: 100%;">
    <el-aside width="250px">
      <navmenu></navmenu>
    </el-aside>
    <el-container>
      <el-header height="40px">
        <div><i class="el-icon-s-home"></i><a href="">首页</a></div>
        <vheader />
      </el-header>
      <el-main>
        <router-view></router-view>
      </el-main>
      <el-footer height="40px">@copy right by pencoo. email: 919182546@qq.com</el-footer>
    </el-container>
  </el-container>
</template>

<script>
import NavMenu from "./NavMenu";
import Header from "./Header";
export default {
  name: "App",
  components: {
    navmenu: NavMenu,
    vheader: Header,
  },
};
</script>

<style scoped>
.el-container{
    height:100%
}
.el-header {
  background-color: #f3f3f4;
  text-align: right;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.el-aside {
  background: #32404f;
  overflow-x: hidden;
  overflow-y: scroll;
}
.el-main {
  color: red;
  height: 900px;
}
.el-footer {
  background: #f3f3f3;
  text-align: center;
  height: 40px;
  line-height: 40px;
}
</style>
