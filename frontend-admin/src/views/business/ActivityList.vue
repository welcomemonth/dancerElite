<template>
  <div class="activity-list">
    <div class="page-header">
      <div>
        <h2 class="page-header__title">活动管理</h2>
        <p class="page-header__desc">线下/线上活动发布与报名状态管理</p>
      </div>
      <div class="page-header__actions">
        <el-button type="primary" :icon="Plus" @click="createActivity" v-permission="'activity:create'">
          新增活动
        </el-button>
      </div>
    </div>

    <div class="filter-bar">
      <el-select v-model="filter.status" placeholder="选择状态" @change="loadActivities" style="width: 160px">
        <el-option label="全部状态" :value="-1" />
        <el-option label="草稿" :value="0" />
        <el-option label="报名中" :value="1" />
        <el-option label="报名截止" :value="2" />
        <el-option label="进行中" :value="3" />
        <el-option label="已结束" :value="4" />
      </el-select>
    </div>

    <el-card>
      <el-table :data="activities" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="活动名称" show-overflow-tooltip />
        <el-table-column prop="location" label="地点" width="150" show-overflow-tooltip />
        <el-table-column label="活动时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.start_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="price" label="费用" width="100">
          <template #default="scope">
            <span :style="{ color: scope.row.price > 0 ? 'var(--brand-700)' : 'var(--success)', fontWeight: 600 }">
              {{ scope.row.price > 0 ? `¥${scope.row.price}` : '免费' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="报名人数" width="120">
          <template #default="scope">
            {{ scope.row.reg_count }}{{ scope.row.max_participants > 0 ? `/${scope.row.max_participants}` : '' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="statusType(scope.row.status)">
              {{ statusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="级别 / 舞种" min-width="200">
          <template #default="scope">
            <el-tag v-for="g in (scope.row.age_groups || [])" :key="'g-' + g" size="small" class="tag-item">{{ g }}</el-tag>
            <el-tag v-for="d in (scope.row.dance_types || [])" :key="'d-' + d" size="small" type="success" class="tag-item">{{ d }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="scope">
            <el-button text type="primary" :icon="Edit" @click="editActivity(scope.row.id)" v-permission="'activity:update'">编辑</el-button>
            <el-button
              text
              type="success"
              :icon="Upload"
              @click="openScoreUpload(scope.row)"
              v-if="scope.row.status === 4"
              v-permission="'activity:update'"
            >
              成绩上传
            </el-button>
            <el-button
              text
              :type="scope.row.status === 1 ? 'warning' : 'success'"
              :icon="scope.row.status === 1 ? VideoPause : VideoPlay"
              @click="toggleStatus(scope.row)"
              v-if="scope.row.status <= 1"
              v-permission="'activity:update_status'"
            >
              {{ scope.row.status === 1 ? '停止报名' : '开放报名' }}
            </el-button>
            <el-button text type="danger" :icon="Delete" @click="deleteActivity(scope.row)" v-permission="'activity:delete'">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadActivities"
        @current-change="loadActivities"
      />
    </el-card>

    <el-dialog v-model="uploadVisible" title="成绩上传" width="520px" @closed="resetUpload">
      <el-form label-width="90px">
        <el-form-item label="活动">
          <span>{{ uploadActivity?.title }}</span>
        </el-form-item>
        <el-form-item label="级别" required>
          <el-select v-model="uploadForm.age_group" placeholder="请选择级别" style="width: 100%">
            <el-option v-for="g in (uploadActivity?.age_groups || [])" :key="g" :label="g" :value="g" />
          </el-select>
        </el-form-item>
        <el-form-item label="舞种" required>
          <el-select v-model="uploadForm.dance_type" placeholder="请选择舞种" style="width: 100%">
            <el-option v-for="d in (uploadActivity?.dance_types || [])" :key="d" :label="d" :value="d" />
          </el-select>
        </el-form-item>
        <el-form-item label="成绩文件" required>
          <el-upload
            :auto-upload="false"
            :limit="1"
            accept=".xlsx"
            :on-change="onFileChange"
            :on-remove="onFileRemove"
            :file-list="fileList"
          >
            <el-button :icon="Upload">选择 .xlsx 文件</el-button>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="uploadVisible = false">取消</el-button>
        <el-button type="primary" :loading="uploading" @click="submitScoreUpload">上传</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, VideoPlay, VideoPause, Upload } from '@element-plus/icons-vue'
import { activityApi, activityResultApi } from '../../api'

const router = useRouter()
const activities = ref([])
const loading = ref(false)

const filter = ref({ status: -1 })
const pagination = ref({ page: 1, page_size: 20, total: 0 })

const statusMap = {
  0: { text: '草稿', type: 'info' },
  1: { text: '报名中', type: 'success' },
  2: { text: '报名截止', type: 'warning' },
  3: { text: '进行中', type: 'primary' },
  4: { text: '已结束', type: 'danger' },
}

const statusText = (s) => statusMap[s]?.text || '未知'
const statusType = (s) => statusMap[s]?.type || 'info'

const loadActivities = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.page_size,
    }
    if (filter.value.status >= 0) {
      params.status = filter.value.status
    }
    const data = await activityApi.list(params)
    activities.value = data.list || []
    pagination.value.total = data.total
  } catch (error) {
    ElMessage.error('加载活动失败')
  } finally {
    loading.value = false
  }
}

const createActivity = () => router.push('/admin/activities/create')
const editActivity = (id) => router.push(`/admin/activities/edit/${id}`)

const toggleStatus = async (row) => {
  const newStatus = row.status === 1 ? 2 : 1
  const action = newStatus === 1 ? '开放报名' : '停止报名'
  try {
    await ElMessageBox.confirm(`确定要${action}吗？`, '提示', { type: 'warning' })
    await activityApi.updateStatus(row.id, { status: newStatus })
    ElMessage.success(`${action}成功`)
    loadActivities()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error(`${action}失败`)
  }
}

const deleteActivity = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该活动吗？删除后不可恢复。', '提示', { type: 'warning' })
  } catch (e) {
    return // 用户取消
  }

  try {
    await activityApi.delete(row.id)
    ElMessage.success('删除成功')
    loadActivities()
  } catch (error) {
    // 后端具体错误（如“活动存在报名记录，无法删除”）已由请求拦截器统一弹出
  }
}

// ===== 成绩上传 =====
const uploadVisible = ref(false)
const uploadActivity = ref(null)
const uploadForm = ref({ age_group: '', dance_type: '' })
const fileList = ref([])
const uploading = ref(false)

const openScoreUpload = (row) => {
  uploadActivity.value = row
  uploadForm.value = {
    age_group: row.age_groups?.[0] || '',
    dance_type: row.dance_types?.[0] || ''
  }
  fileList.value = []
  uploadVisible.value = true
}

const onFileChange = (file) => {
  fileList.value = [file]
}

const onFileRemove = () => {
  fileList.value = []
}

const resetUpload = () => {
  uploadActivity.value = null
  uploadForm.value = { age_group: '', dance_type: '' }
  fileList.value = []
  uploading.value = false
}

const submitScoreUpload = async () => {
  if (!uploadForm.value.age_group || !uploadForm.value.dance_type) {
    ElMessage.warning('请选择级别和舞种')
    return
  }
  const file = fileList.value[0]?.raw
  if (!file) {
    ElMessage.warning('请选择要上传的成绩文件')
    return
  }

  const formData = new FormData()
  formData.append('activity_id', uploadActivity.value.id)
  formData.append('age_group', uploadForm.value.age_group)
  formData.append('dance_type', uploadForm.value.dance_type)
  formData.append('file', file)

  uploading.value = true
  try {
    const res = await activityResultApi.import(formData)
    const { total, imported, skipped, errors } = res || {}
    const failed = errors?.length || 0
    ElMessage.success(`导入完成：共 ${total} 条，成功 ${imported} 条，跳过 ${skipped} 条${failed ? `，失败 ${failed} 条` : ''}`)
    if (failed) {
      console.warn('成绩导入错误明细：', errors)
    }
    uploadVisible.value = false
    loadActivities()
  } catch (error) {
    // 错误已在 request.js 拦截器中统一弹出
  } finally {
    uploading.value = false
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

onMounted(() => loadActivities())
</script>

<style scoped>
.activity-list {
  display: flex;
  flex-direction: column;
}
.tag-item {
  margin-right: 4px;
}
</style>
