<template>
  <el-container class="layout">
    <!-- ============ Header ============ -->
    <el-header class="layout__header">
      <div class="header__inner">
        <div class="brand">
          <el-button
            v-if="isMobile"
            class="menu-toggle"
            text
            @click="drawerOpen = true"
          >
            <el-icon :size="20"><Expand /></el-icon>
          </el-button>
          <div class="brand__mark">远</div>
          <div class="brand__text">
            <h1>远山平台</h1>
            <p>内容 · 活动 · 用户 一体化工作台</p>
          </div>
        </div>

        <div class="header__right">
          <el-dropdown @command="handleCommand" trigger="click">
            <div class="user-chip">
              <div class="user-chip__avatar">
                {{ avatarText }}
              </div>
              <div class="user-chip__meta">
                <span class="user-chip__name">{{ userStore.username || '管理员' }}</span>
                <span class="user-chip__role">{{ userStore.roleDisplay || '系统管理员' }}</span>
              </div>
              <el-icon class="user-chip__caret"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="changePassword">
                  <el-icon><Key /></el-icon> 修改密码
                </el-dropdown-item>
                <el-dropdown-item command="logout" divided>
                  <el-icon><SwitchButton /></el-icon> 退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </el-header>

    <el-container class="layout__body">
      <!-- ============ Aside (desktop) ============ -->
      <el-aside
        v-if="!isMobile"
        width="244px"
        class="layout__aside"
      >
        <div class="aside__search">
          <el-input
            v-model="menuKeyword"
            placeholder="搜索菜单"
            clearable
            size="default"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>

        <el-scrollbar class="aside__scroll">
          <el-menu
            :default-active="activeMenu"
            router
            class="aside__menu"
          >
            <template v-for="menu in filteredMenus" :key="menu.id">
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

              <el-menu-item v-else-if="menu.type === 2" :index="menu.path">
                <el-icon><component :is="menu.icon || 'Document'" /></el-icon>
                <span>{{ menuTitle(menu) }}</span>
              </el-menu-item>
            </template>

            <div v-if="filteredMenus.length === 0" class="aside__empty">
              <el-icon><Search /></el-icon>
              <span>无匹配菜单</span>
            </div>
          </el-menu>
        </el-scrollbar>
      </el-aside>

      <!-- ============ Aside (mobile drawer) ============ -->
      <el-drawer
        v-if="isMobile"
        v-model="drawerOpen"
        direction="ltr"
        size="260px"
        :with-header="false"
      >
        <div class="aside__search">
          <el-input
            v-model="menuKeyword"
            placeholder="搜索菜单"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
        <el-menu
          :default-active="activeMenu"
          router
          class="aside__menu"
          @select="drawerOpen = false"
        >
          <template v-for="menu in filteredMenus" :key="menu.id">
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
            <el-menu-item v-else-if="menu.type === 2" :index="menu.path">
              <el-icon><component :is="menu.icon || 'Document'" /></el-icon>
              <span>{{ menuTitle(menu) }}</span>
            </el-menu-item>
          </template>
        </el-menu>
      </el-drawer>

      <!-- ============ Main ============ -->
      <el-main class="layout__main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed, ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  ArrowDown, Key, SwitchButton, Search, Expand
} from '@element-plus/icons-vue'
import { useUserStore } from '../store/user'
import { usePermissionStore } from '../store/permission'
import { resetDynamicRoutes } from '../router'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const permissionStore = usePermissionStore()

const menuKeyword = ref('')
const drawerOpen = ref(false)
const windowWidth = ref(window.innerWidth)

const isMobile = computed(() => windowWidth.value < 768)

const onResize = () => { windowWidth.value = window.innerWidth }
onMounted(() => window.addEventListener('resize', onResize))
onBeforeUnmount(() => window.removeEventListener('resize', onResize))

const sidebarMenus = computed(() =>
  permissionStore.menus.filter(m => m.type === 1 || m.type === 2)
)

const filteredMenus = computed(() => {
  const kw = menuKeyword.value.trim().toLowerCase()
  if (!kw) return sidebarMenus.value
  const match = (m) => (menuTitle(m) || '').toLowerCase().includes(kw)
  return sidebarMenus.value
    .map(menu => {
      if (menu.children?.length) {
        const kids = menu.children.filter(c => c.type === 2 && match(c))
        if (kids.length) return { ...menu, children: kids }
        if (match(menu)) return menu
        return null
      }
      return match(menu) ? menu : null
    })
    .filter(Boolean)
})

const activeMenu = computed(() => route.path)
const menuTitle = (menu) => menu.title || menu.name

const avatarText = computed(() => {
  const name = userStore.username || ''
  return name ? name.charAt(0).toUpperCase() : 'A'
})

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
/* ============ Container ============ */
.layout {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-page);
}

.layout__body {
  flex: 1;
  overflow: hidden;
}

/* ============ Header ============ */
.layout__header {
  background: linear-gradient(180deg, #ffffff 0%, var(--brand-50) 100%);
  height: 60px;
  padding: 0 var(--sp-6);
  flex-shrink: 0;
  border-bottom: 1px solid var(--border-1);
  box-shadow: var(--shadow-sm);
}

.header__inner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 60px;
}

.brand {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
}

.menu-toggle {
  margin-right: 4px;
}

.brand__mark {
  width: 36px;
  height: 36px;
  border-radius: var(--r-md);
  background: linear-gradient(135deg, var(--brand-600) 0%, var(--brand-800) 100%);
  color: #ffffff;
  display: grid;
  place-items: center;
  font-weight: 700;
  font-size: 16px;
  box-shadow: var(--shadow-brand);
  letter-spacing: -0.5px;
}

.brand__text h1 {
  margin: 0;
  font-size: 17px;
  line-height: 22px;
  font-weight: 700;
  color: var(--text-1);
  letter-spacing: 0;
}

.brand__text p {
  margin: 2px 0 0;
  color: var(--text-3);
  font-size: 12px;
  line-height: 14px;
}

/* ============ Header right ============ */
.header__right {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
}

.user-chip {
  display: flex;
  align-items: center;
  gap: var(--sp-2);
  padding: 4px 10px 4px 4px;
  border-radius: var(--r-full);
  cursor: pointer;
  background: rgba(255, 255, 255, 0.7);
  border: 1px solid var(--border-1);
  transition: all var(--dur-base) var(--ease);
}
.user-chip:hover {
  background: #ffffff;
  border-color: var(--brand-300);
  box-shadow: var(--shadow-sm);
}

.user-chip__avatar {
  width: 30px;
  height: 30px;
  border-radius: var(--r-full);
  background: linear-gradient(135deg, var(--brand-500) 0%, var(--brand-700) 100%);
  color: #fff;
  display: grid;
  place-items: center;
  font-weight: 600;
  font-size: 13px;
  flex-shrink: 0;
}

.user-chip__meta {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  line-height: 1.2;
}

.user-chip__name {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-1);
}

.user-chip__role {
  font-size: 11px;
  color: var(--text-3);
}

.user-chip__caret {
  color: var(--text-3);
  margin-left: 2px;
}

/* ============ Aside ============ */
.layout__aside {
  background: #ffffff;
  height: calc(100vh - 60px);
  border-right: 1px solid var(--border-1);
  display: flex;
  flex-direction: column;
  padding: var(--sp-3) 0 var(--sp-2);
}

.aside__search {
  padding: 0 var(--sp-3) var(--sp-3);
  border-bottom: 1px solid var(--border-2);
}

.aside__scroll {
  flex: 1;
  padding: var(--sp-2) var(--sp-2);
}

.aside__menu {
  border-right: none;
  background: transparent;
}

.aside__empty {
  padding: var(--sp-6) var(--sp-3);
  text-align: center;
  color: var(--text-3);
  font-size: 13px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--sp-2);
}

/* 菜单项基础 */
.aside__menu :deep(.el-menu-item),
.aside__menu :deep(.el-sub-menu__title) {
  border-radius: var(--r-md);
  height: 40px;
  line-height: 40px;
  margin: 2px var(--sp-1);
  font-size: 14px;
  color: var(--text-2);
  transition: all var(--dur-fast) var(--ease);
  position: relative;
  padding-left: 14px !important;
}

.aside__menu :deep(.el-sub-menu .el-menu-item) {
  padding-left: 42px !important;
  min-width: auto;
}

/* hover 态 */
.aside__menu :deep(.el-menu-item:hover),
.aside__menu :deep(.el-sub-menu__title:hover) {
  background: var(--bg-soft);
  color: var(--brand-700);
}

/* 激活态 */
.aside__menu :deep(.el-menu-item.is-active) {
  background: var(--brand-50);
  color: var(--brand-700);
  font-weight: 600;
}

.aside__menu :deep(.el-menu-item.is-active::before) {
  content: '';
  position: absolute;
  left: 0;
  top: 6px;
  bottom: 6px;
  width: 3px;
  border-radius: 0 3px 3px 0;
  background: var(--brand-600);
}

.aside__menu :deep(.el-menu-item .el-icon),
.aside__menu :deep(.el-sub-menu__title .el-icon) {
  color: inherit;
  margin-right: 10px;
}

/* ============ Main ============ */
.layout__main {
  background: var(--bg-page);
  padding: var(--sp-6);
  height: calc(100vh - 60px);
  overflow-y: auto;
}

/* ============ Mobile ============ */
@media (max-width: 768px) {
  .layout__header { padding: 0 var(--sp-4); }
  .brand__text p { display: none; }
  .brand__text h1 { font-size: 15px; }
  .user-chip__meta { display: none; }
  .user-chip { padding: 4px; }
  .layout__main { padding: var(--sp-4); }
}
</style>
