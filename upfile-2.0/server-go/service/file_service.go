package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"upfile-server/config"
	"upfile-server/model"
	"upfile-server/repository"
	"upfile-server/util"
)

// ========== 初始化 ==========

var uploadDir string
var chunkDir string

// 高性能缓冲池：复用大内存缓冲区，避免 GC 压力与频繁系统调用
var (
	chunkCopyBufPool = sync.Pool{
		New: func() any {
			b := make([]byte, 1024*1024) // 1MB 写入缓冲
			return &b
		},
	}
	mergeBufPool = sync.Pool{
		New: func() any {
			b := make([]byte, 2*1024*1024) // 2MB 合并流转缓冲
			return &b
		},
	}
)

func InitDirs() {
	uploadDir, _ = filepath.Abs(config.Cfg.FileServer.UploadDir)
	chunkDir, _ = filepath.Abs(config.Cfg.FileServer.ChunkDir)
	os.MkdirAll(uploadDir, 0755)
	os.MkdirAll(chunkDir, 0755)
	log.Printf("[FileService] uploadDir=%s, chunkDir=%s", uploadDir, chunkDir)
}

// ========== 秒传检查 ==========

// CheckHash 检查文件 Hash 是否已存在（秒传）
func CheckHash(fileHash string) map[string]any {
	file, err := repository.FileSelectByHash(fileHash)
	if err == nil && file != nil {
		if _, err2 := os.Stat(file.StorePath); err2 == nil {
			return map[string]any{"exists": true, "file": file}
		}
	}
	return map[string]any{"exists": false}
}

// InstantUpload 秒传完成 — 复用物理文件
func InstantUpload(fileHash, fileName string, fileSize int64, expireDays, userID int) (*model.FileInfo, error) {
	existing, err := repository.FileSelectByHash(fileHash)
	if err != nil || existing == nil {
		return nil, errors.New("秒传失败：文件不存在，请使用普通上传")
	}
	if _, err := os.Stat(existing.StorePath); err != nil {
		return nil, errors.New("秒传失败：物理文件已丢失，请使用普通上传")
	}

	// 检查用户空间
	user, err := repository.UserSelectByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if user.TotalSpace-user.UsedSpace < fileSize {
		return nil, errors.New("操作失败：空间余额不足以存放该文件")
	}

	newFile := &model.FileInfo{
		FileName:   fileName,
		FileSize:   fileSize,
		FileHash:   fileHash,
		StorePath:  existing.StorePath, // 复用物理文件
		MimeType:   existing.MimeType,
		ExpireDays: expireDays,
		PickupCode: util.GeneratePickupCode(func(code string) bool {
			f, _ := repository.FileSelectByPickupCode(code)
			return f != nil
		}),
		UserID: userID,
	}

	id, err := repository.FileInsert(newFile)
	if err != nil {
		return nil, fmt.Errorf("秒传记录保存失败: %v", err)
	}
	newFile.ID = id

	// 扣减空间配额
	repository.UserAddUsedSpace(userID, fileSize)

	log.Printf("[FileService] 秒传成功: %s (复用: %s, 用户: %d)", fileName, existing.StorePath, userID)
	return newFile, nil
}

// ========== 分片上传 ==========

// InitUpload 初始化分片上传（含断点续传逻辑）
func InitUpload(fileName, fileHash string, fileSize int64, userID int) (map[string]any, error) {
	// 检查空间
	user, err := repository.UserSelectByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if user.TotalSpace-user.UsedSpace < fileSize {
		return nil, errors.New("操作失败：空间余额不足以存放该文件")
	}

	// 断点续传：查找已有 session
	sessions, _ := repository.UploadSessionSelectByUserID(userID)
	for _, s := range sessions {
		if s.FileHash == fileHash {
			if _, err := os.Stat(s.ChunkDir); err == nil {
				log.Printf("[FileService] 断点续传恢复 uploadId=%s", s.UploadID)
				return map[string]any{
					"uploadId":   s.UploadID,
					"chunkSize":  s.ChunkSize,
					"totalChunk": s.TotalChunk,
				}, nil
			}
		}
	}

	// 并发任务数限制：每个用户最多同时执行 5 个上传任务
	if len(sessions) >= 5 {
		return nil, errors.New("并发上传数已达上限（每个用户最多允许同时上传5个任务），请等待前序任务完成")
	}

	// 清理该用户该文件的旧垃圾记录
	repository.UploadSessionDeleteByHashAndUser(fileHash, userID)

	// 创建新会话
	uploadID := generateUploadID()
	chunkSize := calcOptimalChunkSize(fileSize)
	totalChunk := int(math.Ceil(float64(fileSize) / float64(chunkSize)))

	sessionChunkDir := filepath.Join(chunkDir, uploadID)
	os.MkdirAll(sessionChunkDir, 0755)

	session := &model.UploadSession{
		UploadID:   uploadID,
		UserID:     userID,
		FileName:   fileName,
		FileSize:   fileSize,
		FileHash:   fileHash,
		ChunkSize:  chunkSize,
		TotalChunk: totalChunk,
		ChunkDir:   sessionChunkDir,
	}
	if err := repository.UploadSessionInsert(session); err != nil {
		return nil, fmt.Errorf("创建上传会话失败: %v", err)
	}

	return map[string]any{
		"uploadId":   uploadID,
		"chunkSize":  chunkSize,
		"totalChunk": totalChunk,
	}, nil
}

// UploadChunk 上传单个分片（使用 1MB 大缓冲降低系统调用开销）
func UploadChunk(uploadID string, chunkIndex int, src io.Reader, size int64) error {
	session, err := repository.UploadSessionSelectByID(uploadID)
	if err != nil {
		return fmt.Errorf("上传会话不存在: %s", uploadID)
	}

	chunkPath := filepath.Join(session.ChunkDir, strconv.Itoa(chunkIndex))
	os.MkdirAll(filepath.Dir(chunkPath), 0755)

	out, err := os.Create(chunkPath)
	if err != nil {
		return fmt.Errorf("创建分片文件失败: %v", err)
	}
	defer out.Close()

	bufPtr := chunkCopyBufPool.Get().(*[]byte)
	defer chunkCopyBufPool.Put(bufPtr)

	if _, err := io.CopyBuffer(out, src, *bufPtr); err != nil {
		return fmt.Errorf("写入分片失败: %v", err)
	}
	return nil
}

// GetProgress 查询上传进度（通过物理文件检测，无 DB 查询）
func GetProgress(uploadID string) (map[string]any, error) {
	session, err := repository.UploadSessionSelectByID(uploadID)
	if err != nil {
		return nil, fmt.Errorf("上传会话不存在: %s", uploadID)
	}

	entries, _ := os.ReadDir(session.ChunkDir)
	uploaded := make([]int, 0)
	for _, e := range entries {
		if idx, err := strconv.Atoi(e.Name()); err == nil {
			uploaded = append(uploaded, idx)
		}
	}

	return map[string]any{
		"uploadId":      uploadID,
		"totalChunk":    session.TotalChunk,
		"uploadedChunks": uploaded,
		"uploadedCount": len(uploaded),
	}, nil
}

// CompleteUpload 合并分片，完成上传（单次遍历边合并边计算MD5，零二次磁盘IO开销）
func CompleteUpload(uploadID string, expireDays, userID int) (*model.FileInfo, error) {
	session, err := repository.UploadSessionSelectByID(uploadID)
	if err != nil {
		return nil, errors.New("上传会话不存在")
	}

	// 检查分片完整性
	entries, _ := os.ReadDir(session.ChunkDir)
	if len(entries) < session.TotalChunk {
		return nil, fmt.Errorf("分片不完整，已上传: %d/%d", len(entries), session.TotalChunk)
	}

	// 目标文件路径
	dateDir := time.Now().Format("2006/01/02")
	targetDir := filepath.Join(uploadDir, dateDir)
	os.MkdirAll(targetDir, 0755)

	ext := filepath.Ext(session.FileName)
	storedName := session.FileHash[:8] + "_" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ext
	storePath := filepath.Join(targetDir, storedName)

	// 合并分片并同步计算哈希
	out, err := os.Create(storePath)
	if err != nil {
		return nil, fmt.Errorf("创建目标文件失败: %v", err)
	}

	hasher := md5.New()
	mw := io.MultiWriter(out, hasher)

	mergeBufPtr := mergeBufPool.Get().(*[]byte)
	defer mergeBufPool.Put(mergeBufPtr)

	for i := 0; i < session.TotalChunk; i++ {
		chunkPath := filepath.Join(session.ChunkDir, strconv.Itoa(i))
		chunk, err := os.Open(chunkPath)
		if err != nil {
			out.Close()
			os.Remove(storePath)
			return nil, fmt.Errorf("分片文件缺失: %d", i)
		}
		_, err = io.CopyBuffer(mw, chunk, *mergeBufPtr)
		chunk.Close()
		if err != nil {
			out.Close()
			os.Remove(storePath)
			return nil, fmt.Errorf("合并分片失败: %v", err)
		}
	}
	out.Close()

	// 零二次磁盘 I/O 校验：直接使用合并中产生的 MD5
	actualHash := hex.EncodeToString(hasher.Sum(nil))
	if actualHash != session.FileHash {
		os.Remove(storePath)
		return nil, fmt.Errorf("文件 Hash 校验失败，期望: %s，实际: %s", session.FileHash, actualHash)
	}

	// 生成取件码
	pickupCode := util.GeneratePickupCode(func(code string) bool {
		f, _ := repository.FileSelectByPickupCode(code)
		return f != nil
	})

	if expireDays == 0 {
		expireDays = 1
	}

	fileInfo := &model.FileInfo{
		FileName:   session.FileName,
		FileSize:   session.FileSize,
		FileHash:   session.FileHash,
		StorePath:  storePath,
		MimeType:   detectMimeType(session.FileName),
		ExpireDays: expireDays,
		PickupCode: pickupCode,
		UserID:     userID,
	}

	id, err := repository.FileInsert(fileInfo)
	if err != nil {
		os.Remove(storePath)
		return nil, fmt.Errorf("文件记录保存失败: %v", err)
	}
	fileInfo.ID = id

	// 扣减用户空间
	repository.UserAddUsedSpace(userID, session.FileSize)

	// 异步清理分片
	go func() {
		os.RemoveAll(session.ChunkDir)
		repository.UploadSessionDeleteByID(uploadID)
		log.Printf("[FileService] 分片清理完成: %s", uploadID)
	}()

	log.Printf("[FileService] 文件上传完成: %s (%d bytes)", session.FileName, session.FileSize)
	return fileInfo, nil
}

// ListPendingTasks 获取用户所有进行中的上传任务
func ListPendingTasks(userID int) ([]map[string]any, error) {
	sessions, err := repository.UploadSessionSelectByUserID(userID)
	if err != nil {
		return nil, err
	}

	tasks := make([]map[string]any, 0, len(sessions))
	for _, s := range sessions {
		entries, _ := os.ReadDir(s.ChunkDir)
		uploadedCount := len(entries)

		progress := int(math.Round(float64(uploadedCount) / float64(s.TotalChunk) * 95))
		status := "interrupted"
		if uploadedCount >= s.TotalChunk {
			status = "merging"
		}

		tasks = append(tasks, map[string]any{
			"uploadId":      s.UploadID,
			"fileName":      s.FileName,
			"fileSize":      s.FileSize,
			"fileHash":      s.FileHash,
			"totalChunk":    s.TotalChunk,
			"createdAt":     s.CreatedAt,
			"uploadedCount": uploadedCount,
			"progress":      progress,
			"status":        status,
		})
	}
	return tasks, nil
}

// DeleteUploadTask 删除上传任务
func DeleteUploadTask(uploadID string, userID int) error {
	session, err := repository.UploadSessionSelectByID(uploadID)
	if err != nil {
		return errors.New("上传任务不存在")
	}
	if session.UserID != userID {
		return errors.New("无权删除此任务")
	}

	os.RemoveAll(session.ChunkDir)
	repository.UploadSessionDeleteByID(uploadID)
	log.Printf("[FileService] 已删除上传任务: %s (用户: %d)", session.FileName, userID)
	return nil
}

// ========== 文件管理 ==========

// ListFiles 分页查询文件列表
func ListFiles(keyword string, userID int, pageNum, pageSize int) (map[string]any, error) {
	offset := (pageNum - 1) * pageSize
	list, err := repository.FileSelectList(keyword, userID, offset, pageSize)
	if err != nil {
		return nil, err
	}
	total, err := repository.FileCountList(keyword, userID)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"list":     list,
		"total":    total,
		"pageNum":  pageNum,
		"pageSize": pageSize,
	}, nil
}

// GetFileForDownload 获取文件信息（普通下载）
func GetFileForDownload(id int64) (*model.FileInfo, error) {
	f, err := repository.FileSelectByID(id)
	if err != nil {
		return nil, errors.New("文件不存在")
	}
	if _, err := os.Stat(f.StorePath); err != nil {
		return nil, errors.New("文件已丢失")
	}
	repository.FileIncrementDownloads(id)
	return f, nil
}

// GetFileByPickupCode 通过取件码获取文件
func GetFileByPickupCode(code string) (*model.FileInfo, error) {
	if code == "" {
		return nil, errors.New("取件码不能为空")
	}
	f, err := repository.FileSelectByPickupCode(code)
	if err != nil {
		return nil, errors.New("取件码无效或文件已过期、被销毁")
	}
	if _, err := os.Stat(f.StorePath); err != nil {
		return nil, errors.New("文件已损坏或被删除")
	}
	repository.FileIncrementDownloads(f.ID)
	return f, nil
}

// DeleteFile 删除文件并返还配额
func DeleteFile(id int64) error {
	f, err := repository.FileSelectByID(id)
	if err != nil {
		return errors.New("文件不存在")
	}
	os.Remove(f.StorePath)
	repository.FileDeleteByID(id)
	if f.UserID > 0 {
		repository.UserFreeUsedSpace(f.UserID, f.FileSize)
	}
	log.Printf("[FileService] 文件已删除: %s", f.FileName)
	return nil
}

// CleanExpiredFiles 清理过期文件（定时任务）
func CleanExpiredFiles() int {
	files, err := repository.FileSelectExpired()
	if err != nil {
		log.Printf("[FileService] 查询过期文件失败: %v", err)
		return 0
	}
	for _, f := range files {
		os.Remove(f.StorePath)
		if f.UserID > 0 {
			repository.UserFreeUsedSpace(f.UserID, f.FileSize)
		}
		log.Printf("[FileService] 过期文件删除: %s", f.FileName)
	}
	count, _ := repository.FileDeleteExpired()
	if count > 0 {
		log.Printf("[FileService] 本次清理过期文件: %d 个", count)
	}
	return int(count)
}

// ========== 工具 ==========

func calcOptimalChunkSize(fileSize int64) int64 {
	minChunk := int64(5 * 1024 * 1024)  // 5MB
	maxChunkMB := config.Cfg.FileServer.MaxChunkMB
	if maxChunkMB <= 0 {
		maxChunkMB = 100
	}
	maxChunk := int64(maxChunkMB) * 1024 * 1024
	maxParts := int64(10000)

	chunkSize := fileSize / maxParts
	if chunkSize < minChunk {
		return minChunk
	}
	if chunkSize > maxChunk {
		return maxChunk
	}
	return chunkSize
}

func detectMimeType(fileName string) string {
	ext := filepath.Ext(fileName)
	t := mime.TypeByExtension(ext)
	if t == "" {
		return "application/octet-stream"
	}
	return t
}

func generateUploadID() string {
	// 使用时间戳 + 随机数生成唯一 ID（等价 Java UUID）
	return fmt.Sprintf("%x%x",
		time.Now().UnixNano(),
		os.Getpid(),
	) + strconv.FormatInt(time.Now().UnixNano()%0xFFFF, 16)
}
