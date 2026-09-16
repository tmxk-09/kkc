package repository

import "upfile-server/model"

// ========== UserRepo ==========

func UserSelectByID(id int) (*model.User, error) {
	var u model.User
	err := DB.Get(&u, "SELECT * FROM users WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func UserSelectByUsername(username string) (*model.User, error) {
	var u model.User
	err := DB.Get(&u, "SELECT * FROM users WHERE username=?", username)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func UserSelectAll() ([]model.User, error) {
	var users []model.User
	err := DB.Select(&users, "SELECT * FROM users ORDER BY id ASC")
	return users, err
}

func UserInsert(username, passwordHash, phone, sq, sa string, totalSpace int64) (int64, error) {
	res, err := DB.Exec(
		`INSERT INTO users (username, password_hash, phone, security_question, security_answer, total_space) VALUES (?,?,?,?,?,?)`,
		username, passwordHash, phone, sq, sa, totalSpace,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func UserUpdateStatus(id, status int) error {
	_, err := DB.Exec("UPDATE users SET status=? WHERE id=?", status, id)
	return err
}

func UserUpdateTotalSpace(id int, totalSpace int64) error {
	_, err := DB.Exec("UPDATE users SET total_space=? WHERE id=?", totalSpace, id)
	return err
}

func UserUpdatePassword(id int, passwordHash string) error {
	_, err := DB.Exec("UPDATE users SET password_hash=? WHERE id=?", passwordHash, id)
	return err
}

func UserAddUsedSpace(id int, size int64) error {
	_, err := DB.Exec("UPDATE users SET used_space=used_space+? WHERE id=?", size, id)
	return err
}

func UserFreeUsedSpace(id int, size int64) error {
	_, err := DB.Exec("UPDATE users SET used_space=MAX(0, used_space-?) WHERE id=?", size, id)
	return err
}


