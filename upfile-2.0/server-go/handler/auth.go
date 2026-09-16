package handler

import (
	"upfile-server/middleware"
	"upfile-server/service"
	"upfile-server/util"

	"github.com/gin-gonic/gin"
)

// Login POST /api/auth/login
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		util.Fail(c, "参数错误")
		return
	}

	result, err := service.Login(req.Username, req.Password)
	if err != nil {
		util.FailCode(c, 401, err.Error())
		return
	}
	util.Ok(c, result)
}

// Register POST /api/auth/register
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, "参数错误")
		return
	}
	if err := service.Register(req.Username, req.Password); err != nil {
		util.FailCode(c, 400, err.Error())
		return
	}
	util.OkEmpty(c)
}

// GetUserInfo GET /api/user/info
func GetUserInfo(c *gin.Context) {
	username := middleware.GetUsername(c)
	user, err := service.GetUserInfo(username)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Ok(c, user)
}
