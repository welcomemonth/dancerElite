<template>
  <div class="menu-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>菜单管理</span>
          <div class="header-actions">
            <el-button @click="fetchMenus">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button type="primary" @click="openCreateDialog" v-permission="'menu:create'">
              <el-icon><Plus /></el-icon>
              新增菜单
            </el-button>
          </div>
        </div>
      </template>

      <div class="menu-workspace">
        <aside class="menu-tree-panel">
          <div class="tree-title">菜单树</div>
          <el-tree
            :data="treeData"
            node-key="id"
            :props="treeProps"
            :current-node-key="selectedMenuId"
            highlight-current
            default-expand-all
            :expand-on-click-node="false"
            @node-click="handleTreeNodeClick"
          >
            <template #default="{ data }">
              <span class="tree-node">
                <el-icon v-if="data.icon"><component :is="data.icon" /></el-icon>
                <span>{{ data.title }}</span>
                <el-tag v-if="data.id !== 0" size="small" :type="typeTag(data.type)">
                  {{ typeText(data.type) }}
                </el-tag>
              </span>
            </template>
          </el-tree>
        </aside>

        <section class="menu-list-panel">
          <div class="list-toolbar">
            <div>
              <div class="list-title">{{ selectedMenuTitle }}</div>
              <div class="list-subtitle">{{ tableMenus.length }} 条数据</div>
            </div>
            <el-button v-if="selectedMenuId !== 0" @click="selectedMenuId = 0">查看全部</el-button>
          </div>

          <el-table
            :data="tableMenus"
            style="width: 100%"
            v-loading="loading"
            row-key="id"
          >
            <el-table-column prop="title" label="菜单名称" min-width="210">
              <template #default="{ row }">
                <span class="menu-title" :style="{ paddingLeft: `${row.level * 18}px` }">
                  <span class="menu-depth">L{{ row.level + 1 }}</span>
                  <el-icon v-if="row.icon"><component :is="row.icon" /></el-icon>
                  <span>{{ row.title }}</span>
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="parent_title" label="父级" min-width="130">
              <template #default="{ row }">{{ row.parent_title || '顶级菜单' }}</template>
            </el-table-column>
            <el-table-column prop="name" label="标识" min-width="150" />
            <el-table-column prop="path" label="路径" min-width="180">
              <template #default="{ row }">{{ row.path || '-' }}</template>
            </el-table-column>
            <el-table-column prop="component" label="组件" min-width="170">
              <template #default="{ row }">{{ row.component || '-' }}</template>
            </el-table-column>
            <el-table-column prop="permission" label="权限标识" min-width="160">
              <template #default="{ row }">{{ row.permission || '-' }}</template>
            </el-table-column>
            <el-table-column prop="sort" label="排序" width="80" />
            <el-table-column prop="type" label="类型" width="100">
              <template #default="{ row }">
                <el-tag :type="typeTag(row.type)">{{ typeText(row.type) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="90">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                  {{ row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="260" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="editMenu(row)" v-permission="'menu:update'">编辑</el-button>
                <el-button
                  size="small"
                  :type="row.status === 1 ? 'warning' : 'success'"
                  @click="toggleMenuStatus(row)"
                  v-permission="'menu:update_status'"
                >
                  {{ row.status === 1 ? '禁用' : '启用' }}
                </el-button>
                <el-button size="small" type="danger" @click="deleteMenu(row)" v-permission="'menu:delete'">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </section>
      </div>
    </el-card>

    <el-dialog v-model="showCreateDialog" title="新增菜单" width="680px">
      <el-form :model="createForm" :rules="formRules" ref="createFormRef" label-width="110px">
        <MenuFormFields v-model="createForm" :menu-options="menuOptions" />
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="createMenu" :loading="createLoading" v-permission="'menu:create'">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showEditDialog" title="编辑菜单" width="680px">
      <el-form :model="editForm" :rules="formRules" ref="editFormRef" label-width="110px">
        <MenuFormFields v-model="editForm" :menu-options="editMenuOptions" />
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="updateMenu" :loading="editLoading" v-permission="'menu:update'">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, defineComponent, h, onMounted, ref, resolveComponent } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import { menuApi } from '../api'

const loading = ref(false)
const createLoading = ref(false)
const editLoading = ref(false)
const menus = ref([])
const menuList = ref([])
const selectedMenuId = ref(0)
const showCreateDialog = ref(false)
const showEditDialog = ref(false)

const emptyForm = () => ({
  parent_id: 0,
  name: '',
  title: '',
  path: '',
  component: '',
  icon: '',
  sort: 0,
  type: 1,
  permission: '',
  method: ''
})

const createForm = ref(emptyForm())
const editForm = ref({ id: null, ...emptyForm() })
const createFormRef = ref()
const editFormRef = ref()

const formRules = {
  title: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
  name: [
    { required: true, message: '请输入菜单标识', trigger: 'blur' },
    { pattern: /^[a-zA-Z_][a-zA-Z0-9_-]*$/, message: '菜单标识只能包含字母、数字、下划线和横线，且以字母或下划线开头', trigger: 'blur' }
  ],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }]
}

const iconOptions = [
  { label: '文档', value: 'Document' },
  { label: '菜单', value: 'Menu' },
  { label: '用户', value: 'User' },
  { label: '用户（填充）', value: 'UserFilled' },
  { label: '钥匙', value: 'Key' },
  { label: '网格', value: 'Grid' },
  { label: '设置', value: 'Setting' },
  { label: '管理', value: 'Management' },
  { label: '数据分析', value: 'DataAnalysis' },
  { label: '监控', value: 'Monitor' },
  { label: '笔记', value: 'Notebook' },
  { label: '工具', value: 'Tools' }
]

const methodOptions = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE']

const MenuFormFields = defineComponent({
  name: 'MenuFormFields',
  props: {
    modelValue: { type: Object, required: true },
    menuOptions: { type: Array, required: true }
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const update = (key, value) => {
      emit('update:modelValue', { ...props.modelValue, [key]: value })
    }
    const input = (key, placeholder) => h(resolveComponent('el-input'), {
      modelValue: props.modelValue[key],
      placeholder,
      clearable: true,
      'onUpdate:modelValue': value => update(key, value)
    })
    return () => [
      h(resolveComponent('el-form-item'), { label: '父级菜单', prop: 'parent_id' }, () => h(resolveComponent('el-tree-select'), {
        modelValue: props.modelValue.parent_id,
        data: props.menuOptions,
        props: { children: 'children', label: 'title', value: 'id', disabled: 'disabled' },
        placeholder: '请选择父级菜单',
        clearable: true,
        checkStrictly: true,
        renderAfterExpand: false,
        style: 'width: 100%',
        'onUpdate:modelValue': value => update('parent_id', value || 0)
      })),
      h(resolveComponent('el-form-item'), { label: '菜单名称', prop: 'title' }, () => input('title', '请输入菜单名称')),
      h(resolveComponent('el-form-item'), { label: '菜单标识', prop: 'name' }, () => input('name', '请输入菜单标识，如 operation-logs')),
      h(resolveComponent('el-form-item'), { label: '类型', prop: 'type' }, () => h(resolveComponent('el-radio-group'), {
        modelValue: props.modelValue.type,
        'onUpdate:modelValue': value => update('type', value)
      }, () => [
        h(resolveComponent('el-radio'), { label: 1 }, () => '目录'),
        h(resolveComponent('el-radio'), { label: 2 }, () => '菜单'),
        h(resolveComponent('el-radio'), { label: 3 }, () => '按钮/API')
      ])),
      h(resolveComponent('el-form-item'), { label: '路径', prop: 'path' }, () => input('path', '菜单页面路径，如 /admin/menus')),
      h(resolveComponent('el-form-item'), { label: '组件', prop: 'component' }, () => input('component', '组件标识，用于菜单记录')),
      h(resolveComponent('el-form-item'), { label: '权限标识', prop: 'permission' }, () => input('permission', '按钮/API 权限，如 article:create')),
      h(resolveComponent('el-form-item'), { label: '请求方法', prop: 'method' }, () => h(resolveComponent('el-select'), {
        modelValue: props.modelValue.method,
        placeholder: '请选择请求方法',
        clearable: true,
        style: 'width: 100%',
        'onUpdate:modelValue': value => update('method', value)
      }, () => methodOptions.map(method => h(resolveComponent('el-option'), { key: method, label: method, value: method })))),
      h(resolveComponent('el-form-item'), { label: '图标', prop: 'icon' }, () => h(resolveComponent('el-select'), {
        modelValue: props.modelValue.icon,
        placeholder: '请选择图标',
        clearable: true,
        style: 'width: 100%',
        'onUpdate:modelValue': value => update('icon', value)
      }, () => iconOptions.map(icon => h(resolveComponent('el-option'), { key: icon.value, label: icon.label, value: icon.value })))),
      h(resolveComponent('el-form-item'), { label: '排序', prop: 'sort' }, () => h(resolveComponent('el-input-number'), {
        modelValue: props.modelValue.sort,
        min: 0,
        'onUpdate:modelValue': value => update('sort', value ?? 0)
      }))
    ]
  }
})

const menuOptions = computed(() => [
  { id: 0, title: '顶级菜单', children: [] },
  ...buildMenuOptions(menus.value)
])

const editMenuOptions = computed(() => [
  { id: 0, title: '顶级菜单', children: [] },
  ...buildMenuOptions(menus.value, editForm.value.id)
])

const treeProps = {
  children: 'children',
  label: 'title'
}

const treeData = computed(() => [
  { id: 0, title: '全部菜单', icon: 'Menu', children: menus.value }
])

const selectedMenu = computed(() => {
  if (selectedMenuId.value === 0) return null
  return findMenuById(menus.value, selectedMenuId.value)
})

const selectedMenuTitle = computed(() => {
  if (!selectedMenu.value) return '全部菜单'
  return selectedMenu.value.children?.length ? `${selectedMenu.value.title} / 直接下级` : `${selectedMenu.value.title} / 当前菜单`
})

const tableMenus = computed(() => {
  if (selectedMenuId.value === 0) return menuList.value
  if (!selectedMenu.value) return []
  const rows = selectedMenu.value.children?.length ? selectedMenu.value.children : [selectedMenu.value]
  return rows.map(toTableRow)
})

const fetchMenus = async () => {
  loading.value = true
  try {
    const data = await menuApi.list()
    const tree = buildMenuTree(flattenMenus(data || []))
    menus.value = tree
    menuList.value = flattenMenus(tree)
    if (selectedMenuId.value !== 0 && !findMenuById(menus.value, selectedMenuId.value)) {
      selectedMenuId.value = 0
    }
  } catch (error) {
    ElMessage.error('获取菜单列表失败')
  } finally {
    loading.value = false
  }
}

const flattenMenus = (items) => {
  const list = []
  const walk = (nodes) => {
    for (const node of nodes) {
      list.push(toTableRow(node))
      if (node.children?.length) walk(node.children)
    }
  }
  walk(items)
  return list
}

const toTableRow = (node) => {
  const { children, ...row } = node
  return row
}

const buildMenuTree = (list) => {
  const menuMap = new Map()
  const roots = []
  for (const menu of list) {
    menuMap.set(menu.id, { ...menu, children: [], level: 0, parent_title: '' })
  }
  for (const menu of list) {
    const menuItem = menuMap.get(menu.id)
    if (!menu.parent_id) {
      roots.push(menuItem)
    } else {
      const parent = menuMap.get(menu.parent_id)
      if (parent) {
        parent.children.push(menuItem)
      }
      else roots.push(menuItem)
    }
  }
  annotateMenuTree(roots, 0, '')
  return roots
}

const annotateMenuTree = (nodes, level, parentTitle) => {
  for (const node of nodes) {
    node.level = level
    node.parent_title = parentTitle
    if (node.children?.length) annotateMenuTree(node.children, level + 1, node.title)
  }
}

const findMenuById = (nodes, id) => {
  for (const node of nodes) {
    if (node.id === id) return node
    const child = findMenuById(node.children || [], id)
    if (child) return child
  }
  return null
}

const handleTreeNodeClick = (node) => {
  selectedMenuId.value = node.id
}

const buildMenuOptions = (items, disabledRootID = null) => {
  const disabledIDs = disabledRootID ? collectChildIDs(disabledRootID) : new Set()
  if (disabledRootID) disabledIDs.add(disabledRootID)

  const walk = (nodes) => nodes.map(menu => ({
    id: menu.id,
    title: `${menu.title}（${typeText(menu.type)}）`,
    disabled: menu.type === 3 || disabledIDs.has(menu.id),
    children: menu.children?.length ? walk(menu.children) : []
  }))
  return walk(items)
}

const collectChildIDs = (id) => {
  const ids = new Set()
  const walk = (parentID) => {
    for (const item of menuList.value.filter(menu => menu.parent_id === parentID)) {
      ids.add(item.id)
      walk(item.id)
    }
  }
  walk(id)
  return ids
}

const openCreateDialog = () => {
  const parentMenu = selectedMenu.value
  createForm.value = {
    ...emptyForm(),
    parent_id: parentMenu && parentMenu.type !== 3 ? parentMenu.id : 0
  }
  showCreateDialog.value = true
}

const createMenu = async () => {
  if (!createFormRef.value) return
  const valid = await createFormRef.value.validate().catch(() => false)
  if (!valid) return

  createLoading.value = true
  try {
    await menuApi.create(normalizeForm(createForm.value))
    ElMessage.success('菜单创建成功')
    showCreateDialog.value = false
    createForm.value = emptyForm()
    fetchMenus()
  } catch (error) {
    ElMessage.error(error.message || '创建菜单失败')
  } finally {
    createLoading.value = false
  }
}

const editMenu = (menu) => {
  editForm.value = {
    id: menu.id,
    parent_id: menu.parent_id || 0,
    name: menu.name || '',
    title: menu.title || '',
    path: menu.path || '',
    component: menu.component || '',
    icon: menu.icon || '',
    sort: menu.sort || 0,
    type: menu.type || 1,
    permission: menu.permission || '',
    method: menu.method || ''
  }
  showEditDialog.value = true
}

const updateMenu = async () => {
  if (!editFormRef.value) return
  const valid = await editFormRef.value.validate().catch(() => false)
  if (!valid) return

  editLoading.value = true
  try {
    await menuApi.update(editForm.value.id, normalizeForm(editForm.value))
    ElMessage.success('菜单更新成功')
    showEditDialog.value = false
    fetchMenus()
  } catch (error) {
    ElMessage.error(error.message || '更新菜单失败')
  } finally {
    editLoading.value = false
  }
}

const normalizeForm = (form) => ({
  parent_id: form.parent_id || 0,
  name: form.name,
  title: form.title,
  path: form.type === 2 ? form.path : '',
  component: form.type === 2 ? form.component : '',
  icon: form.icon,
  sort: form.sort || 0,
  type: form.type,
  permission: form.type === 3 ? form.permission : '',
  method: form.type === 3 ? form.method : ''
})

const toggleMenuStatus = async (menu) => {
  const newStatus = menu.status === 1 ? 0 : 1
  const action = newStatus === 1 ? '启用' : '禁用'
  try {
    await ElMessageBox.confirm(`确定要${action}菜单 ${menu.title} 吗？`, '确认操作', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await menuApi.updateStatus(menu.id, { status: newStatus })
    ElMessage.success(`${action}成功`)
    fetchMenus()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error(error.message || `${action}失败`)
  }
}

const deleteMenu = async (menu) => {
  try {
    await ElMessageBox.confirm(`确定要删除菜单 ${menu.title} 吗？此操作不可逆！`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await menuApi.delete(menu.id)
    ElMessage.success('删除成功')
    fetchMenus()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error(error.message || '删除失败')
  }
}

const typeText = (type) => {
  const map = { 1: '目录', 2: '菜单', 3: '按钮/API' }
  return map[type] || '未知'
}

const typeTag = (type) => {
  const map = { 1: 'primary', 2: 'success', 3: 'warning' }
  return map[type] || 'info'
}

onMounted(fetchMenus)
</script>

<style scoped>
.menu-management {
  padding: 20px;
}

.card-header,
.header-actions,
.menu-title {
  display: flex;
  align-items: center;
}

.card-header {
  justify-content: space-between;
}

.header-actions {
  gap: 8px;
}

.menu-workspace {
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr);
  gap: 16px;
  min-height: 520px;
}

.menu-tree-panel {
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  padding: 12px;
  background: #fbfcfe;
  overflow: auto;
}

.tree-title {
  margin-bottom: 10px;
  color: #303133;
  font-weight: 600;
}

.tree-node,
.list-toolbar {
  display: flex;
  align-items: center;
}

.tree-node {
  gap: 6px;
  min-width: 0;
}

.tree-node span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.menu-list-panel {
  min-width: 0;
}

.list-toolbar {
  justify-content: space-between;
  margin-bottom: 12px;
}

.list-title {
  color: #303133;
  font-size: 16px;
  font-weight: 600;
  line-height: 24px;
}

.list-subtitle {
  color: #909399;
  font-size: 12px;
  line-height: 18px;
}

.menu-title {
  gap: 8px;
  font-weight: 600;
}

.menu-depth {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 28px;
  height: 20px;
  padding: 0 6px;
  border-radius: 4px;
  background: #eef2ff;
  color: #4338ca;
  font-size: 12px;
  font-weight: 600;
}

@media (max-width: 960px) {
  .menu-workspace {
    grid-template-columns: 1fr;
  }
}
</style>
