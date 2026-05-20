<template>
  <el-container class="layout-container">
    <el-header class="app-header">
      <div class="header-content">
        <div class="brand">
          <div class="brand-mark">远</div>
          <div>
            <h1>远山公益后台管理系统</h1>
            <p>内容、活动与用户运营工作台</p>
          </div>
        </div>
        <div class="header-right">
          <span class="user-info" v-if="userStore.userInfo">
            {{ userStore.username }} · {{ userStore.roleDisplay || '管理员' }}
          </span>
          <el-dropdown @command="handleCommand">
            <el-button class="header-action" circle>
              <el-icon><Setting /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="changePassword">修改密码</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </el-header>

    <el-container>
      <el-aside width="236px" class="app-aside">
        <el-menu
          :default-active="activeMenu"
          router
          class="el-menu-vertical"
          background-color="transparent"
          text-color="#374151"
          active-text-color="#0f766e"
        >
          <template v-for="menu in sidebarMenus" :key="menu.id">
            <!-- 有子菜单的目录 -->
            <el-sub-menu v-if="menu.children?.length" :index="'dir-' + menu.id">
              <template #title>
                <el-icon><component :is="menu.icon || 'Folder'" /></el-icon>
                <span>{{ menuTitle(menu) }}</span>
              </template>
              <el-menu-item
                v-for="child in menu.children.filter(c => c.type === 2)"
                :key="child.id"
                :index="child.path"
              >
                <el-icon><component :is="child.icon || 'Document'" /></el-icon>
                <span>{{ menuTitle(child) }}</span>
              </el-menu-item>
            </el-sub-menu>

            <!-- 没有子菜单的直接菜单项 (type=2) -->
            <el-menu-item v-else-if="menu.type === 2" :index="menu.path">
              <el-icon><component :is="menu.icon || 'Document'" /></el-icon>
              <span>{{ menuTitle(menu) }}</span>
            </el-menu-item>
          </template>
        </el-menu>
      </el-aside>

      <el-main class="app-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Setting } from '@element-plus/icons-vue'
import { useUserStore } from '../store/user'
import { usePermissionStore } from '../store/permission'
import { resetDynamicRoutes } from '../router'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const permissionStore = usePermissionStore()

// 侧边栏只显示目录(type=1)和菜单(type=2)，不显示按钮权限(type=3)
const sidebarMenus = computed(() => {
  return permissionStore.menus.filter(m => m.type === 1 || m.type === 2)
})

const activeMenu = computed(() => {
  return route.path
})

const menuTitle = (menu) => {
  return menu.title || menu.name
}

const handleCommand = (command) => {
  if (command === 'logout') {
    logout()
  } else if (command === 'changePassword') {
    router.push('/admin/change-password')
  }
}

const logout = () => {
  userStore.logout()
  permissionStore.reset()
  resetDynamicRoutes()
  ElMessage.success('退出成功')
  router.push('/login')
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7f6;
}

.app-header {
  background: #ffffff;
  color: #111827;
  height: 64px;
  padding: 0 24px;
  flex-shrink: 0;
  border-bottom: 1px solid #e5e7eb;
  box-shadow: 0 1px 2px rgba(17, 24, 39, 0.04);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 64px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-mark {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: #0f766e;
  color: #ffffff;
  display: grid;
  place-items: center;
  font-weight: 700;
  box-shadow: 0 8px 18px rgba(15, 118, 110, 0.18);
}

.header-content h1 {
  margin: 0;
  font-size: 18px;
  line-height: 24px;
  font-weight: 700;
  letter-spacing: 0;
}

.header-content p {
  margin: 2px 0 0;
  color: #6b7280;
  font-size: 12px;
  line-height: 16px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-info {
  color: #4b5563;
  font-size: 14px;
  padding: 6px 10px;
  background: #f3f4f6;
  border-radius: 6px;
}

.header-action {
  border: 1px solid #e5e7eb;
  color: #374151;
}

.el-container {
  flex: 1;
  overflow: hidden;
}

.app-aside {
  background: #ffffff;
  height: calc(100vh - 64px);
  overflow-y: auto;
  border-right: 1px solid #e5e7eb;
  padding: 14px 10px;
}

.el-menu-vertical {
  height: 100%;
  border-right: none;
}

.el-menu-vertical :deep(.el-menu-item),
.el-menu-vertical :deep(.el-sub-menu__title) {
  border-radius: 8px;
  height: 42px;
  line-height: 42px;
  margin: 2px 0;
}

.el-menu-vertical :deep(.el-menu-item.is-active) {
  background: #e6f4f1;
  font-weight: 600;
}

.app-main {
  background: #f5f7f6;
  padding: 22px;
  height: calc(100vh - 64px);
  overflow-y: auto;
}

@media (max-width: 768px) {
  .app-header {
    padding: 0 14px;
  }

  .brand p {
    display: none;
  }

  .app-aside {
    width: 192px !important;
  }

  .app-main {
    padding: 14px;
  }
}
</style>
