<template>
  <div class="login">
    <!-- 左侧品牌展示 -->
    <aside class="login__brand">
      <div class="brand-glow brand-glow--1" />
      <div class="brand-glow brand-glow--2" />
      <div class="brand-glow brand-glow--3" />

      <div class="brand-content">
        <div class="brand-mark">
          <span>远</span>
        </div>
        <h1 class="brand-title">远山平台</h1>
        <p class="brand-subtitle">一站式内容、活动与用户运营工作台</p>

        <ul class="brand-features">
          <li>
            <el-icon><DataAnalysis /></el-icon>
            <div>
              <h3>数据洞察</h3>
              <p>多维度运营数据实时可视</p>
            </div>
          </li>
          <li>
            <el-icon><Promotion /></el-icon>
            <div>
              <h3>高效协同</h3>
              <p>菜单 / 角色 / 用户精细化权限</p>
            </div>
          </li>
          <li>
            <el-icon><MagicStick /></el-icon>
            <div>
              <h3>代码生成</h3>
              <p>表单到 CRUD 一键生成业务模块</p>
            </div>
          </li>
        </ul>
      </div>

      <div class="brand-footer">
        © {{ new Date().getFullYear() }} 远山平台 · 让运营更简单
      </div>
    </aside>

    <!-- 右侧登录表单 -->
    <main class="login__panel">
      <div class="login-card">
        <div class="login-card__head">
          <div class="brand-mark brand-mark--sm">
            <span>远</span>
          </div>
          <div>
            <h2>欢迎回来</h2>
            <p>请使用管理员账号登录</p>
          </div>
        </div>

        <el-form
          :model="form"
          label-position="top"
          @submit.prevent="onSubmit"
          class="login-form"
        >
          <el-form-item label="用户名">
            <el-input
              v-model="form.username"
              placeholder="请输入用户名"
              size="large"
              :prefix-icon="User"
            />
          </el-form-item>

          <el-form-item label="密码">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="请输入密码"
              show-password
              size="large"
              :prefix-icon="Lock"
              @keyup.enter="onSubmit"
            />
          </el-form-item>

          <el-button
            type="primary"
            native-type="submit"
            :loading="loading"
            size="large"
            class="login-btn"
            @click="onSubmit"
          >
            登 录
          </el-button>
        </el-form>

        <p class="login-tip">
          <el-icon><InfoFilled /></el-icon>
          忘记密码请联系系统管理员
        </p>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  User, Lock, InfoFilled,
  DataAnalysis, Promotion, MagicStick
} from '@element-plus/icons-vue'
import { useUserStore } from '../store/user'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const form = ref({
  username: '',
  password: ''
})

const onSubmit = async () => {
  if (!form.value.username || !form.value.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }

  loading.value = true
  try {
    await userStore.login(form.value)
    ElMessage.success('登录成功')
    router.push('/admin/articles')
  } catch (error) {
    // 错误已在 request.js 拦截器中处理
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login {
  height: 100vh;
  display: flex;
  background: var(--bg-page);
  overflow: hidden;
}

/* ============ 左侧品牌区 ============ */
.login__brand {
  flex: 1.2;
  position: relative;
  background: linear-gradient(135deg, var(--brand-700) 0%, var(--brand-900) 50%, #0c2870 100%);
  color: #fff;
  padding: 56px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  overflow: hidden;
}

.brand-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  pointer-events: none;
}
.brand-glow--1 {
  width: 380px; height: 380px;
  background: rgba(96, 165, 250, 0.30);
  top: -120px; right: -100px;
}
.brand-glow--2 {
  width: 320px; height: 320px;
  background: rgba(59, 130, 246, 0.25);
  bottom: -100px; left: -80px;
}
.brand-glow--3 {
  width: 220px; height: 220px;
  background: rgba(147, 197, 253, 0.18);
  top: 40%; left: 40%;
}

.brand-content {
  position: relative;
  z-index: 1;
}

.brand-mark {
  width: 56px;
  height: 56px;
  border-radius: var(--r-lg);
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.25);
  display: grid;
  place-items: center;
  color: #fff;
  font-size: 24px;
  font-weight: 700;
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.18);
  margin-bottom: 32px;
}

.brand-title {
  font-size: 38px;
  font-weight: 700;
  line-height: 1.2;
  letter-spacing: -0.5px;
  margin: 0 0 12px;
  color: #ffffff;
}

.brand-subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.78);
  margin: 0 0 48px;
}

.brand-features {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.brand-features li {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.brand-features .el-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--r-md);
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.18);
  font-size: 20px;
  display: grid;
  place-items: center;
  flex-shrink: 0;
  color: #bfdbfe;
}

.brand-features h3 {
  font-size: 15px;
  font-weight: 600;
  margin: 0 0 4px;
  color: #fff;
}

.brand-features p {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.65);
  margin: 0;
  line-height: 1.5;
}

.brand-footer {
  position: relative;
  z-index: 1;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.55);
}

/* ============ 右侧表单区 ============ */
.login__panel {
  flex: 1;
  display: grid;
  place-items: center;
  padding: 40px;
}

.login-card {
  width: 100%;
  max-width: 380px;
}

.login-card__head {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 32px;
}

.brand-mark--sm {
  width: 44px;
  height: 44px;
  border-radius: var(--r-md);
  background: linear-gradient(135deg, var(--brand-600) 0%, var(--brand-800) 100%);
  border: none;
  font-size: 18px;
  box-shadow: var(--shadow-brand);
  margin-bottom: 0;
}

.login-card__head h2 {
  margin: 0;
  font-size: 22px;
  font-weight: 700;
  color: var(--text-1);
  line-height: 1.3;
}

.login-card__head p {
  margin: 4px 0 0;
  font-size: 13px;
  color: var(--text-3);
}

.login-form :deep(.el-form-item__label) {
  font-size: 13px;
  color: var(--text-2);
  padding-bottom: 6px;
}

.login-btn {
  width: 100%;
  height: 46px;
  font-size: 15px;
  font-weight: 600;
  letter-spacing: 2px;
  background: linear-gradient(135deg, var(--brand-600) 0%, var(--brand-700) 100%);
  border: none;
  margin-top: 8px;
  border-radius: var(--r-md);
}
.login-btn:hover {
  background: linear-gradient(135deg, var(--brand-500) 0%, var(--brand-600) 100%);
  box-shadow: var(--shadow-brand);
}

.login-tip {
  margin: 20px 0 0;
  font-size: 12px;
  color: var(--text-3);
  display: flex;
  align-items: center;
  gap: 6px;
  justify-content: center;
}

/* ============ Mobile ============ */
@media (max-width: 900px) {
  .login {
    flex-direction: column;
  }
  .login__brand {
    flex: none;
    padding: 32px 24px;
    min-height: 280px;
  }
  .brand-title {
    font-size: 28px;
  }
  .brand-subtitle {
    margin-bottom: 24px;
  }
  .brand-features {
    display: none;
  }
  .brand-footer {
    display: none;
  }
  .brand-mark {
    margin-bottom: 16px;
  }
  .login__panel {
    flex: 1;
    padding: 24px;
  }
}

@media (max-width: 480px) {
  .login__brand {
    min-height: 200px;
    padding: 24px;
  }
  .brand-title {
    font-size: 22px;
    margin-bottom: 6px;
  }
}
</style>
