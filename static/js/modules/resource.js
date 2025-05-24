import axios from 'axios';

// 资源管理模块方法
export const resourceMethods = {
    // 加载数据
    loadData() {
        axios.get('/api/resources', {
            params: {
                page: this.currentPage,
                size: this.pageSize,
                title: this.searchForm.title,
                type: this.searchForm.type
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
            title: '',
            type: ''
        }
        this.handleSearch()
    },

    // 新增
    handleAdd() {
        this.dialogType = 'add'
        this.form = {
            title: '',
            url: '',
            type: '',
            content: ''
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
        this.$confirm('确认删除该资源吗？', '提示', {
            type: 'warning'
        }).then(() => {
            axios.delete(`/api/resources/${row.id}`).then(() => {
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
            axios.post('/api/resources', this.form).then(() => {
                this.$message.success('添加成功')
                this.dialogVisible = false
                this.loadData()
            }).catch(error => {
                this.$message.error('添加失败：' + error.message)
            })
        } else {
            axios.put(`/api/resources/${this.form.id}`, this.form).then(() => {
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