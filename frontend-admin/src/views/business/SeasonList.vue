<template>
  <div class="season-list">
    <div class="page-header">
      <div>
        <h2 class="page-header__title">赛季管理</h2>
        <p class="page-header__desc">年度赛季维护，用于积分排行与赛事归档</p>
      </div>
      <div class="page-header__actions">
        <el-button type="primary" :icon="Plus" @click="openCreate" v-permission="'season:create'">
          新增赛季
        </el-button>
      </div>
    </div>

    <el-card>
      <el-table :data="seasons" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="year" label="年份" width="100" />
        <el-table-column prop="name" label="赛季名称" show-overflow-tooltip />
        <el-table-column label="状态" width="110">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'info'">
              {{ scope.row.status === 'active' ? '进行中' : '已归档' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="开始日期" width="130">
          <template #default="scope">{{ formatDate(scope.row.start_date) }}</template>
        </el-table-column>
        <el-table-column label="结束日期" width="130">
          <template #default="scope">{{ formatDate(scope.row.end_date) }}</template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">{{ formatDateTime(scope.row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="scope">
            <el-button text type="primary" :icon="Edit" @click="openEdit(scope.row)" v-permission="'season:update'">编辑</el-button>
            <el-button
              text
              :type="scope.row.status === 'active' ? 'warning' : 'success'"
              @click="toggleStatus(scope.row)"
              v-permission="'season:update_status'"
            >
              {{ scope.row.status === 'active' ? '归档' : '激活' }}
            </el-button>
            <el-button text type="danger" :icon="Delete" @click="deleteSeason(scope.row)" v-permission="'season:delete'">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="showDialog" :title="editing ? '编辑赛季' : '新增赛季'" width="480px" @closed="editing = null">
      <el-form :model="form" label-width="90px">
        <el-form-item label="年份" required>
          <el-input-number v-model="form.year" :min="2000" :max="2100" :controls="false" style="width: 100%" placeholder="如 2026" />
        </el-form-item>
        <el-form-item label="赛季名称" required>
          <el-input v-model="form.name" placeholder="如 2026赛季" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status" style="width: 100%">
            <el-option label="进行中" value="active" />
            <el-option label="已归档" value="archived" />
          </el-select>
        </el-form-item>
        <el-form-item label="开始日期">
          <el-date-picker v-model="form.start_date" type="date" placeholder="选择开始日期" style="width: 100%" />
        </el-form-item>
        <el-form-item label="结束日期">
          <el-date-picker v-model="form.end_date" type="date" placeholder="选择结束日期" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="saveSeason" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import { seasonApi } from '../../api'

const seasons = ref([])
const loading = ref(false)
const saving = ref(false)
const showDialog = ref(false)
const editing = ref(null) // 正在编辑的行，null 表示新增

const form = ref({
  year: new Date().getFullYear(),
  name: '',
  status: 'active',
  start_date: null,
  end_date: null
})

// 获取赛季列表
const loadSeasons = async () => {
  loading.value = true
  try {
    seasons.value = await seasonApi.list()
  } catch (error) {
    ElMessage.error('加载赛季失败')
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  editing.value = null
  form.value = {
    year: new Date().getFullYear(),
    name: '',
    status: 'active',
    start_date: null,
    end_date: null
  }
  showDialog.value = true
}

const openEdit = (row) => {
  editing.value = row
  form.value = {
    year: row.year,
    name: row.name,
    status: row.status,
    start_date: row.start_date ? new Date(row.start_date) : null,
    end_date: row.end_date ? new Date(row.end_date) : null
  }
  showDialog.value = true
}

// 保存（新增/编辑）
const saveSeason = async () => {
  if (!form.value.year) {
    ElMessage.warning('请输入年份')
    return
  }
  if (!form.value.name) {
    ElMessage.warning('请输入赛季名称')
    return
  }

  saving.value = true
  try {
    if (editing.value) {
      await seasonApi.update(editing.value.id, form.value)
      ElMessage.success('更新成功')
    } else {
      await seasonApi.create(form.value)
      ElMessage.success('创建成功')
    }
    showDialog.value = false
    loadSeasons()
  } catch (error) {
    // 错误已在 request.js 拦截器中统一提示
  } finally {
    saving.value = false
  }
}

// 激活 / 归档切换
const toggleStatus = async (row) => {
  const target = row.status === 'active' ? 'archived' : 'active'
  const actionText = target === 'active' ? '激活' : '归档'
  try {
    await ElMessageBox.confirm(`确定要${actionText}「${row.name}」吗？`, '提示', { type: 'warning' })
    await seasonApi.updateStatus(row.id, { status: target })
    ElMessage.success(`${actionText}成功`)
    loadSeasons()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error(`${actionText}失败`)
  }
}

// 删除赛季
const deleteSeason = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要删除「${row.name}」吗？`, '提示', { type: 'warning' })
    await seasonApi.delete(row.id)
    ElMessage.success('删除成功')
    loadSeasons()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

const formatDateTime = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

onMounted(() => loadSeasons())
</script>

<style scoped>
.season-list {
  display: flex;
  flex-direction: column;
}
</style>
