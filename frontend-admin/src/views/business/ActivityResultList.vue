<template>
  <div class="result-list">
    <div class="page-header">
      <div>
        <h2 class="page-header__title">成绩管理</h2>
        <p class="page-header__desc">已结束活动的成绩查询、编辑与批量导入</p>
      </div>
      <div class="page-header__actions">
        <el-button type="primary" :icon="Upload" @click="openImport" v-permission="'activity_result:create'">
          导入成绩
        </el-button>
      </div>
    </div>

    <div class="filter-bar">
      <el-select v-model="filter.activity_id" placeholder="选择活动" clearable @change="onFilterChange" style="width: 220px">
        <el-option label="全部活动" :value="0" />
        <el-option v-for="act in activities" :key="act.id" :label="act.title" :value="act.id" />
      </el-select>
      <el-select v-model="filter.season_id" placeholder="选择赛季" clearable @change="onFilterChange" style="width: 160px">
        <el-option label="全部赛季" :value="0" />
        <el-option v-for="s in seasons" :key="s.id" :label="s.name" :value="s.id" />
      </el-select>
      <el-select v-model="filter.dance_type" placeholder="选择舞种" clearable @change="onFilterChange" style="width: 150px">
        <el-option v-for="d in DANCE_TYPE_OPTIONS" :key="d" :label="d" :value="d" />
      </el-select>
      <el-select v-model="filter.age_group" placeholder="选择级别" clearable @change="onFilterChange" style="width: 130px">
        <el-option v-for="g in AGE_GROUP_OPTIONS" :key="g" :label="g" :value="g" />
      </el-select>
      <el-button :icon="Refresh" @click="resetFilter">重置</el-button>
    </div>

    <el-card>
      <el-table :data="results" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="activity_title" label="活动名称" show-overflow-tooltip min-width="160" />
        <el-table-column prop="player_name" label="选手" width="120" />
        <el-table-column prop="age_group" label="级别" width="90" />
        <el-table-column prop="dance_type" label="舞种" width="110" />
        <el-table-column prop="rank" label="名次" width="80" />
        <el-table-column prop="points" label="积分" width="80" />
        <el-table-column prop="award" label="奖项" width="110" />
        <el-table-column prop="participant_num" label="参赛人数" width="100" />
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="scope">
            <el-button text type="primary" :icon="Edit" @click="openEdit(scope.row)" v-permission="'activity_result:update'">编辑</el-button>
            <el-button text type="danger" :icon="Delete" @click="deleteResult(scope.row)" v-permission="'activity_result:delete'">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadResults"
        @current-change="loadResults"
      />
    </el-card>

    <!-- 导入成绩对话框 -->
    <el-dialog v-model="importVisible" title="导入成绩" width="520px" @closed="resetImport">
      <el-form label-width="90px">
        <el-form-item label="活动" required>
          <el-select v-model="importForm.activity_id" placeholder="请选择活动" style="width: 100%" @change="onImportActivityChange">
            <el-option v-for="act in activities" :key="act.id" :label="act.title" :value="act.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="级别" required>
          <el-select v-model="importForm.age_group" placeholder="请选择级别" style="width: 100%">
            <el-option v-for="g in (importActivity?.age_groups || [])" :key="g" :label="g" :value="g" />
          </el-select>
        </el-form-item>
        <el-form-item label="舞种" required>
          <el-select v-model="importForm.dance_type" placeholder="请选择舞种" style="width: 100%">
            <el-option v-for="d in (importActivity?.dance_types || [])" :key="d" :label="d" :value="d" />
          </el-select>
        </el-form-item>
        <el-form-item label="成绩文件" required>
          <el-upload
            :auto-upload="false"
            :limit="1"
            accept=".xlsx"
            :on-change="onImportFileChange"
            :on-remove="onImportFileRemove"
            :file-list="importFileList"
          >
            <el-button :icon="Upload">选择 .xlsx 文件</el-button>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="importVisible = false">取消</el-button>
        <el-button type="primary" :loading="importing" @click="submitImport">上传</el-button>
      </template>
    </el-dialog>

    <!-- 编辑成绩对话框 -->
    <el-dialog v-model="editVisible" title="编辑成绩" width="480px" @closed="editingRow = null">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="名次">
          <el-input-number v-model="editForm.rank" :min="1" :controls="false" style="width: 100%" />
        </el-form-item>
        <el-form-item label="积分">
          <el-input-number v-model="editForm.points" :min="0" :controls="false" style="width: 100%" />
        </el-form-item>
        <el-form-item label="奖项">
          <el-input v-model="editForm.award" placeholder="如 冠军 / 金奖" />
        </el-form-item>
        <el-form-item label="参赛人数">
          <el-input-number v-model="editForm.participant_num" :min="0" :controls="false" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveEdit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, Edit, Delete, Refresh } from '@element-plus/icons-vue'
import { activityResultApi, activityApi, seasonApi } from '../../api'

const route = useRoute()

const results = ref([])
const activities = ref([])
const seasons = ref([])
const loading = ref(false)
const saving = ref(false)
const importing = ref(false)

const filter = ref({ activity_id: 0, season_id: 0, dance_type: '', age_group: '' })
const pagination = ref({ page: 1, page_size: 20, total: 0 })

const AGE_GROUP_OPTIONS = ['U11', 'U13', 'U15']
const DANCE_TYPE_OPTIONS = ['古典舞', '民族民间舞']

// ===== 列表加载 =====
const loadActivities = async () => {
  try {
    const data = await activityApi.list({ page: 1, page_size: 100 })
    activities.value = data.list || []
  } catch (error) {
    console.error('加载活动列表失败:', error)
  }
}

const loadSeasons = async () => {
  try {
    seasons.value = await seasonApi.list()
  } catch (error) {
    console.error('加载赛季列表失败:', error)
  }
}

const loadResults = async () => {
  loading.value = true
  try {
    const params = { page: pagination.value.page, page_size: pagination.value.page_size }
    if (filter.value.activity_id > 0) params.activity_id = filter.value.activity_id
    if (filter.value.season_id > 0) params.season_id = filter.value.season_id
    if (filter.value.dance_type) params.dance_type = filter.value.dance_type
    if (filter.value.age_group) params.age_group = filter.value.age_group

    const data = await activityResultApi.list(params)
    results.value = data.list || []
    pagination.value.total = data.total
  } catch (error) {
    ElMessage.error('加载成绩失败')
  } finally {
    loading.value = false
  }
}

const onFilterChange = () => {
  pagination.value.page = 1
  loadResults()
}

const resetFilter = () => {
  filter.value = { activity_id: 0, season_id: 0, dance_type: '', age_group: '' }
  pagination.value.page = 1
  loadResults()
}

// ===== 导入成绩 =====
const importVisible = ref(false)
const importForm = ref({ activity_id: null, age_group: '', dance_type: '' })
const importFileList = ref([])

const importActivity = computed(() => activities.value.find((a) => a.id === importForm.value.activity_id))

const openImport = () => {
  importForm.value = { activity_id: filter.value.activity_id || null, age_group: '', dance_type: '' }
  importFileList.value = []
  importVisible.value = true
}

const onImportActivityChange = (val) => {
  const act = activities.value.find((a) => a.id === val)
  importForm.value.age_group = act?.age_groups?.[0] || ''
  importForm.value.dance_type = act?.dance_types?.[0] || ''
}

const onImportFileChange = (file) => {
  importFileList.value = [file]
}

const onImportFileRemove = () => {
  importFileList.value = []
}

const resetImport = () => {
  importForm.value = { activity_id: null, age_group: '', dance_type: '' }
  importFileList.value = []
  importing.value = false
}

const submitImport = async () => {
  if (!importForm.value.activity_id) {
    ElMessage.warning('请选择活动')
    return
  }
  if (!importForm.value.age_group || !importForm.value.dance_type) {
    ElMessage.warning('请选择级别和舞种')
    return
  }
  const file = importFileList.value[0]?.raw
  if (!file) {
    ElMessage.warning('请选择要上传的成绩文件')
    return
  }

  const formData = new FormData()
  formData.append('activity_id', importForm.value.activity_id)
  formData.append('age_group', importForm.value.age_group)
  formData.append('dance_type', importForm.value.dance_type)
  formData.append('file', file)

  importing.value = true
  try {
    const res = await activityResultApi.import(formData)
    const { total, imported, skipped, errors } = res || {}
    const failed = errors?.length || 0
    ElMessage.success(`导入完成：共 ${total} 条，成功 ${imported} 条，跳过 ${skipped} 条${failed ? `，失败 ${failed} 条` : ''}`)
    if (failed) console.warn('成绩导入错误明细：', errors)
    importVisible.value = false
    loadResults()
  } catch (error) {
    // 错误已在 request.js 拦截器中统一弹出
  } finally {
    importing.value = false
  }
}

// ===== 编辑 / 删除 =====
const editVisible = ref(false)
const editingRow = ref(null)
const editForm = ref({ rank: 1, points: 0, award: '', participant_num: 0 })

const openEdit = (row) => {
  editingRow.value = row
  editForm.value = {
    rank: row.rank,
    points: row.points,
    award: row.award,
    participant_num: row.participant_num
  }
  editVisible.value = true
}

const saveEdit = async () => {
  const row = editingRow.value
  if (!row) return

  saving.value = true
  try {
    await activityResultApi.update(row.id, {
      activity_id: row.activity_id,
      player_id: row.player_id,
      season_id: row.season_id,
      registration_id: row.registration_id,
      dance_type: row.dance_type,
      age_group: row.age_group,
      rank: editForm.value.rank,
      points: editForm.value.points,
      award: editForm.value.award,
      participant_num: editForm.value.participant_num
    })
    ElMessage.success('更新成功')
    editVisible.value = false
    loadResults()
  } catch (error) {
    // 错误已在 request.js 拦截器中统一弹出
  } finally {
    saving.value = false
  }
}

const deleteResult = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要删除选手「${row.player_name}」的这条成绩吗？`, '提示', { type: 'warning' })
  } catch (e) {
    return
  }

  try {
    await activityResultApi.delete(row.id)
    ElMessage.success('删除成功')
    loadResults()
  } catch (error) {
    // 错误已在 request.js 拦截器中统一弹出
  }
}

onMounted(() => {
  // 支持从活动列表「查看成绩」跳转并预筛选
  const activityId = Number(route.query.activity_id || 0)
  if (activityId > 0) filter.value.activity_id = activityId
  loadActivities()
  loadSeasons()
  loadResults()
})
</script>

<style scoped>
.result-list {
  display: flex;
  flex-direction: column;
}
</style>
