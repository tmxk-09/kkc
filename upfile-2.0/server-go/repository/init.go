package repository

import (
	"log"

	"upfile-server/config"

	"golang.org/x/crypto/bcrypt"
)

// initSchema 等价于 Java DbInitConfig：建表 + 兼容升级 + 预埋 admin
func initSchema() error {
	stmts := []string{
		// ===== files 表 =====
		`CREATE TABLE IF NOT EXISTS files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			file_name TEXT NOT NULL,
			file_size INTEGER NOT NULL,
			file_hash TEXT NOT NULL,
			store_path TEXT NOT NULL,
			mime_type TEXT DEFAULT '',
			created_at DATETIME DEFAULT (datetime('now','localtime')),
			downloads INTEGER DEFAULT 0,
			expire_days INTEGER DEFAULT 1,
			expire_at DATETIME DEFAULT NULL,
			pickup_code TEXT DEFAULT NULL,
			user_id INTEGER DEFAULT 1
		)`,
		`CREATE INDEX IF NOT EXISTS idx_files_hash ON files(file_hash)`,
		`CREATE INDEX IF NOT EXISTS idx_files_name ON files(file_name)`,
		`CREATE INDEX IF NOT EXISTS idx_files_expire ON files(expire_at)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_files_pickup ON files(pickup_code)`,
		`CREATE INDEX IF NOT EXISTS idx_files_user ON files(user_id)`,

		// ===== users 表 =====
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			phone TEXT,
			security_question TEXT,
			security_answer TEXT,
			password_hash TEXT NOT NULL,
			total_space INTEGER DEFAULT 1073741824,
			used_space INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT (datetime('now','localtime')),
			status INTEGER DEFAULT 1
		)`,

		// ===== space_gift_records 表 =====
		`CREATE TABLE IF NOT EXISTS space_gift_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			from_user_id INTEGER NOT NULL,
			to_user_id INTEGER NOT NULL,
			gift_size INTEGER NOT NULL,
			created_at DATETIME DEFAULT (datetime('now','localtime')),
			remark TEXT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_gift_from ON space_gift_records(from_user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_gift_to ON space_gift_records(to_user_id)`,

		// ===== upload_sessions 表 =====
		`CREATE TABLE IF NOT EXISTS upload_sessions (
			upload_id TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL DEFAULT 1,
			file_name TEXT NOT NULL,
			file_size INTEGER NOT NULL,
			file_hash TEXT NOT NULL,
			chunk_size INTEGER NOT NULL,
			total_chunk INTEGER NOT NULL,
			chunk_dir TEXT NOT NULL,
			created_at DATETIME DEFAULT (datetime('now','localtime'))
		)`,

		// ===== upload_chunks 表（保留但不主动写入）=====
		`CREATE TABLE IF NOT EXISTS upload_chunks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			upload_id TEXT NOT NULL,
			chunk_index INTEGER NOT NULL,
			chunk_size INTEGER NOT NULL,
			chunk_hash TEXT DEFAULT '',
			created_at DATETIME DEFAULT (datetime('now','localtime')),
			UNIQUE(upload_id, chunk_index)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_chunks_upload ON upload_chunks(upload_id)`,

		// ===== tickets 表 =====
		`CREATE TABLE IF NOT EXISTS tickets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			reply_content TEXT,
			status INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT (datetime('now','localtime')),
			updated_at DATETIME DEFAULT (datetime('now','localtime'))
		)`,

		// ===== messages 表 =====
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			is_read INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT (datetime('now','localtime'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_user ON messages(user_id)`,
	}

	// ALTER TABLE 兼容升级（字段不存在时新增，忽略错误）
	alterStmts := []string{
		`ALTER TABLE files ADD COLUMN expire_days INTEGER DEFAULT 1`,
		`ALTER TABLE files ADD COLUMN expire_at DATETIME DEFAULT NULL`,
		`ALTER TABLE files ADD COLUMN pickup_code TEXT DEFAULT NULL`,
		`ALTER TABLE files ADD COLUMN user_id INTEGER DEFAULT 1`,
		`ALTER TABLE users ADD COLUMN phone TEXT`,
		`ALTER TABLE users ADD COLUMN security_question TEXT`,
		`ALTER TABLE users ADD COLUMN security_answer TEXT`,
		`ALTER TABLE upload_sessions ADD COLUMN user_id INTEGER DEFAULT 1`,
	}

	for _, sql := range stmts {
		if _, err := DB.Exec(sql); err != nil {
			return err
		}
	}
	for _, sql := range alterStmts {
		DB.Exec(sql) // 忽略字段已存在的错误
	}

	// 预埋 admin 用户
	if err := seedAdmin(); err != nil {
		log.Printf("[DB] admin 用户预埋失败（可能已存在）: %v", err)
	}

	log.Println("[DB] 数据库 Schema 初始化完成")
	return nil
}

// seedAdmin 预埋管理员账号（若不存在）
func seedAdmin() error {
	var count int
	if err := DB.Get(&count, "SELECT COUNT(*) FROM users WHERE username='admin'"); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	adminPass := config.Cfg.FileServer.AdminPass
	if adminPass == "" {
		adminPass = "admin123"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = DB.Exec(
		`INSERT INTO users (username, password_hash, phone, security_question, security_answer, total_space)
		 VALUES ('admin', ?, '13800000000', '超级管理员预设', 'admin', 1073741824000)`,
		string(hash),
	)
	if err != nil {
		return err
	}

	// 转移历史无主文件给 admin
	DB.Exec("UPDATE files SET user_id = (SELECT id FROM users WHERE username='admin') WHERE user_id = 1 OR user_id IS NULL")

	log.Println("[DB] admin 用户已预埋")
	return nil
}
