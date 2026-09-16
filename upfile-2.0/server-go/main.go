package main

import (
	"fmt"
	"log"
	"os"

	"upfile-server/config"
	"upfile-server/handler"
	"upfile-server/middleware"
	"upfile-server/repository"
	"upfile-server/service"
	"upfile-server/task"

	"github.com/gin-gonic/gin"
)

var (
	Version   = "2.0.0"
	BuildTime = "dev"
)

func main() {
	log.Printf("[Server] 🚀 正在启动 Upfile 服务 (版本: %s, 构建时间: %s)...", Version, BuildTime)

	// ===== 1. 加载配置 =====
	cfgFile := os.Getenv("CONFIG_FILE")
	config.Load(cfgFile)

	// ===== 2. 初始化数据库 =====
	repository.Init()

	// ===== 3. 初始化文件目录 =====
	service.InitDirs()

	// ===== 4. 启动定时清理任务 =====
	cronJob := task.StartCleanupTask()
	defer cronJob.Stop()

	// ===== 5. 配置 Gin =====
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 全局 CORS 处理
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// ===== 6. 路由注册 =====
	api := r.Group("/api")

	// --- 公开路由（无需登录）---
	auth := api.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/register", handler.Register)
	}

	// 文件公开访问（支持直连极速下载、多线程并发分段下载、取件码）
	filePublic := api.Group("/file")
	{
		filePublic.GET("/pickup/check/:code", handler.CheckPickupCode)
		filePublic.GET("/pickup/:code", handler.PickupFile)
		filePublic.HEAD("/pickup/:code", handler.PickupFile)
		filePublic.GET("/download/:id", handler.DownloadFile)
		filePublic.HEAD("/download/:id", handler.DownloadFile)
	}

	// --- 需登录的路由 ---
	protected := api.Group("")
	protected.Use(middleware.Auth())
	{
		// 上传模块
		upload := protected.Group("/upload")
		{
			upload.POST("/check", handler.CheckHash)
			upload.POST("/instant", handler.InstantUpload)
			upload.POST("/init", handler.InitUpload)
			upload.POST("/chunk", handler.UploadChunk)
			upload.POST("/complete", handler.CompleteUpload)
			upload.POST("/progress", handler.GetProgress)
			upload.POST("/tasks", handler.ListUploadTasks)
			upload.POST("/task/delete", handler.DeleteUploadTask)
		}

		// 文件管理
		file := protected.Group("/file")
		{
			file.POST("/list", handler.ListFiles)
			file.POST("/delete", handler.DeleteFile)
		}

		// 用户模块
		user := protected.Group("/user")
		{
			user.GET("/info", handler.GetUserInfo)
			user.GET("/list", handler.GetUserList)
			user.POST("/status", handler.UpdateUserStatus)
			user.POST("/space", handler.UpdateUserSpace)
			user.POST("/password", handler.UpdateUserPassword)
		}

		// 工单模块
		ticket := protected.Group("/ticket")
		{
			ticket.POST("/create", handler.CreateTicket)
			ticket.POST("/list", handler.GetTicketList)
			ticket.POST("/reply", handler.ReplyTicket)
		}

		// 站内信模块
		message := protected.Group("/message")
		{
			message.GET("/list", handler.GetMessageList)
			message.POST("/list", handler.GetMessageList)
			message.POST("/read", handler.ReadMessage)
			message.POST("/read_all", handler.ReadAllMessages)
			message.POST("/send", handler.SendMessage)
		}

		// 空间赠送模块
		space := protected.Group("/space")
		{
			space.POST("/gift", handler.GiftSpace)
		}
	}

	// ===== 7. 启动服务 =====
	addr := fmt.Sprintf(":%d", config.Cfg.Server.Port)
	log.Printf("[Server] upfile Go 服务启动，监听 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("[Server] 启动失败: %v", err)
	}
}
