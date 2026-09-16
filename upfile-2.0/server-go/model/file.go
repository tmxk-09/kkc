package model

// FileInfo 文件信息实体，对应 files 表
type FileInfo struct {
	ID          int64  `db:"id" json:"id"`
	FileName    string `db:"file_name" json:"fileName"`
	FileSize    int64  `db:"file_size" json:"fileSize"`
	FileHash    string `db:"file_hash" json:"fileHash"`
	StorePath   string `db:"store_path" json:"storePath"`
	MimeType    string `db:"mime_type" json:"mimeType"`
	CreatedAt   string `db:"created_at" json:"createdAt"`
	Downloads   int    `db:"downloads" json:"downloads"`
	ExpireDays  int    `db:"expire_days" json:"expireDays"`
	ExpireAt    string `db:"expire_at" json:"expireAt"`
	PickupCode  string `db:"pickup_code" json:"pickupCode"`
	UserID      int    `db:"user_id" json:"userId"`
}
