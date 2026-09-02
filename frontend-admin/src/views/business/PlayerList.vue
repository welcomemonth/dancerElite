<template>
  <div class="player-list">
    <div class="page-header">
      <div>
        <h2 class="page-header__title">选手管理</h2>
        <p class="page-header__desc">查看并维护选手资料（当前年龄组仅影响后续新赛季，不影响历史数据）</p>
      </div>
    </div>

    <div class="filter-bar">
      <el-input v-model="filter.name" placeholder="姓名" clearable style="width: 160px" @keyup.enter="loadPlayers" @clear="loadPlayers" />
      <el-input v-model="filter.institution" placeholder="所属机构" clearable style="width: 180px" @keyup.enter="loadPlayers" @clear="loadPlayers" />
      <el-select v-model="filter.age_group" placeholder="当前年龄组" clearable style="width: 140px" @change="loadPlayers">
        <el-option v-for="g in AGE_GROUP_OPTIONS" :key="g" :label="g" :value="g" />
      </el-select>
      <el-button type="primary" @click="loadPlayers">查询</el-button>
    </div>

    <el-card>
      <el-table :data="players" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="real_name" label="姓名" width="120" show-overflow-tooltip />
        <el-table-column label="性别" width="80">
          <template #default="scope">{{ genderText(scope.row.gender) }}</template>
        </el-table-column>
        <el-table-column prop="institution" label="所属机构" min-width="160" show-overflow-tooltip />
        <el-table-column prop="teacher" label="指导老师" width="110" show-overflow-tooltip />
        <el-table-column prop="phone" label="联系电话" width="130" show-overflow-tooltip />
        <el-table-column label="出生日期" width="120">
          <template #default="scope">
            {{ scope.row.birth_year }}-{{ pad(scope.row.birth_month) }}-{{ pad(scope.row.birth_day) }}
          </template>
        </el-table-column>
        <el-table-column label="当前年龄组" width="110">
          <template #default="scope">
            <el-tag v-if="scope.row.age_group" size="small">{{ scope.row.age_group }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <el-button text type="primary" :icon="Edit" @click="openEdit(scope.row)" v-permission="'player:update'">编辑</el-button>
            <el-button text type="danger" :icon="Delete" @click="deletePlayer(scope.row)" v-permission="'player:delete'">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadPlayers"
        @current-change="loadPlayers"
      />
    </el-card>

    <!-- 编辑选手资料 -->
    <el-dialog v-model="editVisible" title="编辑选手资料" width="520px">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="姓名">
          <el-input v-model="editForm.real_name" />
        </el-form-item>
        <el-form-item label="性别">
          <el-select v-model="editForm.gender" style="width: 100%">
            <el-option label="男" value="male" />
            <el-option label="女" value="female" />
          </el-select>
        </el-form-item>
        <el-form-item label="联系电话">
          <el-input v-model="editForm.phone" />
        </el-form-item>
        <el-form-item label="所属机构">
          <el-input v-model="editForm.institution" />
        </el-form-item>
        <el-form-item label="指导老师">
          <el-input v-model="editForm.teacher" />
        </el-form-item>
        <el-form-item label="当前年龄组">
          <el-select v-model="editForm.age_group" style="width: 100%">
            <el-option v-for="g in AGE_GROUP_OPTIONS" :key="g" :label="g" :value="g" />
          </el-select>
          <div class="form-tip">仅影响之后新赛季的成绩归属，历史赛季数据不变</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="savePlayer">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, Delete } from '@element-plus/icons-vue'
import { playerApi } from '../../api'

const AGE_GROUP_OPTIONS = ['U11', 'U13', 'U15']

const players = ref([])
const loading = ref(false)
const filter = ref({ name: '', institution: '', age_group: '' })
const pagination = ref({ page: 1, page_size: 20, total: 0 })

const editVisible = ref(false)
const saving = ref(false)
const editForm = ref({ id: null, real_name: '', gender: 'male', phone: '', institution: '', teacher: '', age_group: '' })

const genderText = (g) => (g === 'male' ? '男' : g === 'female' ? '女' : g || '-')
const pad = (n) => String(n).padStart(2, '0')

const loadPlayers = async () => {
  loading.value = true
  try {
    const params = { page: pagination.value.page, page_size: pagination.value.page_size }
    if (filter.value.name) params.name = filter.value.name
    if (filter.value.institution) params.institution = filter.value.institution
    if (filter.value.age_group) params.age_group = filter.value.age_group
    const data = await playerApi.list(params)
    players.value = data.list || []
    pagination.value.total = data.total
  } catch (e) {
    ElMessage.error('加载选手失败')
  } finally {
    loading.value = false
  }
}

const openEdit = (row) => {
  editForm.value = {
    id: row.id,
    real_name: row.real_name,
    gender: row.gender,
    phone: row.phone,
    institution: row.institution,
    teacher: row.teacher,
    age_group: row.age_group,
  }
  editVisible.value = true
}

const savePlayer = async () => {
  saving.value = true
  try {
    const { id, ...data } = editForm.value
    await playerApi.update(id, data)
    ElMessage.success('保存成功')
    editVisible.value = false
    loadPlayers()
  } catch (e) {
    // 错误已由拦截器提示
  } finally {
    saving.value = false
  }
}

const deletePlayer = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除选手「${row.real_name}」吗？`, '提示', { type: 'warning' })
  } catch (e) {
    return
  }
  try {
    await playerApi.delete(row.id)
    ElMessage.success('删除成功')
    loadPlayers()
  } catch (e) {
    // 错误已由拦截器提示
  }
}

onMounted(() => loadPlayers())
</script>

<style scoped>
.player-list {
  display: flex;
  flex-direction: column;
}
.filter-bar {
  display: flex;
  gap: var(--sp-2);
  margin-bottom: var(--sp-3);
}
.form-tip {
  margin-top: 4px;
  color: var(--text-3);
  font-size: 12px;
}
</style>
