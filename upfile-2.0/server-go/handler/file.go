package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"upfile-server/middleware"
	"upfile-server/model"
	"upfile-server/service"
	"upfile-server/util"

	"github.com/gin-gonic/gin"
)

// ListFiles POST /api/file/list
func ListFiles(c *gin.Context) {
	var req struct {
		Keyword  string `json:"keyword"`
		PageNum  int    `json:"pageNum"`
		PageSize int    `json:"pageSize"`
	}
	c.ShouldBindJSON(&req)
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	userID := middleware.GetUserID(c)
	username := middleware.GetUsername(c)

	// admin 查全部
	queryUserID := userID
	if username == "admin" {
		queryUserID = 0
	}

	result, err := service.ListFiles(req.Keyword, queryUserID, req.PageNum, req.PageSize)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Ok(c, result)
}

// DeleteFile POST /api/file/delete
func DeleteFile(c *gin.Context) {
	var req struct {
		ID int64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		util.Fail(c, "id不能为空")
		return
	}
	if err := service.DeleteFile(req.ID); err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.OkEmpty(c)
}

// DownloadFile GET /api/file/download/:id
func DownloadFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "id格式错误"})
		return
	}

	fileInfo, err := service.GetFileForDownload(id)
	if err != nil {
		c.JSON(404, gin.H{"code": 404, "msg": err.Error()})
		return
	}
	serveFileDownload(c, fileInfo)
}

// CheckPickupCode GET /api/file/pickup/check/:code
func CheckPickupCode(c *gin.Context) {
	code := strings.ToUpper(c.Param("code"))
	_, err := service.GetFileByPickupCode(code)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.OkEmpty(c)
}

// PickupFile GET /api/file/pickup/:code
func PickupFile(c *gin.Context) {
	code := strings.ToUpper(c.Param("code"))
	fileInfo, err := service.GetFileByPickupCode(code)
	if err != nil {
		c.JSON(404, gin.H{"code": 404, "msg": err.Error()})
		return
	}
	serveFileDownload(c, fileInfo)
}

// serveFileDownload 支持 HTTP Range 断点续传与内核零拷贝 sendfile
func serveFileDownload(c *gin.Context, fileInfo *model.FileInfo) {
	f, err := os.Open(fileInfo.StorePath)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "文件打开失败"})
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "文件信息获取失败"})
		return
	}

	// 文件名 URL 编码 (RFC 5987 / UTF-8)
	encodedName := url.QueryEscape(fileInfo.FileName)
	encodedName = strings.ReplaceAll(encodedName, "+", "%20")

	// 设置响应头
	if fileInfo.MimeType != "" {
		c.Header("Content-Type", fileInfo.MimeType)
	} else {
		c.Header("Content-Type", "application/octet-stream")
	}
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, encodedName))
	c.Header("Accept-Ranges", "bytes")
	c.Header("X-Accel-Buffering", "no")
	c.Header("Access-Control-Expose-Headers", "Content-Range, Accept-Ranges, Content-Length, Content-Disposition")

	// 使用 http.ServeContent：
	// 1. Linux 内核级 sendfile 零拷贝直通，CPU与内存开销近乎为零
	// 2. 完美支持多线程分段下载工具（IDM、Aria2、浏览器多连接加速）与 Range 断点续传
	// 3. 避免堆内存大缓冲分配
	http.ServeContent(c.Writer, c.Request, fileInfo.FileName, stat.ModTime(), f)
}
