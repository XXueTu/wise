import { createRouter, createWebHistory } from 'vue-router'
import ModelManager from '../components/ModelManager.vue'
import ResourceManager from '../components/ResourceManager.vue'

const routes = [
  {
    path: '/',
    redirect: '/resource'
  },
  {
    path: '/resource',
    name: 'Resource',
    component: ResourceManager
  },
  {
    path: '/model',
    name: 'Model',
    component: ModelManager
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router 