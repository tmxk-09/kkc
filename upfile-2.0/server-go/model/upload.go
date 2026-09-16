package model

// UploadSession 上传会话，对应 upload_sessions 表
type UploadSession struct {
	UploadID   string `db:"upload_id" json:"uploadId"`
	UserID     int    `db:"user_id" json:"userId"`
	FileName   string `db:"file_name" json:"fileName"`
	FileSize   int64  `db:"file_size" json:"fileSize"`
	FileHash   string `db:"file_hash" json:"fileHash"`
	ChunkSize  int64  `db:"chunk_size" json:"chunkSize"`
	TotalChunk int    `db:"total_chunk" json:"totalChunk"`
	ChunkDir   string `db:"chunk_dir" json:"chunkDir"`
	CreatedAt  string `db:"created_at" json:"createdAt"`
}
