<template>
  <div class="login-container">
    <el-card class="login-card">
      <div class="login-head">
        <div class="brand-mark">远</div>
        <h2 class="title">远山公益后台管理系统</h2>
        <p>请使用后台账号登录工作台</p>
      </div>

      <el-form :model="form" label-position="top" @submit.prevent="onSubmit">
        <el-form-item label="用户名">
          <el-input v-model="form.username" placeholder="请输入用户名" size="large" />
        </el-form-item>

        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password size="large" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading" size="large" class="login-button">
            登录
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
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
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 24px;
  box-sizing: border-box;
  background:
    linear-gradient(135deg, rgba(15, 118, 110, 0.10), rgba(245, 158, 11, 0.08)),
    #f5f7f6;
}

.login-card {
  width: min(420px, 100%);
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  box-shadow: 0 18px 45px rgba(17, 24, 39, 0.12);
}

.login-head {
  text-align: center;
  margin-bottom: 26px;
}

.brand-mark {
  width: 44px;
  height: 44px;
  margin: 0 auto 14px;
  border-radius: 8px;
  display: grid;
  place-items: center;
  background: #0f766e;
  color: #ffffff;
  font-size: 20px;
  font-weight: 700;
}

.title {
  margin: 0;
  color: #111827;
  font-size: 22px;
  line-height: 30px;
  letter-spacing: 0;
}

.login-head p {
  margin: 8px 0 0;
  color: #6b7280;
  font-size: 14px;
}

.login-button {
  width: 100%;
  margin-top: 6px;
}
</style>
