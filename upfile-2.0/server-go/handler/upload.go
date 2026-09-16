package handler

import (
	"strconv"

	"upfile-server/middleware"
	"upfile-server/service"
	"upfile-server/util"

	"github.com/gin-gonic/gin"
)

// CheckHash POST /api/upload/check
func CheckHash(c *gin.Context) {
	var req struct {
		FileHash string `json:"fileHash"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.FileHash == "" {
		util.Fail(c, "fileHash不能为空")
		return
	}
	util.Ok(c, service.CheckHash(req.FileHash))
}

// InstantUpload POST /api/upload/instant
func InstantUpload(c *gin.Context) {
	var req struct {
		FileHash   string `json:"fileHash"`
		FileName   string `json:"fileName"`
		FileSize   int64  `json:"fileSize"`
		ExpireDays int    `json:"expireDays"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.FileHash == "" || req.FileName == "" || req.FileSize == 0 {
		util.Fail(c, "参数不完整：fileHash、fileName、fileSize不能为空")
		return
	}
	if req.ExpireDays == 0 {
		req.ExpireDays = 1
	}

	userID := middleware.GetUserID(c)
	fileInfo, err := service.InstantUpload(req.FileHash, req.FileName, req.FileSize, req.ExpireDays, userID)
	if err != nil {
		util.Fail(c, "秒传失败: "+err.Error())
		return
	}
	util.Ok(c, fileInfo)
}

// InitUpload POST /api/upload/init
func InitUpload(c *gin.Context) {
	var req struct {
		FileName string `json:"fileName"`
		FileHash string `json:"fileHash"`
		FileSize int64  `json:"fileSize"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.FileName == "" || req.FileHash == "" || req.FileSize == 0 {
		util.Fail(c, "参数不完整")
		return
	}

	userID := middleware.GetUserID(c)
	result, err := service.InitUpload(req.FileName, req.FileHash, req.FileSize, userID)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Ok(c, result)
}

// UploadChunk POST /api/upload/chunk  (multipart/form-data)
func UploadChunk(c *gin.Context) {
	uploadID := c.PostForm("uploadId")
	chunkIndexStr := c.PostForm("chunkIndex")
	if uploadID == "" || chunkIndexStr == "" {
		util.Fail(c, "uploadId 和 chunkIndex 不能为空")
		return
	}

	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil {
		util.Fail(c, "chunkIndex 格式错误")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		util.Fail(c, "分片文件缺失")
		return
	}

	src, err := file.Open()
	if err != nil {
		util.Fail(c, "打开分片失败")
		return
	}
	defer src.Close()

	if err := service.UploadChunk(uploadID, chunkIndex, src, file.Size); err != nil {
		util.Fail(c, "分片上传失败: "+err.Error())
		return
	}
	util.OkEmpty(c)
}

// CompleteUpload POST /api/upload/complete
func CompleteUpload(c *gin.Context) {
	var req struct {
		UploadID   string `json:"uploadId"`
		ExpireDays int    `json:"expireDays"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.UploadID == "" {
		util.Fail(c, "uploadId不能为空")
		return
	}
	if req.ExpireDays == 0 {
		req.ExpireDays = 1
	}

	userID := middleware.GetUserID(c)
	fileInfo, err := service.CompleteUpload(req.UploadID, req.ExpireDays, userID)
	if err != nil {
		util.Fail(c, "合并失败: "+err.Error())
		return
	}
	util.Ok(c, fileInfo)
}

// GetProgress POST /api/upload/progress
func GetProgress(c *gin.Context) {
	var req struct {
		UploadID string `json:"uploadId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.UploadID == "" {
		util.Fail(c, "uploadId不能为空")
		return
	}
	result, err := service.GetProgress(req.UploadID)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Ok(c, result)
}

// ListUploadTasks POST /api/upload/tasks
func ListUploadTasks(c *gin.Context) {
	userID := middleware.GetUserID(c)
	tasks, err := service.ListPendingTasks(userID)
	if err != nil {
		util.Fail(c, "获取任务列表失败: "+err.Error())
		return
	}
	util.Ok(c, tasks)
}

// DeleteUploadTask POST /api/upload/task/delete
func DeleteUploadTask(c *gin.Context) {
	var req struct {
		UploadID string `json:"uploadId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.UploadID == "" {
		util.Fail(c, "uploadId不能为空")
		return
	}
	userID := middleware.GetUserID(c)
	if err := service.DeleteUploadTask(req.UploadID, userID); err != nil {
		util.Fail(c, "删除任务失败: "+err.Error())
		return
	}
	util.OkEmpty(c)
}
