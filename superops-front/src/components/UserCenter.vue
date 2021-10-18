<template>
    <div>
        <div>
            <el-tabs v-model="activeName" @tab-click="handleClick">
                <el-tab-pane label="个人中心" name="first"><uinfo v-if="componentlist.ucent"></uinfo></el-tab-pane>
                <el-tab-pane label="Api Token管理" name="second"><apitoken v-if="componentlist.apitoken"></apitoken></el-tab-pane>
            </el-tabs>
        </div>
    </div>
</template>
<script>
import uinfo from './UserCenterUserinfo'
import apitoken from './UserCenterApiToken'
export default {
    data(){
        return {
            activeName: 'first',
            componentlist: {
                ucent: true,
                apitoken: false
            }
        }
    },
    components:{
        uinfo,
        apitoken
    },
    methods: {
      handleClick(tab) {
          if (tab.name == "first") {
              this.componentlist.ucent = true
              this.componentlist.apitoken = false
          } else {
              this.componentlist.ucent = false
              this.componentlist.apitoken = true
          }
      }
    }
}
</script>
<style scoped>

</style>