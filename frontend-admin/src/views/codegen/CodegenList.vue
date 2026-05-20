<template>
  <div class="codegen-list">
    <div class="page-header">
      <div>
        <h2 class="page-header__title">代码生成</h2>
        <p class="page-header__desc">基于数据库表结构一键生成 Go 后端 CRUD 代码</p>
      </div>
      <div class="page-header__actions">
        <el-button type="primary" :icon="Plus" @click="$router.push('/admin/codegen/create')" v-permission="'codegen:create'">
          新建配置
        </el-button>
      </div>
    </div>

    <el-card>
      <el-table :data="list" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="table_name" label="数据库表" width="180" />
        <el-table-column prop="module_name" label="模块名" width="150" />
        <el-table-column prop="display_name" label="显示名称" width="150" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.generated ? 'success' : 'info'">
              {{ row.generated ? '已生成' : '未生成' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="生成时间" width="180">
          <template #default="{ row }">
            {{ row.generated_at ? formatDate(row.generated_at) : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="320">
          <template #default="{ row }">
            <el-button text type="primary" :icon="Edit" @click="$router.push(`/admin/codegen/edit/${row.id}`)" v-permission="'codegen:update'">编辑</el-button>
            <el-button text type="warning" :icon="View" @click="handlePreview(row)" v-permission="'codegen:preview'">预览</el-button>
            <el-button text type="success" :icon="MagicStick" @click="handleGenerate(row)" v-permission="'codegen:update_generate'">
              {{ row.generated ? '重新生成' : '生成' }}
            </el-button>
            <el-button text type="danger" :icon="Delete" @click="handleDelete(row)" v-permission="'codegen:delete'">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @current-change="loadList"
        @size-change="loadList"
      />
    </el-card>

    <!-- 代码预览弹窗 -->
    <el-dialog v-model="previewVisible" title="代码预览" width="80%" top="5vh">
      <el-tabs v-model="previewTab">
        <el-tab-pane label="Model" name="model">
          <pre class="code-block">{{ previewCode.model_code }}</pre>
        </el-tab-pane>
        <el-tab-pane label="Service" name="service">
          <pre class="code-block">{{ previewCode.service_code }}</pre>
        </el-tab-pane>
        <el-tab-pane label="Handler" name="handler">
          <pre class="code-block">{{ previewCode.handler_code }}</pre>
        </el-tab-pane>
        <el-tab-pane label="Router 代码片段" name="router">
          <pre class="code-block">{{ previewCode.router_code }}</pre>
        </el-tab-pane>
        <el-tab-pane label="API 代码片段" name="api">
          <pre class="code-block">{{ previewCode.api_code }}</pre>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, View, MagicStick, Delete } from '@element-plus/icons-vue'
import { codegenApi } from '../../api'

const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)

const previewVisible = ref(false)
const previewTab = ref('model')
const previewCode = ref({})

const loadList = async () => {
  loading.value = true
  try {
    const data = await codegenApi.list({ page: page.value, page_size: pageSize.value })
    list.value = data.list || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

const handlePreview = async (row) => {
  try {
    const data = await codegenApi.preview(row.id)
    previewCode.value = data
    previewTab.value = 'model'
    previewVisible.value = true
  } catch (e) {
    ElMessage.error('预览失败: ' + (e.message || '未知错误'))
  }
}

const handleGenerate = async (row) => {
  const action = row.generated ? '重新生成（将覆盖已有文件）' : '生成代码'
  try {
    await ElMessageBox.confirm(
      `确定要${action}吗？生成后需要重新编译后端并重启服务。`,
      '确认生成',
      { type: 'warning' }
    )
    const data = await codegenApi.generate(row.id)
    previewCode.value = data
    previewTab.value = 'router'
    previewVisible.value = true
    ElMessage.success('代码生成成功！请将 Router 和 API 代码片段手动添加到对应文件中，然后重新编译。')
    loadList()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('生成失败: ' + (e.message || '未知错误'))
    }
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该配置吗？', '确认删除', { type: 'warning' })
    await codegenApi.delete(row.id)
    ElMessage.success('删除成功')
    loadList()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const formatDate = (str) => {
  if (!str) return ''
  return new Date(str).toLocaleString('zh-CN')
}

onMounted(loadList)
</script>

<style scoped>
.codegen-list {
  display: flex;
  flex-direction: column;
}
.code-block {
  background: #0f172a;
  color: #e2e8f0;
  padding: 16px;
  border-radius: var(--r-md);
  overflow-x: auto;
  font-size: 13px;
  line-height: 1.6;
  max-height: 60vh;
  white-space: pre;
  font-family: ui-monospace, 'Consolas', 'Monaco', monospace;
  margin: 0;
}
</style>
