<template>
    <el-form ref="userinfos" :model="userinfo" label-width="80px">
        <el-form-item label="用户名称">
            <el-input v-model="userinfo['username']" style="width:300px" disabled></el-input>
        </el-form-item>
        <el-form-item label="手机号">
            <el-input v-model="userinfo['phone']" style="width:300px" disabled></el-input>
        </el-form-item>
        <el-form-item label="邮箱">
            <el-input v-model="userinfo['email']" style="width:300px" disabled></el-input>
        </el-form-item>
        <el-form-item label="介绍">
            <el-input v-model="userinfo['context']" style="width:300px" disabled></el-input>
        </el-form-item>
        <el-form-item label="个性签名">
            <el-input v-model="userinfo['autograph']" style="width:300px" disabled></el-input>
        </el-form-item>
    </el-form>
</template>

<script>
export default {
    data(){
        return {
            userinfo: {}
        }
    },
    created(){
        this.userinfo = JSON.parse(window.sessionStorage.getItem("userinfo"))
    }
}
</script>

<style scoped>

</style>