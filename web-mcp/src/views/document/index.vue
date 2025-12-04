<template>
  <div class="document-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-button @click="goBack" :icon="ArrowLeft">返回</el-button>
            <span class="title">文档列表 - {{ libraryName }}</span>
          </div>
          <el-upload
            :action="uploadUrl"
            :data="{ library_id: libraryId }"
            :show-file-list="false"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :before-upload="beforeUpload"
            accept=".md,.pdf,.docx"
          >
            <el-button type="primary">
              <el-icon><Upload /></el-icon>
              上传文档
            </el-button>
          </el-upload>
        </div>
      </template>

      <el-table :data="documents" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" show-overflow-tooltip />
        <el-table-column prop="file_type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag>{{ row.file_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="file_size" label="大小" width="100">
          <template #default="{ row }">
            {{ formatSize(row.file_size) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="上传时间" width="180" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-popconfirm title="确定删除该文档？" @confirm="handleDelete(row.id)">
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
        @size-change="fetchDocuments"
        @current-change="fetchDocuments"
        style="margin-top: 20px; justify-content: flex-end;"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { getDocuments, deleteDocument } from '@/api/document'
import { getLibrary } from '@/api/library'
import type { Document } from '@/api/document'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()

const libraryId = computed(() => Number(route.params.id))
const libraryName = ref('')

const loading = ref(false)
const documents = ref<Document[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const uploadUrl = computed(() => `/api/v1/documents/upload`)

const fetchLibrary = async () => {
  const res = await getLibrary(libraryId.value)
  if (res.code === 0) {
    libraryName.value = res.data.name
  }
}

const fetchDocuments = async () => {
  loading.value = true
  try {
    const res = await getDocuments({
      library_id: libraryId.value,
      page: page.value,
      page_size: pageSize.value
    })
    if (res.code === 0) {
      documents.value = res.data.list
      total.value = res.data.total
    }
  } finally {
    loading.value = false
  }
}

const handleDelete = async (id: number) => {
  await deleteDocument(id)
  ElMessage.success('删除成功')
  fetchDocuments()
}

const handleUploadSuccess = () => {
  ElMessage.success('上传成功，正在处理中...')
  fetchDocuments()
}

const handleUploadError = () => {
  ElMessage.error('上传失败')
}

const beforeUpload = (file: File) => {
  const allowedTypes = ['.md', '.pdf', '.docx']
  const ext = file.name.substring(file.name.lastIndexOf('.')).toLowerCase()
  if (!allowedTypes.includes(ext)) {
    ElMessage.error('只支持 .md, .pdf, .docx 格式')
    return false
  }
  return true
}

const formatSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(1) + ' MB'
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    active: 'success',
    processing: 'warning',
    processed: 'success',
    failed: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: '待处理',
    processing: '处理中',
    processed: '已完成',
    failed: '失败'
  }
  return map[status] || status
}

const goBack = () => {
  router.push('/')
}

onMounted(() => {
  fetchLibrary()
  fetchDocuments()
})
</script>

<style scoped>
.document-page {
  max-width: 880px;
  margin: 0 auto;
  padding: 40px 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title {
  font-size: 16px;
  font-weight: 500;
  color: #1f2937;
}

:deep(.el-card) {
  border-radius: 12px;
  border: 1px solid #e5e7eb;
}
</style>
