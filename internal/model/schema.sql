-- 资源表
CREATE TABLE IF NOT EXISTS resources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    rid TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    type TEXT NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP,'localtime')),
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP,'localtime'))
);

-- 模型表
CREATE TABLE IF NOT EXISTS models (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    mid TEXT NOT NULL,
    `name` TEXT NOT NULL,
    description TEXT NOT NULL,
    `type` TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP,'localtime')),
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP,'localtime'))
);

-- 标签表
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uid TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    color TEXT NOT NULL,
    icon TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP,'localtime')),
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP,'localtime'))
);

-- 任务表
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tid TEXT NOT NULL,
    name TEXT NOT NULL,
    types TEXT NOT NULL,
    status TEXT NOT NULL,
    current_state TEXT NOT NULL,
    total_steps INTEGER NOT NULL,
    current_step INTEGER NOT NULL,
    retry_count INTEGER NULL,
    params TEXT NOT NULL,
    result TEXT NOT NULL,
    duration INTEGER NOT NULL,
    error TEXT NOT NULL,
    extend TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP,'localtime')),
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP,'localtime'))
);

-- 任务计划表
CREATE TABLE IF NOT EXISTS task_plans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tid TEXT NOT NULL,
    pid TEXT NOT NULL,
    types TEXT NOT NULL,
    `name` TEXT NOT NULL,
    `index` INTEGER NOT NULL,
    `status` TEXT NOT NULL,
    params TEXT NOT NULL,
    result TEXT NOT NULL,
    duration INTEGER NOT NULL,
    error TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP,'localtime')),
    updated_at TIMESTAMP NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP,'localtime'))
); 