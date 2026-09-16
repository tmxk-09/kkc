package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"upfile-server/model"
	"upfile-server/repository"

	"golang.org/x/crypto/bcrypt"
)

// GiftSpace 赠送存储空间（原子事务）
func GiftSpace(fromUserID int, toUsername string, giftSizeMB int64, verifyPassword string) error {
	if giftSizeMB <= 0 {
		return errors.New("赠送量必须大于 0")
	}
	giftBytes := giftSizeMB * 1024 * 1024

	// 1. 校验赠送方身份（二次密码确认）
	fromUser, err := repository.UserSelectByID(fromUserID)
	if err != nil {
		return errors.New("用户不存在")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(fromUser.PasswordHash), []byte(verifyPassword)); err != nil {
		return errors.New("密码验证失败，操作已拒绝")
	}

	// 2. 查找接收方
	toUser, err := repository.UserSelectByUsername(toUsername)
	if err != nil {
		return errors.New("目标用户名不存在")
	}
	if toUser.ID == fromUserID {
		return errors.New("不能赠送给自己")
	}

	// 3. 检查可用空间
	freeSpace := fromUser.TotalSpace - fromUser.UsedSpace
	if freeSpace < giftBytes {
		return fmt.Errorf("操作失败：自身可用空间不足（当前剩余 %.2f MB）", float64(freeSpace)/1024/1024)
	}

	// 4. 事务：扣减赠送方 + 增加接收方 + 插入流水
	tx, err := repository.DB.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 扣减赠送方 total_space
	if _, err := tx.Exec("UPDATE users SET total_space=total_space-? WHERE id=?", giftBytes, fromUserID); err != nil {
		tx.Rollback()
		return fmt.Errorf("扣减空间失败: %v", err)
	}

	// 增加接收方 total_space
	if _, err := tx.Exec("UPDATE users SET total_space=total_space+? WHERE id=?", giftBytes, toUser.ID); err != nil {
		tx.Rollback()
		return fmt.Errorf("增加目标空间失败: %v", err)
	}

	// 插入流水记录
	if _, err := tx.Exec(
		"INSERT INTO space_gift_records (from_user_id, to_user_id, gift_size, remark) VALUES (?,?,?,?)",
		fromUserID, toUser.ID, giftBytes,
		fmt.Sprintf("%s 赠送给 %s %.0f MB", fromUser.Username, toUsername, float64(giftBytes)/1024/1024),
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("记录流水失败: %v", err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("事务提交失败: %v", err)
	}

	log.Printf("[SpaceService] %s 赠送 %s %.0f MB 空间成功",
		fromUser.Username, toUsername, float64(giftBytes)/1024/1024)
	return nil
}

// GetGiftRecords 获取用户的赠送流水（简单查询）
func GetGiftRecords(userID int) ([]model.SpaceGiftRecord, error) {
	var records []model.SpaceGiftRecord
	err := repository.DB.Select(&records,
		"SELECT * FROM space_gift_records WHERE from_user_id=? OR to_user_id=? ORDER BY id DESC LIMIT 50",
		userID, userID)
	// 忽略 no rows 错误
	if errors.Is(err, sql.ErrNoRows) {
		return []model.SpaceGiftRecord{}, nil
	}
	return records, err
}
