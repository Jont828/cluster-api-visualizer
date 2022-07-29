import Vue from 'vue'
import VueRouter from 'vue-router'
import ManagementCluster from '../views/ManagementCluster.vue'
import DescribeCluster from '../views/DescribeCluster.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'ManagementCluster',
    component: ManagementCluster
  },
  {
    path: '/clusters/',
    name: 'DescribeCluster',
    component: DescribeCluster,
    props: true
  },
  {
    path: '*',
    component: ManagementCluster
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
