// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css'
import axios from 'axios'
import Vuex from 'vuex'
// import store from './store'

Vue.use(Vuex)
Vue.prototype.$axios = axios
axios.defaults.baseURL = "http://127.0.0.1:8080/api/v1"
Vue.use(ElementUI)
Vue.config.productionTip = false

// axios.get('/static/config.json').then((res) => {
//   Vue.prototype.base_url = res.base_url;
// })

//axios拦截器，用于在请求时添加token
axios.interceptors.request.use(function (config) {
  config.headers['Access-Token'] = window.sessionStorage.getItem('Access-Token');
  return config
}, function (error) {
  return Promise.reject(error)
})

// 路由守卫，用于判断登录状态
router.beforeEach((to, from, next) => {
  if (window.sessionStorage.getItem('Access-Token')) {
    if (to.path == '/login' || to.path == '/register') {
      next({name:"main"})
    }
    next()
  } else {
    if (to.path == '/login' || to.path == '/register') {
      next()
    } else {
      console.log("hello")
      next({
        name: 'login'
      })
    }
  }
})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
//   store,
  components: {
    App
  },
  template: '<App/>'
})
