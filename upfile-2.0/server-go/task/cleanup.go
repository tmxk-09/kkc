package task

import (
	"log"

	"upfile-server/service"

	"github.com/robfig/cron/v3"
)

// StartCleanupTask 启动过期文件清理定时任务（每小时执行一次）
// 等价 Java FileCleanupTask + @Scheduled
func StartCleanupTask() *cron.Cron {
	c := cron.New()

	// 每小时整点执行
	c.AddFunc("0 * * * *", func() {
		log.Println("[CleanupTask] 开始清理过期文件...")
		count := service.CleanExpiredFiles()
		if count > 0 {
			log.Printf("[CleanupTask] 本次清理过期文件: %d 个", count)
		}
	})

	c.Start()
	log.Println("[CleanupTask] 过期文件清理任务已启动（每小时执行）")
	return c
}
