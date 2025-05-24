// 资源管理组件模板
export const resourceTemplate = `
    <div>
        <el-card>
            <template #header>
                <div class="card-header">
                    <span>资源列表</span>
                    <el-button type="primary" @click="handleAdd">新增资源</el-button>
                </div>
            </template>
            
            <!-- 搜索表单 -->
            <el-form :inline="true" :model="searchForm" class="search-form">
                <el-form-item label="标题">
                    <el-input v-model="searchForm.title" placeholder="请输入标题"></el-input>
                </el-form-item>
                <el-form-item label="类型">
                    <el-select v-model="searchForm.type" placeholder="请选择类型">
                        <el-option label="全部" value=""></el-option>
                        <el-option label="微信公众号" value="wechat"></el-option>
                        <el-option label="知乎" value="zhihu"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="handleSearch">搜索</el-button>
                    <el-button @click="resetSearch">重置</el-button>
                </el-form-item>
            </el-form>

            <!-- 数据表格 -->
            <el-table :data="tableData" style="width: 100%">
                <el-table-column prop="id" label="ID" width="80"></el-table-column>
                <el-table-column prop="title" label="标题"></el-table-column>
                <el-table-column prop="type" label="类型" width="120"></el-table-column>
                <el-table-column prop="createdAt" label="创建时间" width="180"></el-table-column>
                <el-table-column label="操作" width="200">
                    <template #default="scope">
                        <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
                        <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>

            <!-- 分页 -->
            <div class="pagination">
                <el-pagination
                    v-model:current-page="currentPage"
                    v-model:page-size="pageSize"
                    :page-sizes="[10, 20, 50, 100]"
                    :total="total"
                    layout="total, sizes, prev, pager, next"
                    @size-change="handleSizeChange"
                    @current-change="handleCurrentChange">
                </el-pagination>
            </div>
        </el-card>

        <!-- 编辑对话框 -->
        <el-dialog
            v-model="dialogVisible"
            :title="dialogType === 'add' ? '新增资源' : '编辑资源'"
            width="50%">
            <el-form :model="form" label-width="80px">
                <el-form-item label="标题">
                    <el-input v-model="form.title"></el-input>
                </el-form-item>
                <el-form-item label="URL">
                    <el-input v-model="form.url"></el-input>
                </el-form-item>
                <el-form-item label="类型">
                    <el-select v-model="form.type" placeholder="请选择类型">
                        <el-option label="微信公众号" value="wechat"></el-option>
                        <el-option label="知乎" value="zhihu"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="内容">
                    <el-input type="textarea" v-model="form.content" rows="4"></el-input>
                </el-form-item>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="dialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="handleSubmit">确定</el-button>
                </span>
            </template>
        </el-dialog>
    </div>
`;

// 模型管理组件模板
export const modelTemplate = `
    <div>
        <el-card>
            <template #header>
                <div class="card-header">
                    <span>模型列表</span>
                    <el-button type="primary" @click="handleAdd">新增模型</el-button>
                </div>
            </template>
            
            <!-- 搜索表单 -->
            <el-form :inline="true" :model="searchForm" class="search-form">
                <el-form-item label="模型名称">
                    <el-input v-model="searchForm.modelName" placeholder="请输入模型名称"></el-input>
                </el-form-item>
                <el-form-item label="标签">
                    <el-select v-model="searchForm.tag" placeholder="请选择标签">
                        <el-option label="全部" value=""></el-option>
                        <el-option label="function" value="function"></el-option>
                        <el-option label="chat" value="chat"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="状态">
                    <el-select v-model="searchForm.status" placeholder="请选择状态">
                        <el-option label="全部" value=""></el-option>
                        <el-option label="可用" value="active"></el-option>
                        <el-option label="不可用" value="inactive"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="handleSearch">搜索</el-button>
                    <el-button @click="resetSearch">重置</el-button>
                </el-form-item>
            </el-form>

            <!-- 数据表格 -->
            <el-table :data="tableData" style="width: 100%">
                <el-table-column prop="id" label="ID" width="80"></el-table-column>
                <el-table-column prop="modelName" label="模型名称"></el-table-column>
                <el-table-column prop="modelRealName" label="实际名称"></el-table-column>
                <el-table-column prop="type" label="类型" width="120"></el-table-column>
                <el-table-column prop="tag" label="标签" width="120"></el-table-column>
                <el-table-column prop="status" label="状态" width="100">
                    <template #default="scope">
                        <el-tag :type="scope.row.status === 'active' ? 'success' : 'danger'">
                            {{ scope.row.status === 'active' ? '可用' : '不可用' }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="createdAt" label="创建时间" width="180"></el-table-column>
                <el-table-column label="操作" width="200">
                    <template #default="scope">
                        <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
                        <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>

            <!-- 分页 -->
            <div class="pagination">
                <el-pagination
                    v-model:current-page="currentPage"
                    v-model:page-size="pageSize"
                    :page-sizes="[10, 20, 50, 100]"
                    :total="total"
                    layout="total, sizes, prev, pager, next"
                    @size-change="handleSizeChange"
                    @current-change="handleCurrentChange">
                </el-pagination>
            </div>
        </el-card>

        <!-- 编辑对话框 -->
        <el-dialog
            v-model="dialogVisible"
            :title="dialogType === 'add' ? '新增模型' : '编辑模型'"
            width="50%">
            <el-form :model="form" label-width="100px">
                <el-form-item label="基础URL">
                    <el-input v-model="form.baseUrl"></el-input>
                </el-form-item>
                <el-form-item label="配置">
                    <el-input type="textarea" v-model="form.config" rows="4"></el-input>
                </el-form-item>
                <el-form-item label="类型">
                    <el-input v-model="form.type"></el-input>
                </el-form-item>
                <el-form-item label="模型名称">
                    <el-input v-model="form.modelName"></el-input>
                </el-form-item>
                <el-form-item label="实际名称">
                    <el-input v-model="form.modelRealName"></el-input>
                </el-form-item>
                <el-form-item label="状态">
                    <el-select v-model="form.status">
                        <el-option label="可用" value="active"></el-option>
                        <el-option label="不可用" value="inactive"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="标签">
                    <el-select v-model="form.tag">
                        <el-option label="function" value="function"></el-option>
                        <el-option label="chat" value="chat"></el-option>
                    </el-select>
                </el-form-item>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="dialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="handleSubmit">确定</el-button>
                </span>
            </template>
        </el-dialog>
    </div>
`; 