<template>
  <el-container class="app-container">
    <el-aside width="220px" class="app-aside">
      <div class="logo-container">
        <el-icon class="logo-icon"><Monitor /></el-icon>
        <h1 class="title">管理系统</h1>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="app-menu"
        :router="true"
        :collapse="isCollapse"
      >
        <el-menu-item index="resource">
          <el-icon><Document /></el-icon>
          <template #title>资源管理</template>
        </el-menu-item>
        <el-menu-item index="model">
          <el-icon><Connection /></el-icon>
          <template #title>模型管理</template>
        </el-menu-item>
      </el-menu>
    </el-aside>
    
    <el-container class="main-container">
      <el-header class="app-header">
        <div class="header-left">
          <el-button
            type="text"
            class="collapse-btn"
            @click="toggleCollapse"
          >
            <el-icon>
              <component :is="isCollapse ? 'Expand' : 'Fold'" />
            </el-icon>
          </el-button>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item>{{ currentPage }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-dropdown>
            <span class="user-info">
              <el-avatar :size="32" src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" />
              <span class="username">管理员</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>个人信息</el-dropdown-item>
                <el-dropdown-item>修改密码</el-dropdown-item>
                <el-dropdown-item divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-main class="app-main">
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { Connection, Document, Monitor } from '@element-plus/icons-vue'
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const isCollapse = ref(false)
const activeMenu = computed(() => route.path.slice(1) || 'resource')
const currentPage = computed(() => {
  const pathMap = {
    resource: '资源管理',
    model: '模型管理'
  }
  return pathMap[route.path.slice(1)] || '首页'
})

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}
</script>

<style scoped>
.app-container {
  height: 100vh;
  background-color: #f5f7fa;
}

.app-aside {
  background-color: #001529;
  transition: width 0.3s;
  overflow: hidden;
}

.logo-container {
  height: 64px;
  padding: 16px;
  display: flex;
  align-items: center;
  background-color: #002140;
}

.logo-icon {
  font-size: 24px;
  color: #fff;
  margin-right: 12px;
}

.title {
  color: #fff;
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  white-space: nowrap;
}

.app-menu {
  border-right: none;
  background-color: #001529;
}

.app-menu :deep(.el-menu-item) {
  color: rgba(255, 255, 255, 0.65);
}

.app-menu :deep(.el-menu-item.is-active) {
  color: #fff;
  background-color: #1890ff;
}

.app-menu :deep(.el-menu-item:hover) {
  color: #fff;
  background-color: #001f3f;
}

.main-container {
  display: flex;
  flex-direction: column;
}

.app-header {
  background-color: #fff;
  border-bottom: 1px solid #e8e8e8;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 64px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.collapse-btn {
  font-size: 20px;
  padding: 0;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 8px;
}

.username {
  font-size: 14px;
  color: #333;
}

.app-main {
  padding: 24px;
  overflow-y: auto;
}

:deep(.el-menu--collapse) {
  width: 64px;
}

:deep(.el-menu--collapse) .logo-container {
  padding: 16px 0;
  justify-content: center;
}

:deep(.el-menu--collapse) .title {
  display: none;
}
</style> 