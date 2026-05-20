<template>
  <div class="operation-log-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>操作日志</span>
          <el-button @click="fetchLogs">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>

      <el-form :model="query" inline class="filter-bar">
        <el-form-item label="用户">
          <el-input v-model="query.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="模块">
          <el-input v-model="query.module" placeholder="如 articles" clearable />
        </el-form-item>
        <el-form-item label="操作">
          <el-select v-model="query.action" placeholder="全部操作" clearable style="width: 140px">
            <el-option label="新增" value="create" />
            <el-option label="更新" value="update" />
            <el-option label="删除" value="delete" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="logs" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="90" />
        <el-table-column prop="username" label="用户" min-width="120">
          <template #default="{ row }">{{ row.username || '-' }}</template>
        </el-table-column>
        <el-table-column prop="module" label="模块" min-width="130" />
        <el-table-column prop="action" label="操作" width="100">
          <template #default="{ row }">
            <el-tag :type="actionTag(row.action)">{{ actionText(row.action) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP" min-width="140" />
        <el-table-column prop="created_at" label="时间" min-width="180">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="详情" min-width="220">
          <template #default="{ row }">
            <span class="detail-text">{{ formatDetail(row.detail) }}</span>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="query.page"
        v-model:page-size="query.page_size"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchLogs"
        @current-change="fetchLogs"
      />
    </el-card>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, Search } from '@element-plus/icons-vue'
import { operationLogApi } from '../../api'

const loading = ref(false)
const logs = ref([])
const total = ref(0)
const query = ref({
  page: 1,
  page_size: 20,
  username: '',
  module: '',
  action: ''
})

const fetchLogs = async () => {
  loading.value = true
  try {
    const data = await operationLogApi.list(query.value)
    logs.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    ElMessage.error('加载操作日志失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  query.value.page = 1
  fetchLogs()
}

const resetSearch = () => {
  query.value = {
    page: 1,
    page_size: 20,
    username: '',
    module: '',
    action: ''
  }
  fetchLogs()
}

const actionText = (action) => {
  const map = {
    create: '新增',
    update: '更新',
    delete: '删除'
  }
  return map[action] || action || '-'
}

const actionTag = (action) => {
  const map = {
    create: 'success',
    update: 'warning',
    delete: 'danger'
  }
  return map[action] || 'info'
}

const formatDate = (value) => {
  if (!value) return '-'
  return new Date(value).toLocaleString('zh-CN', { hour12: false })
}

const formatDetail = (detail) => {
  if (!detail) return '-'
  if (typeof detail === 'string') return detail
  return `${detail.method || ''} ${detail.path || ''}`.trim() || JSON.stringify(detail)
}

onMounted(fetchLogs)
</script>

<style scoped>
.operation-log-page {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.filter-bar {
  margin-bottom: 12px;
}

.detail-text {
  color: #64748b;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", monospace;
  font-size: 12px;
}
</style>
