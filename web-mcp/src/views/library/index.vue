<template>
  <div class="library-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>库列表</span>
          <el-button type="primary" @click="showCreateDialog">
            <el-icon><Plus /></el-icon>
            新建库
          </el-button>
        </div>
      </template>

      <el-table :data="libraries" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="default_version" label="默认版本" width="120" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="goToDocuments(row.id)">
              文档
            </el-button>
            <el-button type="primary" link @click="showEditDialog(row)">
              编辑
            </el-button>
            <el-popconfirm title="确定删除该库？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="fetchLibraries"
        @current-change="fetchLibraries"
        style="margin-top: 20px; justify-content: flex-end;"
      />
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑库' : '新建库'"
      width="500px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入库名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入库描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getLibraries, createLibrary, updateLibrary, deleteLibrary } from '@/api/library'
import type { LibraryListItem } from '@/api/library'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()

const loading = ref(false)
const libraries = ref<LibraryListItem[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const submitting = ref(false)
const formRef = ref<FormInstance>()

const form = reactive({
  name: '',
  description: ''
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入库名称', trigger: 'blur' }]
}

const fetchLibraries = async () => {
  loading.value = true
  try {
    const res = await getLibraries({
      page: page.value,
      page_size: pageSize.value
    })
    libraries.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

const showCreateDialog = () => {
  isEdit.value = false
  editId.value = null
  form.name = ''
  form.description = ''
  dialogVisible.value = true
}

const showEditDialog = (row: LibraryListItem) => {
  isEdit.value = true
  editId.value = row.id
  form.name = row.name
  form.description = ''  // 列表项没有 description，需要从详情获取
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    submitting.value = true
    try {
      if (isEdit.value && editId.value) {
        await updateLibrary(editId.value, form)
        ElMessage.success('更新成功')
      } else {
        await createLibrary(form)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      fetchLibraries()
    } finally {
      submitting.value = false
    }
  })
}

const handleDelete = async (id: number) => {
  await deleteLibrary(id)
  ElMessage.success('删除成功')
  fetchLibraries()
}

const goToDocuments = (id: number) => {
  router.push(`/libraries/${id}/documents`)
}

onMounted(() => {
  fetchLibraries()
})
</script>

<style scoped>
.library-page {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
