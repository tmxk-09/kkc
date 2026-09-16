package repository

import "upfile-server/model"

// ========== UploadRepo ==========

func UploadSessionInsert(s *model.UploadSession) error {
	_, err := DB.Exec(
		`INSERT INTO upload_sessions (upload_id, user_id, file_name, file_size, file_hash, chunk_size, total_chunk, chunk_dir)
		 VALUES (?,?,?,?,?,?,?,?)`,
		s.UploadID, s.UserID, s.FileName, s.FileSize, s.FileHash, s.ChunkSize, s.TotalChunk, s.ChunkDir,
	)
	return err
}

func UploadSessionSelectByID(uploadID string) (*model.UploadSession, error) {
	var s model.UploadSession
	err := DB.Get(&s, "SELECT * FROM upload_sessions WHERE upload_id=?", uploadID)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func UploadSessionSelectByUserID(userID int) ([]model.UploadSession, error) {
	var sessions []model.UploadSession
	err := DB.Select(&sessions,
		"SELECT * FROM upload_sessions WHERE user_id=? ORDER BY created_at DESC", userID)
	return sessions, err
}

func UploadSessionDeleteByID(uploadID string) error {
	_, err := DB.Exec("DELETE FROM upload_sessions WHERE upload_id=?", uploadID)
	return err
}

func UploadSessionDeleteByHashAndUser(fileHash string, userID int) error {
	_, err := DB.Exec("DELETE FROM upload_sessions WHERE file_hash=? AND user_id=?", fileHash, userID)
	return err
}
