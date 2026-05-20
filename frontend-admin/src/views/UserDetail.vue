<template>
  <div class="user-detail">
    <div class="page-header">
      <div class="page-header__back">
        <el-button text :icon="ArrowLeft" @click="$router.back()">返回</el-button>
        <div>
          <h2 class="page-header__title">用户详情</h2>
          <p class="page-header__desc">小程序用户的基本信息</p>
        </div>
      </div>
    </div>

    <el-card>
      <el-descriptions v-if="user" :column="1" border>
        <el-descriptions-item label="ID">{{ user.id }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ user.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item label="OpenID">{{ user.openid || '-' }}</el-descriptions-item>
        <el-descriptions-item label="姓名">{{ user.name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ user.created_at }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ user.updated_at }}</el-descriptions-item>
      </el-descriptions>

      <div v-else class="loading">
        加载中...
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { userApi } from '../api'

const route = useRoute()
const user = ref(null)

onMounted(async () => {
  try {
    user.value = await userApi.get(route.params.id)
  } catch (error) {
    console.error('加载用户失败', error)
  }
})
</script>

<style scoped>
.user-detail {
  display: flex;
  flex-direction: column;
}
.page-header__back {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
}
.loading {
  text-align: center;
  padding: var(--sp-10);
  color: var(--text-3);
}
</style>
