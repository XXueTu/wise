package model

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func InitDB() *bun.DB {
	// 确保data目录存在
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		panic(err)
	}

	// 数据库文件路径
	dbPath := filepath.Join(dataDir, "wise.db")

	// 连接数据库
	sqldb, err := sql.Open(sqliteshim.ShimName, dbPath)
	if err != nil {
		panic(err)
	}

	// 设置连接池参数
	sqldb.SetMaxOpenConns(1) // SQLite 只支持一个写连接
	sqldb.SetMaxIdleConns(1)

	// 创建 bun DB 实例
	db := bun.NewDB(sqldb, sqlitedialect.New())
	// 创建
	db.NewCreateTable().Model((*Resource)(nil)).IfNotExists().Exec(context.Background(),
		(*Resource)(nil),
		(*Models)(nil),
	)

	// 初始化表数据
	NewModelsModel(db).InitData()
	return db
}

// 自定义数据库日志输出
