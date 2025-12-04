<template>
  <div class="search-page">
    <el-card class="search-card">
      <template #header>
        <span>搜索测试</span>
      </template>

      <el-form :model="searchForm" inline>
        <el-form-item label="选择库">
          <el-select v-model="searchForm.library_id" placeholder="请选择库" style="width: 200px">
            <el-option
              v-for="lib in libraries"
              :key="lib.id"
              :label="`${lib.name} (${lib.version})`"
              :value="lib.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="搜索模式">
          <el-select v-model="searchForm.mode" placeholder="全部" style="width: 120px" clearable>
            <el-option label="代码" value="code" />
            <el-option label="信息" value="info" />
          </el-select>
        </el-form-item>
        <el-form-item label="查询内容" style="flex: 1">
          <el-input
            v-model="searchForm.query"
            placeholder="请输入搜索关键词"
            @keyup.enter="handleSearch"
            clearable
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch" :loading="loading">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="result-card" v-if="results.length > 0 || searched">
      <template #header>
        <div class="result-header">
          <span>搜索结果</span>
          <span class="total">共 {{ total }} 条</span>
        </div>
      </template>

      <div v-if="results.length === 0" class="empty">
        <el-empty description="暂无搜索结果" />
      </div>

      <div v-else class="result-list">
        <div v-for="item in results" :key="item.chunk_id" class="result-item">
          <div class="result-meta">
            <el-tag size="small" :type="item.chunk_type === 'code' ? 'success' : 'info'">
              {{ item.chunk_type }}
            </el-tag>
            <span class="score">
              综合分数: {{ item.score.toFixed(3) }}
              <el-tooltip content="向量相似度 / BM25分数">
                <span class="score-detail">
                  ({{ item.vector_score.toFixed(3) }} / {{ item.bm25_score.toFixed(3) }})
                </span>
              </el-tooltip>
            </span>
          </div>
          <div class="result-content">
            <pre>{{ item.content }}</pre>
          </div>
        </div>
      </div>

      <el-pagination
        v-if="total > 0"
        v-model:current-page="searchForm.page"
        :page-size="searchForm.limit"
        :total="total"
        layout="prev, pager, next"
        @current-change="handleSearch"
        style="margin-top: 20px; justify-content: center;"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getLibraries } from '@/api/library'
import { searchDocuments } from '@/api/search'
import type { Library } from '@/api/library'
import type { SearchResultItem } from '@/api/search'
import { ElMessage } from 'element-plus'

const libraries = ref<Library[]>([])
const loading = ref(false)
const searched = ref(false)
const results = ref<SearchResultItem[]>([])
const total = ref(0)

const searchForm = reactive({
  library_id: null as number | null,
  query: '',
  mode: '',
  page: 1,
  limit: 10
})

const fetchLibraries = async () => {
  const res = await getLibraries({ page: 1, page_size: 100 })
  if (res.code === 0) {
    libraries.value = res.data.list
  }
}

const handleSearch = async () => {
  if (!searchForm.library_id) {
    ElMessage.warning('请先选择一个库')
    return
  }
  if (!searchForm.query.trim()) {
    ElMessage.warning('请输入搜索关键词')
    return
  }

  loading.value = true
  searched.value = true
  
  try {
    const res = await searchDocuments({
      library_id: searchForm.library_id,
      query: searchForm.query,
      mode: searchForm.mode || undefined,
      page: searchForm.page,
      limit: searchForm.limit
    })
    
    if (res.code === 0) {
      results.value = res.data.results
      total.value = res.data.total
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchLibraries()
})
</script>

<style scoped>
.search-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.search-card :deep(.el-form) {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.total {
  color: #909399;
  font-size: 14px;
}

.result-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.result-item {
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 16px;
  background-color: #fafafa;
}

.result-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.score {
  font-size: 13px;
  color: #606266;
}

.score-detail {
  color: #909399;
  font-size: 12px;
}

.result-content {
  background-color: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 12px;
  overflow-x: auto;
}

.result-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.5;
}

.empty {
  padding: 40px 0;
}
</style>
