import Vue from 'vue'
import Router from 'vue-router'
import login from '@/components/UserLogin'
import main from '@/components/main'

Vue.use(Router)

var routes = [
  {
    path: '/',
    name: 'main',
    component: main,
    children: [
      {
        path: 'user/usercenter',
        component: () => import(`@/components/UserCenter`)
      }
    ]
  },
  {
    path: '/login',
    name: 'login',
    component: login
  },
  {
    path: '/register',
    name: 'register',
    component: () => import(`@/components/UserRegister`)
  }
];

export default new Router({routes})
