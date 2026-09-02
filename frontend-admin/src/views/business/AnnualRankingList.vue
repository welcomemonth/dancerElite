<template>
  <div class="ranking-list">
    <div class="page-header">
      <div>
        <h2 class="page-header__title">年度积分榜</h2>
        <p class="page-header__desc">按「最高 3 场积分之和」自动统计，成绩录入后自动重算</p>
      </div>
      <div class="page-header__actions">
        <el-button type="primary" :icon="Refresh" @click="recalculate" v-permission="'annual_ranking:create'">
          重算积分榜
        </el-button>
      </div>
    </div>

    <div class="filter-bar">
      <el-select v-model="filter.season_id" placeholder="选择赛季" clearable @change="onFilterChange" style="width: 180px">
        <el-option label="全部赛季" :value="0" />
        <el-option v-for="s in seasons" :key="s.id" :label="s.name" :value="s.id" />
      </el-select>
      <el-select v-model="filter.age_group" placeholder="选择级别" clearable @change="onFilterChange" style="width: 130px">
        <el-option v-for="g in AGE_GROUP_OPTIONS" :key="g" :label="g" :value="g" />
      </el-select>
      <el-select v-model="filter.dance_type" placeholder="选择舞种" clearable @change="onFilterChange" style="width: 150px">
        <el-option v-for="d in DANCE_TYPE_OPTIONS" :key="d" :label="d" :value="d" />
      </el-select>
      <el-button :icon="Refresh" @click="resetFilter">重置</el-button>
    </div>

    <el-card>
      <el-table :data="rankings" v-loading="loading">
        <el-table-column prop="rank" label="名次" width="70" />
        <el-table-column prop="player_name" label="选手" width="120" />
        <el-table-column prop="institution" label="机构" show-overflow-tooltip min-width="140" />
        <el-table-column prop="age_group" label="级别" width="80" />
        <el-table-column prop="dance_type" label="舞种" width="110" />
        <el-table-column prop="total_points" label="总积分" width="90" />
        <el-table-column label="直通决赛" width="90">
          <template #default="scope">
            <el-tag v-if="scope.row.is_direct_advance" type="warning" size="small">直通</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="score_count" label="计入场次" width="90" />
        <el-table-column label="排名变化" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.previous_rank === 0" type="info" size="small">新上榜</el-tag>
            <el-tag v-else-if="scope.row.rank_change > 0" type="success" size="small">↑{{ scope.row.rank_change }}</el-tag>
            <el-tag v-else-if="scope.row.rank_change < 0" type="danger" size="small">↓{{ -scope.row.rank_change }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="上次重算" width="180">
          <template #default="scope">{{ formatDate(scope.row.last_updated_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="scope">
            <el-button text type="primary" :icon="Edit" @click="openEdit(scope.row)" v-permission="'annual_ranking:update'">编辑</el-button>
            <el-button text type="danger" :icon="Delete" @click="deleteRanking(scope.row)" v-permission="'annual_ranking:delete'">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadRankings"
        @current-change="loadRankings"
      />
    </el-card>

    <!-- 编辑对话框 -->
    <el-dialog v-model="editVisible" title="编辑积分" width="420px" @closed="editingRow = null">
      <el-form :model="editForm" label-width="90px">
        <el-form-item label="总积分">
          <el-input-number v-model="editForm.total_points" :min="0" :controls="false" style="width: 100%" />
        </el-form-item>
        <el-form-item label="名次">
          <el-input-number v-model="editForm.rank" :min="1" :controls="false" style="width: 100%" />
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
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, Delete, Refresh } from '@element-plus/icons-vue'
import { annualRankingApi, seasonApi } from '../../api'

const rankings = ref([])
const seasons = ref([])
const loading = ref(false)
const saving = ref(false)

const filter = ref({ season_id: 0, age_group: '', dance_type: '' })
const pagination = ref({ page: 1, page_size: 20, total: 0 })

const AGE_GROUP_OPTIONS = ['U11', 'U13', 'U15']
const DANCE_TYPE_OPTIONS = ['古典舞', '民族民间舞']

const loadSeasons = async () => {
  try {
    seasons.value = await seasonApi.list()
  } catch (error) {
    console.error('加载赛季列表失败:', error)
  }
}

const loadRankings = async () => {
  loading.value = true
  try {
    const params = { page: pagination.value.page, page_size: pagination.value.page_size }
    if (filter.value.season_id > 0) params.season_id = filter.value.season_id
    if (filter.value.age_group) params.age_group = filter.value.age_group
    if (filter.value.dance_type) params.dance_type = filter.value.dance_type

    const data = await annualRankingApi.list(params)
    rankings.value = data.list || []
    pagination.value.total = data.total
  } catch (error) {
    ElMessage.error('加载积分榜失败')
  } finally {
    loading.value = false
  }
}

const onFilterChange = () => {
  pagination.value.page = 1
  loadRankings()
}

const resetFilter = () => {
  filter.value = { season_id: 0, age_group: '', dance_type: '' }
  pagination.value.page = 1
  loadRankings()
}

// ===== 重算 =====
const recalculate = async () => {
  if (!filter.value.season_id) {
    ElMessage.warning('请先选择赛季')
    return
  }
  const season = seasons.value.find((s) => s.id === filter.value.season_id)
  try {
    await ElMessageBox.confirm(
      `确定要重算「${season?.name || filter.value.season_id}」的积分榜吗？将根据当前所有成绩重新计算。`,
      '提示',
      { type: 'warning' }
    )
  } catch (e) {
    return
  }

  try {
    const res = await annualRankingApi.recalculate({ season_id: filter.value.season_id })
    ElMessage.success(`重算完成，共 ${res.count} 条`)
    loadRankings()
  } catch (error) {
    // 错误已在 request.js 拦截器中统一弹出
  }
}

// ===== 编辑 / 删除 =====
const editVisible = ref(false)
const editingRow = ref(null)
const editForm = ref({ total_points: 0, rank: 1 })

const openEdit = (row) => {
  editingRow.value = row
  editForm.value = { total_points: row.total_points, rank: row.rank }
  editVisible.value = true
}

const saveEdit = async () => {
  const row = editingRow.value
  if (!row) return

  saving.value = true
  try {
    await annualRankingApi.update(row.id, {
      total_points: editForm.value.total_points,
      rank: editForm.value.rank
    })
    ElMessage.success('更新成功')
    editVisible.value = false
    loadRankings()
  } catch (error) {
    // 错误已在 request.js 拦截器中统一弹出
  } finally {
    saving.value = false
  }
}

const deleteRanking = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要删除选手「${row.player_name}」的这条积分榜记录吗？`, '提示', { type: 'warning' })
  } catch (e) {
    return
  }

  try {
    await annualRankingApi.delete(row.id)
    ElMessage.success('删除成功')
    loadRankings()
  } catch (error) {
    // 错误已在 request.js 拦截器中统一弹出
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

onMounted(() => {
  loadSeasons()
  loadRankings()
})
</script>

<style scoped>
.ranking-list {
  display: flex;
  flex-direction: column;
}
</style>
