<template>
  <div class="resource-manager">
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="资源名称">
          <el-input
            v-model="searchForm.resourceName"
            placeholder="请输入资源名称"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="资源类型">
          <el-select
            v-model="searchForm.resourceType"
            placeholder="请选择资源类型"
            clearable
          >
            <el-option label="文章" value="article" />
            <el-option label="视频" value="video" />
            <el-option label="图片" value="image" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            clearable
          >
            <el-option label="待处理" value="pending" />
            <el-option label="处理中" value="processing" />
            <el-option label="已完成" value="completed" />
            <el-option label="失败" value="failed" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetSearch">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="table-card">
      <template #header>
        <div class="card-header">
          <span class="title">资源列表</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增资源
          </el-button>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column type="index" width="50" align="center" />
        <el-table-column prop="resourceName" label="资源名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="resourceType" label="资源类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getResourceTypeTag(row.resourceType)">
              {{ getResourceTypeLabel(row.resourceType) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTag(row.status)">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button
                type="primary"
                link
                @click="handleEdit(row)"
              >
                编辑
              </el-button>
              <el-button
                type="danger"
                link
                @click="handleDelete(row)"
              >
                删除
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          :current-page="currentPage"
          :page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增资源' : '编辑资源'"
      width="600px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="资源名称" prop="resourceName">
          <el-input v-model="form.resourceName" placeholder="请输入资源名称" />
        </el-form-item>
        <el-form-item label="资源类型" prop="resourceType">
          <el-select v-model="form.resourceType" placeholder="请选择资源类型">
            <el-option label="文章" value="article" />
            <el-option label="视频" value="video" />
            <el-option label="图片" value="image" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status" placeholder="请选择状态">
            <el-option label="待处理" value="pending" />
            <el-option label="处理中" value="processing" />
            <el-option label="已完成" value="completed" />
            <el-option label="失败" value="failed" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { Plus, Refresh, Search } from '@element-plus/icons-vue';
import { resourceMethods } from '../modules/resource.js';

export default {
    name: 'ResourceManager',
    data() {
        return {
            tableData: [],
            currentPage: 1,
            pageSize: 10,
            total: 0,
            dialogVisible: false,
            dialogType: 'add',
            searchForm: {
                resourceName: '',
                resourceType: '',
                status: ''
            },
            form: {
                resourceName: '',
                resourceType: '',
                status: ''
            },
            loading: false,
            rules: {
                resourceName: [
                    { required: true, message: '请输入资源名称', trigger: 'blur' }
                ],
                resourceType: [
                    { required: true, message: '请选择资源类型', trigger: 'change' }
                ],
                status: [
                    { required: true, message: '请选择状态', trigger: 'change' }
                ]
            }
        }
    },
    methods: {
        ...resourceMethods,
        getResourceTypeTag(type) {
            const typeMap = {
                article: '',
                video: 'success',
                image: 'warning'
            }
            return typeMap[type] || ''
        },
        getResourceTypeLabel(type) {
            const typeMap = {
                article: '文章',
                video: '视频',
                image: '图片'
            }
            return typeMap[type] || type
        },
        getStatusTag(status) {
            const statusMap = {
                pending: 'info',
                processing: 'warning',
                completed: 'success',
                failed: 'danger'
            }
            return statusMap[status] || ''
        },
        getStatusLabel(status) {
            const statusMap = {
                pending: '待处理',
                processing: '处理中',
                completed: '已完成',
                failed: '失败'
            }
            return statusMap[status] || status
        }
    },
    mounted() {
        this.loadData()
    }
}
</script>

<style scoped>
.resource-manager {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.search-card {
  margin-bottom: 16px;
}

.table-card {
  flex: 1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header .title {
  font-size: 16px;
  font-weight: 500;
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

:deep(.el-card__header) {
  padding: 12px 20px;
}

:deep(.el-form--inline .el-form-item) {
  margin-right: 16px;
  margin-bottom: 0;
}

:deep(.el-button-group) {
  display: flex;
  gap: 8px;
}

:deep(.el-table) {
  margin: 16px 0;
}

:deep(.el-dialog__body) {
  padding: 20px 40px;
}
</style> 