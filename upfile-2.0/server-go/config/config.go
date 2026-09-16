package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config 全局配置结构体
type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`

	FileServer struct {
		UploadDir      string `mapstructure:"upload_dir"`
		ChunkDir       string `mapstructure:"chunk_dir"`
		JwtSecret      string `mapstructure:"jwt_secret"`
		JwtExpireHours int    `mapstructure:"jwt_expire_hours"`
		AdminUser      string `mapstructure:"admin_user"`
		AdminPass      string `mapstructure:"admin_pass"`
		MaxChunkMB     int    `mapstructure:"max_chunk_mb"`
		DbPath         string `mapstructure:"db_path"`
	} `mapstructure:"file_server"`
}

var Cfg Config

// Load 加载配置（优先 env 覆盖）
func Load(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/app")
	}

	// 环境变量支持，如 FILE_SERVER_JWT_SECRET
	viper.SetEnvPrefix("FILE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("[CONFIG] 未找到配置文件，使用默认值: %v", err)
	}

	// 设置默认值
	viper.SetDefault("server.port", 9102)
	viper.SetDefault("file_server.upload_dir", "./uploads/files")
	viper.SetDefault("file_server.chunk_dir", "./uploads/chunks")
	viper.SetDefault("file_server.jwt_secret", "file-server-secret-key-2024")
	viper.SetDefault("file_server.jwt_expire_hours", 72)
	viper.SetDefault("file_server.admin_user", "admin")
	viper.SetDefault("file_server.admin_pass", "L@#$12345699101")
	viper.SetDefault("file_server.max_chunk_mb", 100)
	viper.SetDefault("file_server.db_path", "./data/fileserver.db")

	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("[CONFIG] 配置解析失败: %v", err)
	}

	log.Printf("[CONFIG] 已加载配置，服务端口: %d", Cfg.Server.Port)
}
