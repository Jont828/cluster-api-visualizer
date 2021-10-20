import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import TargetCluster from '../views/TargetCluster.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/target-cluster/:id',
    name: 'TargetCluster',
    component: TargetCluster,
    props: true
  },
  {
    path: '*',
    component: Home
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
