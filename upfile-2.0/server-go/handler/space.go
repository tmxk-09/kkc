package handler

import (
	"upfile-server/middleware"
	"upfile-server/service"
	"upfile-server/util"

	"github.com/gin-gonic/gin"
)

// GiftSpace POST /api/space/gift
func GiftSpace(c *gin.Context) {
	var req struct {
		ToUsername     string `json:"toUsername"`
		GiftSizeMb     int64  `json:"giftSizeMb"`
		VerifyPassword string `json:"verifyPassword"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, "参数错误")
		return
	}
	if req.ToUsername == "" || req.GiftSizeMb <= 0 || req.VerifyPassword == "" {
		util.Fail(c, "参数不完整")
		return
	}

	fromUserID := middleware.GetUserID(c)
	if err := service.GiftSpace(fromUserID, req.ToUsername, req.GiftSizeMb, req.VerifyPassword); err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.OkEmpty(c)
}
