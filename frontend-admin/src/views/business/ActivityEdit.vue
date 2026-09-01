<template>
  <div class="activity-edit">
    <div class="page-header">
      <div class="page-header__back">
        <el-button text :icon="ArrowLeft" @click="$router.back()">返回</el-button>
        <div>
          <h2 class="page-header__title">{{ isEdit ? '编辑活动' : '新增活动' }}</h2>
          <p class="page-header__desc">{{ isEdit ? '修改活动信息' : '发布一个新活动并开放报名' }}</p>
        </div>
      </div>
      <div class="page-header__actions">
        <el-button :icon="DocumentChecked" @click="saveActivity(0)" v-permission="isEdit ? 'activity:update' : 'activity:create'">保存草稿</el-button>
        <el-button type="primary" :icon="Promotion" @click="saveActivity(1)" v-permission="isEdit ? 'activity:update' : 'activity:create'">保存并开放报名</el-button>
      </div>
    </div>

    <el-card>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px" class="activity-form">
        <el-form-item label="活动名称" prop="title">
          <el-input v-model="form.title" placeholder="请输入活动名称" />
        </el-form-item>

        <el-form-item label="活动简介" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入活动简介" />
        </el-form-item>

        <el-form-item label="缩略图">
          <div class="thumbnail-upload">
            <el-upload
              :action="uploadUrl"
              :headers="uploadHeaders"
              :show-file-list="false"
              :on-success="handleThumbnailSuccess"
              accept="image/*"
            >
              <img v-if="form.thumbnail" :src="form.thumbnail" class="thumbnail" />
              <el-icon v-else class="upload-icon"><Plus /></el-icon>
            </el-upload>
            <el-button v-if="form.thumbnail" type="danger" size="small" @click="form.thumbnail = ''" v-permission="isEdit ? 'activity:update' : 'activity:create'">删除</el-button>
          </div>
        </el-form-item>

        <el-form-item label="活动地点">
          <el-input v-model="form.location" placeholder="请输入活动地点" />
        </el-form-item>

        <el-form-item label="级别组合" prop="age_groups">
          <el-select v-model="form.age_groups" multiple collapse-tags placeholder="可多选，U11/U13/U15（可扩展）" style="width: 100%">
            <el-option v-for="g in AGE_GROUP_OPTIONS" :key="g" :label="g" :value="g" />
          </el-select>
        </el-form-item>

        <el-form-item label="舞种类型" prop="dance_types">
          <el-select v-model="form.dance_types" multiple collapse-tags placeholder="可多选，古典舞/民族民间舞（可扩展）" style="width: 100%">
            <el-option v-for="d in DANCE_TYPE_OPTIONS" :key="d" :label="d" :value="d" />
          </el-select>
          <span class="form-tip">{{ sessionCount }} 场赛事（级别数 × 舞种数）</span>
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="活动开始时间" prop="start_time">
              <el-date-picker
                v-model="form.start_time"
                type="datetime"
                placeholder="选择开始时间"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="活动结束时间" prop="end_time">
              <el-date-picker
                v-model="form.end_time"
                type="datetime"
                placeholder="选择结束时间"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="报名开始时间">
              <el-date-picker
                v-model="form.reg_start_time"
                type="datetime"
                placeholder="选择报名开始时间"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="报名截止时间">
              <el-date-picker
                v-model="form.reg_end_time"
                type="datetime"
                placeholder="选择报名截止时间"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="人数上限">
              <el-input-number v-model="form.max_participants" :min="0" placeholder="0表示不限" />
              <span class="form-tip">0表示不限制人数</span>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="活动费用">
              <el-input-number v-model="form.price" :min="0" :precision="2" :step="0.01" />
              <span class="form-tip">0表示免费</span>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="活动详情">
          <RichEditor v-model="form.content" :height="400" />
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus, ArrowLeft, DocumentChecked, Promotion } from '@element-plus/icons-vue'
import RichEditor from '../../components/RichEditor.vue'
import { activityApi } from '../../api'

const router = useRouter()
const route = useRoute()
const formRef = ref()

// 级别/舞种选项，后续可在此扩展
const AGE_GROUP_OPTIONS = ['U11', 'U13', 'U15']
const DANCE_TYPE_OPTIONS = ['古典舞', '民族民间舞']

const isEdit = computed(() => !!route.params.id)

const form = ref({
  title: '',
  description: '',
  content: '',
  thumbnail: '',
  location: '',
  start_time: '',
  end_time: '',
  reg_start_time: null,
  reg_end_time: null,
  max_participants: 0,
  price: 0,
  status: 0,
  age_groups: [],
  dance_types: []
})

// 场次数量 = 级别数 × 舞种数（笛卡尔积）
const sessionCount = computed(() => {
  const groups = form.value.age_groups || []
  const types = form.value.dance_types || []
  return groups.length * types.length
})

const rules = {
  title: [{ required: true, message: '请输入活动名称', trigger: 'blur' }],
  start_time: [{ required: true, message: '请选择活动开始时间', trigger: 'change' }],
  end_time: [{ required: true, message: '请选择活动结束时间', trigger: 'change' }],
}

const uploadUrl = (import.meta.env.VITE_API_BASE_URL || '') + '/api/admin/upload/image'
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${localStorage.getItem('token') || ''}`
}))

const handleThumbnailSuccess = (response) => {
  form.value.thumbnail = response.url
}

const loadActivity = async () => {
  if (!isEdit.value) return
  try {
    const data = await activityApi.get(route.params.id)
    form.value = {
      title: data.title,
      description: data.description,
      content: data.content,
      thumbnail: data.thumbnail,
      location: data.location,
      start_time: data.start_time,
      end_time: data.end_time,
      reg_start_time: data.reg_start_time,
      reg_end_time: data.reg_end_time,
      max_participants: data.max_participants,
      price: data.price,
      status: data.status,
      age_groups: data.age_groups || [],
      dance_types: data.dance_types || [],
    }
  } catch (error) {
    ElMessage.error('加载活动失败')
    router.back()
  }
}

const saveActivity = async (status) => {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  const data = { ...form.value, status }

  try {
    if (isEdit.value) {
      await activityApi.update(route.params.id, data)
      ElMessage.success('更新成功')
    } else {
      await activityApi.create(data)
      ElMessage.success('创建成功')
    }
    router.push('/admin/activities')
  } catch (error) {
    // 错误已在 request.js 拦截器中处理
  }
}

onMounted(() => loadActivity())
</script>

<style scoped>
.activity-edit {
  display: flex;
  flex-direction: column;
}
.page-header__back {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
}
.activity-form {
  max-width: 1000px;
}
.thumbnail-upload {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
}
.thumbnail {
  width: 200px;
  height: 150px;
  object-fit: cover;
  border: 1px solid var(--border-1);
  border-radius: var(--r-md);
}
.upload-icon {
  width: 200px;
  height: 150px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px dashed var(--border-1);
  border-radius: var(--r-md);
  cursor: pointer;
  font-size: 40px;
  color: var(--text-3);
  background: var(--bg-soft);
  transition: all var(--dur-base) var(--ease);
}
.upload-icon:hover {
  border-color: var(--brand-400);
  color: var(--brand-600);
  background: var(--brand-50);
}
.form-tip {
  margin-left: 10px;
  color: var(--text-3);
  font-size: 12px;
}
</style>
