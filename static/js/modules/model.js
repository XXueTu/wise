import axios from 'axios';

// 模型管理模块方法
export const modelMethods = {
    // 加载数据
    loadData() {
        axios.get('/api/models', {
            params: {
                page: this.currentPage,
                size: this.pageSize,
                modelName: this.searchForm.modelName,
                tag: this.searchForm.tag,
                status: this.searchForm.status
            }
        }).then(response => {
            this.tableData = response.data.list
            this.total = response.data.total
        }).catch(error => {
            this.$message.error('加载数据失败：' + error.message)
        })
    },

    // 搜索
    handleSearch() {
        this.currentPage = 1
        this.loadData()
    },

    // 重置搜索
    resetSearch() {
        this.searchForm = {
            modelName: '',
            tag: '',
            status: ''
        }
        this.handleSearch()
    },

    // 新增
    handleAdd() {
        this.dialogType = 'add'
        this.form = {
            baseUrl: '',
            config: '',
            type: '',
            modelName: '',
            modelRealName: '',
            status: 'active',
            tag: 'function'
        }
        this.dialogVisible = true
    },

    // 编辑
    handleEdit(row) {
        this.dialogType = 'edit'
        this.form = { ...row }
        this.dialogVisible = true
    },

    // 删除
    handleDelete(row) {
        this.$confirm('确认删除该模型吗？', '提示', {
            type: 'warning'
        }).then(() => {
            axios.delete(`/api/models/${row.id}`).then(() => {
                this.$message.success('删除成功')
                this.loadData()
            }).catch(error => {
                this.$message.error('删除失败：' + error.message)
            })
        }).catch(() => {})
    },

    // 提交表单
    handleSubmit() {
        if (this.dialogType === 'add') {
            axios.post('/api/models', this.form).then(() => {
                this.$message.success('添加成功')
                this.dialogVisible = false
                this.loadData()
            }).catch(error => {
                this.$message.error('添加失败：' + error.message)
            })
        } else {
            axios.put(`/api/models/${this.form.id}`, this.form).then(() => {
                this.$message.success('更新成功')
                this.dialogVisible = false
                this.loadData()
            }).catch(error => {
                this.$message.error('更新失败：' + error.message)
            })
        }
    },

    // 分页大小改变
    handleSizeChange(val) {
        this.pageSize = val
        this.loadData()
    },

    // 页码改变
    handleCurrentChange(val) {
        this.currentPage = val
        this.loadData()
    }
}; 