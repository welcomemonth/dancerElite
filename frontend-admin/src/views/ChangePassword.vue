<template>
  <div class="change-password">
    <div class="cp-card">
      <div class="cp-card__head">
        <div class="cp-icon">
          <el-icon><Key /></el-icon>
        </div>
        <div>
          <h2>修改密码</h2>
          <p>定期更换密码可以提升账号安全性</p>
        </div>
      </div>

      <el-form :model="form" :rules="rules" ref="formRef" label-position="top" size="large">
        <el-form-item label="原密码" prop="oldPassword">
          <el-input
            v-model="form.oldPassword"
            type="password"
            placeholder="请输入原密码"
            show-password
            :prefix-icon="Lock"
          />
        </el-form-item>

        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="form.newPassword"
            type="password"
            placeholder="请输入新密码（至少6位）"
            show-password
            :prefix-icon="Lock"
          />
        </el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="form.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            show-password
            :prefix-icon="Lock"
          />
        </el-form-item>

        <div class="cp-actions">
          <el-button @click="resetForm">重置</el-button>
          <el-button type="primary" @click="changePassword" :loading="loading">
            修改密码
          </el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { Key, Lock } from '@element-plus/icons-vue'
import { useUserStore } from '../store/user'
import { usePermissionStore } from '../store/permission'
import { resetDynamicRoutes } from '../router'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const formRef = ref()

const form = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== form.value.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  oldPassword: [
    { required: true, message: '请输入原密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const changePassword = async () => {
  if (!formRef.value) return
  
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await userStore.changePassword({
      old_password: form.value.oldPassword,
      new_password: form.value.newPassword
    })

    ElMessage.success('密码修改成功，请重新登录')

    // 清除登录状态
    userStore.logout()
    usePermissionStore().reset()
    resetDynamicRoutes()
    router.push('/login')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '修改密码失败')
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  form.value = {
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  }
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}
</script>

<style scoped>
.change-password {
  display: grid;
  place-items: start center;
  padding-top: var(--sp-6);
}

.cp-card {
  width: 100%;
  max-width: 460px;
  background: var(--bg-card);
  border: 1px solid var(--border-1);
  border-radius: var(--r-lg);
  box-shadow: var(--shadow-sm);
  padding: var(--sp-6);
}

.cp-card__head {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
  margin-bottom: var(--sp-5);
  padding-bottom: var(--sp-4);
  border-bottom: 1px solid var(--border-2);
}

.cp-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--r-md);
  background: linear-gradient(135deg, var(--brand-600), var(--brand-800));
  color: #fff;
  display: grid;
  place-items: center;
  font-size: 20px;
  box-shadow: var(--shadow-brand);
  flex-shrink: 0;
}

.cp-card__head h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  color: var(--text-1);
}

.cp-card__head p {
  margin: 4px 0 0;
  font-size: 12px;
  color: var(--text-3);
}

.cp-actions {
  display: flex;
  gap: var(--sp-2);
  justify-content: flex-end;
  margin-top: var(--sp-2);
}

.cp-actions .el-button {
  min-width: 100px;
}
</style>
