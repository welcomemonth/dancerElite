<template>
  <div class="article-edit">
    <div class="page-header">
      <div class="page-header__back">
        <el-button text :icon="ArrowLeft" @click="$router.back()">返回</el-button>
        <div>
          <h2 class="page-header__title">{{ isEdit ? '编辑文章' : '新增文章' }}</h2>
          <p class="page-header__desc">{{ isEdit ? '修改文章内容' : '创建一篇新文章' }}</p>
        </div>
      </div>
      <div class="page-header__actions">
        <el-button :icon="DocumentChecked" @click="saveArticle(0)" v-permission="isEdit ? 'article:update' : 'article:create'">保存草稿</el-button>
        <el-button type="primary" :icon="Promotion" @click="saveArticle(1)" v-permission="isEdit ? 'article:update' : 'article:create'">发布文章</el-button>
      </div>
    </div>

    <el-card>
      <el-form :model="form" label-width="100px" class="article-form">
        <el-form-item label="所属栏目" required>
          <el-select v-model="form.column_id" placeholder="请选择栏目">
            <el-option
              v-for="column in columns"
              :key="column.id"
              :label="column.name"
              :value="column.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="文章标题" required>
          <el-input v-model="form.title" placeholder="请输入文章标题" />
        </el-form-item>

        <el-form-item label="作者">
          <el-input v-model="form.author" placeholder="请输入作者" />
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
            <el-button v-if="form.thumbnail" type="danger" size="small" @click="form.thumbnail = ''" v-permission="isEdit ? 'article:update' : 'article:create'">删除</el-button>
          </div>
        </el-form-item>

        <el-form-item label="文章内容">
          <RichEditor v-model="form.content" :height="500" />
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
import RichEditor from '../components/RichEditor.vue'
import { articleApi, columnApi } from '../api'

const router = useRouter()
const route = useRoute()

const isEdit = computed(() => !!route.params.id)
const columns = ref([])
const form = ref({
  column_id: null,
  title: '',
  author: '',
  thumbnail: '',
  content: '',
  status: 0
})

// 缩略图上传配置
const uploadUrl = (import.meta.env.VITE_API_BASE_URL || '') + '/api/admin/upload/image'
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${localStorage.getItem('token') || ''}`
}))

const loadColumns = async () => {
  try {
    columns.value = await columnApi.list()
  } catch (error) {
    ElMessage.error('加载栏目失败')
  }
}

const loadArticle = async () => {
  if (!isEdit.value) return

  try {
    const data = await articleApi.get(route.params.id)
    form.value = {
      column_id: data.column_id,
      title: data.title,
      author: data.author,
      thumbnail: data.thumbnail,
      content: data.content,
      status: data.status
    }
  } catch (error) {
    ElMessage.error('加载文章失败')
    router.back()
  }
}

const handleThumbnailSuccess = (response) => {
  form.value.thumbnail = response.url
}

const saveArticle = async (status) => {
  if (!form.value.column_id) {
    ElMessage.warning('请选择栏目')
    return
  }
  if (!form.value.title) {
    ElMessage.warning('请输入文章标题')
    return
  }

  const data = { ...form.value, status }

  try {
    if (isEdit.value) {
      await articleApi.update(route.params.id, data)
      ElMessage.success('更新成功')
    } else {
      await articleApi.create(data)
      ElMessage.success('创建成功')
    }
    router.push('/admin/articles')
  } catch (error) {
    // 错误已在 request.js 拦截器中处理
  }
}

onMounted(() => {
  loadColumns()
  loadArticle()
})
</script>

<style scoped>
.article-edit {
  display: flex;
  flex-direction: column;
}
.page-header__back {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
}
.article-form {
  max-width: 1200px;
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
</style>
