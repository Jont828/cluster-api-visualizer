import Vue from 'vue'
import VueRouter from 'vue-router'
import ClusterOverview from '../views/ClusterOverview.vue'
import TargetCluster from '../views/TargetCluster.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'ClusterOverview',
    component: ClusterOverview
  },
  {
    path: '/clusters/',
    name: 'TargetCluster',
    component: TargetCluster,
    props: true
  },
  {
    path: '*',
    component: ClusterOverview
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
