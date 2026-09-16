package handler

import (
	"upfile-server/middleware"
	"upfile-server/repository"
	"upfile-server/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetUserList GET /api/user/list (Admin only)
func GetUserList(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		util.Fail(c, "权限不足：非法越权探测")
		return
	}
	users, err := repository.UserSelectAll()
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	// 清除密码哈希
	for i := range users {
		users[i].PasswordHash = ""
		users[i].SecurityAnswer = ""
	}
	util.Ok(c, users)
}

// UpdateUserStatus POST /api/user/status (Admin only)
func UpdateUserStatus(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		util.Fail(c, "权限不足")
		return
	}
	var req struct {
		ID     int `json:"id"`
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		util.Fail(c, "缺少账户标志或状态属性")
		return
	}
	if err := repository.UserUpdateStatus(req.ID, req.Status); err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.OkEmpty(c)
}

// UpdateUserSpace POST /api/user/space (Admin only)
func UpdateUserSpace(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		util.Fail(c, "权限不足")
		return
	}
	var req struct {
		ID         int   `json:"id"`
		TotalSpace int64 `json:"totalSpace"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		util.Fail(c, "缺少账户标志或最新配额")
		return
	}
	if err := repository.UserUpdateTotalSpace(req.ID, req.TotalSpace); err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.OkEmpty(c)
}

// UpdateUserPassword POST /api/user/password (Admin only)
func UpdateUserPassword(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		util.Fail(c, "权限不足")
		return
	}
	var req struct {
		ID       int    `json:"id"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 || req.Password == "" {
		util.Fail(c, "缺少账户标志或新密码不能为空")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		util.Fail(c, "密码加密失败")
		return
	}
	if err := repository.UserUpdatePassword(req.ID, string(hash)); err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.OkEmpty(c)
}
