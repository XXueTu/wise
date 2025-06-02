-- 资源表
CREATE TABLE IF NOT EXISTS resources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    url TEXT NOT NULL, -- 资源URL
    title TEXT NOT NULL, -- 资源标题
    describe TEXT NOT NULL, -- 资源描述
    content TEXT NOT NULL, -- 资源内容
    type TEXT NOT NULL, -- 资源类型（如：wechat, zhihu等）
    tags TEXT NOT NULL, -- 资源标签
    created_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime')),
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime'))
);

-- 模型表
CREATE TABLE IF NOT EXISTS models (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    base_url TEXT NOT NULL, -- 资源链接
    config TEXT NOT NULL, -- 资源配置
    type TEXT NOT NULL, -- 资源类型
    model_name TEXT NOT NULL, -- 模型名称
    model_real_name TEXT NOT NULL, -- 模型真实名称
    status TEXT NOT NULL, -- 状态
    tag TEXT NOT NULL, -- 标签
    created_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime')),
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime'))
);

-- 标签表
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uid TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    color TEXT NOT NULL,
    icon TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime')),
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime'))
);

-- 任务表
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tid TEXT NOT NULL, -- 任务唯一标识
    name TEXT NOT NULL, -- 任务名称
    types TEXT NOT NULL, -- 任务类型
    status TEXT NOT NULL, -- 任务状态
    current_state TEXT NOT NULL, -- 当前状态机
    total_steps INTEGER NOT NULL, -- 总步骤
    current_step INTEGER NOT NULL, -- 当前步骤
    retry_count INTEGER, -- 重试次数
    params TEXT NOT NULL, -- 任务参数
    result TEXT NOT NULL, -- 任务结果
    duration INTEGER NOT NULL, -- 任务耗时 ms
    error TEXT NOT NULL, -- 任务错误
    extend TEXT NOT NULL, -- 扩展字段
    created_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime')), -- 创建时间
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime')) -- 更新时间
);

-- 任务计划表
CREATE TABLE IF NOT EXISTS task_plans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tid TEXT NOT NULL, -- 任务唯一标识
    pid TEXT NOT NULL, -- 任务计划唯一标识
    before_pid TEXT NOT NULL, -- 上一个任务计划唯一标识
    next TEXT NOT NULL, -- 下一个任务类型
    types TEXT NOT NULL, -- 任务类型
    name TEXT NOT NULL, -- 任务名称
    `index` INTEGER NOT NULL, -- 任务计划索引
    status TEXT NOT NULL, -- 任务状态
    params TEXT NOT NULL, -- 任务参数
    result TEXT NOT NULL, -- 任务结果
    duration INTEGER NOT NULL, -- 任务耗时 ms
    error TEXT NOT NULL, -- 任务错误
    created_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime')), -- 创建时间
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime')) -- 更新时间
);