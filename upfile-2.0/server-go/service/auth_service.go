package service

import (
	"errors"
	"fmt"

	"upfile-server/config"
	"upfile-server/model"
	"upfile-server/repository"
	"upfile-server/util"

	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证服务

// Login 用户登录，返回 JWT Token + 用户信息
func Login(username, password string) (map[string]any, error) {
	user, err := repository.UserSelectByUsername(username)
	if err != nil {
		// 兼容 admin（配置文件密码）
		if username == config.Cfg.FileServer.AdminUser {
			return loginAdmin(username, password)
		}
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查封禁状态
	if user.Status == 0 {
		return nil, errors.New("账号已被封禁，请联系管理员")
	}

	token, err := util.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("Token 生成失败: %v", err)
	}

	user.PasswordHash = ""
	return map[string]any{
		"token":      token,
		"username":   user.Username,
		"userId":     user.ID,
		"totalSpace": user.TotalSpace,
		"usedSpace":  user.UsedSpace,
		"user":       user,
	}, nil
}

// loginAdmin 配置文件 admin 兜底登录（向后兼容）
func loginAdmin(username, password string) (map[string]any, error) {
	// 先尝试从数据库找
	user, err := repository.UserSelectByUsername(username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := util.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("Token 生成失败: %v", err)
	}
	user.PasswordHash = ""
	return map[string]any{
		"token":      token,
		"username":   user.Username,
		"userId":     user.ID,
		"totalSpace": user.TotalSpace,
		"usedSpace":  user.UsedSpace,
		"user":       user,
	}, nil
}

// Register 用户注册
func Register(username, password string) error {
	if len(username) < 2 || len(username) > 30 {
		return errors.New("用户名长度须在 2~30 个字符之间")
	}
	if len(password) < 6 {
		return errors.New("密码长度不能少于 6 位")
	}

	// 检查是否已存在
	if _, err := repository.UserSelectByUsername(username); err == nil {
		return errors.New("用户名已被占用")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	// 默认 1GB 空间
	defaultSpace := int64(1073741824)
	newUserID, err := repository.UserInsert(username, string(hash), "", "", "", defaultSpace)
	if err != nil {
		return fmt.Errorf("注册失败: %v", err)
	}

	// 自动为新注册用户同步全局广播系统消息
	_ = repository.MessageCopyBroadcastForUser(int(newUserID))

	return nil
}

// GetUserInfo 获取当前用户信息
func GetUserInfo(username string) (*model.User, error) {
	user, err := repository.UserSelectByUsername(username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	user.PasswordHash = ""
	user.SecurityAnswer = ""
	return user, nil
}
