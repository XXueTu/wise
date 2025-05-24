<template>
  <div class="model-manager">
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="模型名称">
          <el-input
            v-model="searchForm.modelName"
            placeholder="请输入模型名称"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="标签">
          <el-select
            v-model="searchForm.tag"
            placeholder="请选择标签"
            clearable
          >
            <el-option label="NLP" value="nlp" />
            <el-option label="CV" value="cv" />
            <el-option label="推荐" value="recommendation" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            clearable
          >
            <el-option label="启用" value="active" />
            <el-option label="禁用" value="inactive" />
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
          <span class="title">模型列表</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增模型
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
        <el-table-column prop="modelName" label="模型名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="tag" label="标签" width="120">
          <template #default="{ row }">
            <el-tag :type="getTagType(row.tag)">
              {{ getTagLabel(row.tag) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">
              {{ row.status === 'active' ? '启用' : '禁用' }}
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
      :title="dialogType === 'add' ? '新增模型' : '编辑模型'"
      width="600px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="模型名称" prop="modelName">
          <el-input v-model="form.modelName" placeholder="请输入模型名称" />
        </el-form-item>
        <el-form-item label="标签" prop="tag">
          <el-select v-model="form.tag" placeholder="请选择标签">
            <el-option label="NLP" value="nlp" />
            <el-option label="CV" value="cv" />
            <el-option label="推荐" value="recommendation" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio label="active">启用</el-radio>
            <el-radio label="inactive">禁用</el-radio>
          </el-radio-group>
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
import { modelMethods } from '../modules/model.js';

export default {
    name: 'ModelManager',
    data() {
        return {
            tableData: [],
            currentPage: 1,
            pageSize: 10,
            total: 0,
            dialogVisible: false,
            dialogType: 'add',
            searchForm: {
                modelName: '',
                tag: '',
                status: ''
            },
            form: {
                baseUrl: '',
                config: '',
                type: '',
                modelName: '',
                modelRealName: '',
                status: 'active',
                tag: 'function'
            },
            rules: {
                modelName: [
                    { required: true, message: '请输入模型名称', trigger: 'blur' }
                ],
                tag: [
                    { required: true, message: '请选择标签', trigger: 'change' }
                ],
                status: [
                    { required: true, message: '请选择状态', trigger: 'change' }
                ]
            }
        }
    },
    methods: {
        ...modelMethods,
        getTagType(tag) {
            const tagMap = {
                nlp: 'success',
                cv: 'warning',
                recommendation: 'info'
            }
            return tagMap[tag] || ''
        },
        getTagLabel(tag) {
            const tagMap = {
                nlp: 'NLP',
                cv: 'CV',
                recommendation: '推荐'
            }
            return tagMap[tag] || tag
        }
    },
    mounted() {
        this.loadData()
    }
}
</script>

<style scoped>
.model-manager {
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