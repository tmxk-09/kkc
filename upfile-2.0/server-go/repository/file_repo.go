package repository

import "upfile-server/model"

// ========== FileRepo ==========

func FileSelectByHash(hash string) (*model.FileInfo, error) {
	var f model.FileInfo
	err := DB.Get(&f, "SELECT * FROM files WHERE file_hash=? LIMIT 1", hash)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func FileSelectByID(id int64) (*model.FileInfo, error) {
	var f model.FileInfo
	err := DB.Get(&f, "SELECT * FROM files WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func FileSelectByPickupCode(code string) (*model.FileInfo, error) {
	var f model.FileInfo
	err := DB.Get(&f, "SELECT * FROM files WHERE pickup_code=?", code)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// FileSelectList 分页查询（admin 传 userID=0 查全部）
func FileSelectList(keyword string, userID int, offset, limit int) ([]model.FileInfo, error) {
	var files []model.FileInfo
	kw := "%" + keyword + "%"
	var err error
	if userID == 0 {
		err = DB.Select(&files,
			"SELECT * FROM files WHERE file_name LIKE ? ORDER BY id DESC LIMIT ? OFFSET ?",
			kw, limit, offset)
	} else {
		err = DB.Select(&files,
			"SELECT * FROM files WHERE file_name LIKE ? AND user_id=? ORDER BY id DESC LIMIT ? OFFSET ?",
			kw, userID, limit, offset)
	}
	return files, err
}

func FileCountList(keyword string, userID int) (int, error) {
	var total int
	kw := "%" + keyword + "%"
	var err error
	if userID == 0 {
		err = DB.Get(&total, "SELECT COUNT(*) FROM files WHERE file_name LIKE ?", kw)
	} else {
		err = DB.Get(&total, "SELECT COUNT(*) FROM files WHERE file_name LIKE ? AND user_id=?", kw, userID)
	}
	return total, err
}

func FileInsert(f *model.FileInfo) (int64, error) {
	res, err := DB.Exec(
		`INSERT INTO files (file_name, file_size, file_hash, store_path, mime_type, expire_days, expire_at, pickup_code, user_id)
		 VALUES (?,?,?,?,?,?,
		   CASE WHEN ? = -1 THEN NULL ELSE datetime('now','localtime','+' || ? || ' days') END,
		   ?,?)`,
		f.FileName, f.FileSize, f.FileHash, f.StorePath, f.MimeType,
		f.ExpireDays, f.ExpireDays, f.ExpireDays,
		f.PickupCode, f.UserID,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func FileDeleteByID(id int64) error {
	_, err := DB.Exec("DELETE FROM files WHERE id=?", id)
	return err
}

func FileIncrementDownloads(id int64) error {
	_, err := DB.Exec("UPDATE files SET downloads=downloads+1 WHERE id=?", id)
	return err
}

// FileSelectExpired 查询所有已过期文件
func FileSelectExpired() ([]model.FileInfo, error) {
	var files []model.FileInfo
	err := DB.Select(&files,
		"SELECT * FROM files WHERE expire_days != -1 AND expire_at IS NOT NULL AND expire_at < datetime('now','localtime')")
	return files, err
}

func FileDeleteExpired() (int64, error) {
	res, err := DB.Exec(
		"DELETE FROM files WHERE expire_days != -1 AND expire_at IS NOT NULL AND expire_at < datetime('now','localtime')")
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
