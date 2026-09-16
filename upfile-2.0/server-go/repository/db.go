package repository

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"upfile-server/config"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var DB *sqlx.DB

// Init 初始化 SQLite 数据库连接
func Init() {
	dbPath := config.Cfg.FileServer.DbPath

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("[DB] 创建数据库目录失败: %v", err)
	}

	dsn := fmt.Sprintf("file:%s?_journal_mode=WAL&_busy_timeout=5000&cache=shared", dbPath)
	db, err := sqlx.Connect("sqlite", dsn)
	if err != nil {
		log.Fatalf("[DB] SQLite 连接失败: %v", err)
	}

	// SQLite 最佳实践：限制最大连接数为 1（防止 WAL 并发写冲突）
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	DB = db
	log.Printf("[DB] SQLite 已连接: %s", dbPath)

	// 执行建表初始化
	if err := initSchema(); err != nil {
		log.Fatalf("[DB] 数据库初始化失败: %v", err)
	}
}
