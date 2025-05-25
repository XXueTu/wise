package model

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/zeromicro/go-zero/core/logx"
)

type logxWriter struct{}

// 移除 ANSI 颜色代码的正则表达式
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func (w *logxWriter) Write(p []byte) (n int, err error) {
	// 移除 ANSI 颜色代码
	cleanStr := ansiRegex.ReplaceAllString(string(p), "")

	// 提取 SQL 查询和时间信息
	parts := strings.Split(cleanStr, "  ")
	if len(parts) >= 3 {
		timeInfo := parts[0]
		duration := parts[1]
		query := strings.Join(parts[2:], "  ")

		// 清理 SQL 查询
		query = strings.ReplaceAll(query, "\"", "")  // 移除双引号
		query = strings.ReplaceAll(query, "\\", "")  // 移除反斜杠
		query = strings.ReplaceAll(query, "\n", " ") // 将换行替换为空格
		query = strings.ReplaceAll(query, "  ", " ") // 移除多余的空格
		query = strings.TrimSpace(query)             // 移除首尾空格

		// 格式化输出
		logx.Infof("[SQL] %s | %s\n%s", timeInfo, duration, query)
	} else {
		logx.Info(cleanStr)
	}

	return len(p), nil
}

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
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.WithWriter(&logxWriter{}),
	))

	// 创建表
	db.NewCreateTable().Model((*Resource)(nil)).IfNotExists().Exec(context.Background(),
		(*Resource)(nil),
	)
	db.NewCreateTable().Model((*Models)(nil)).IfNotExists().Exec(context.Background(),
		(*Models)(nil),
	)
	db.NewCreateTable().Model((*Tags)(nil)).IfNotExists().Exec(context.Background(),
		(*Tags)(nil),
	)

	// 初始化表数据
	NewModelsModel(db).InitData()
	return db
}

// 自定义数据库日志输出
