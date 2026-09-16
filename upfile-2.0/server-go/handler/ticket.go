package handler

import (
	"upfile-server/middleware"
	"upfile-server/model"
	"upfile-server/repository"
	"upfile-server/util"

	"github.com/gin-gonic/gin"
)

// CreateTicket POST /api/ticket/create
func CreateTicket(c *gin.Context) {
	username := middleware.GetUsername(c)
	user, err := repository.UserSelectByUsername(username)
	if err != nil {
		util.Fail(c, "用户不存在")
		return
	}

	var req model.Ticket
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, "参数错误")
		return
	}
	if req.Title == "" {
		util.Fail(c, "工单标题不能为空")
		return
	}
	if req.Content == "" {
		util.Fail(c, "工单内容不能为空")
		return
	}
	req.UserID = user.ID
	if err := repository.TicketInsert(&req); err != nil {
		util.Fail(c, "创建工单失败: "+err.Error())
		return
	}
	util.OkEmpty(c)
}

// GetTicketList POST /api/ticket/list
func GetTicketList(c *gin.Context) {
	username := middleware.GetUsername(c)
	if middleware.IsAdmin(c) {
		tickets, err := repository.TicketSelectAll()
		if err != nil {
			util.Fail(c, err.Error())
			return
		}
		util.Ok(c, tickets)
		return
	}

	user, err := repository.UserSelectByUsername(username)
	if err != nil {
		util.Fail(c, "用户不存在")
		return
	}
	tickets, err := repository.TicketSelectByUserID(user.ID)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Ok(c, tickets)
}

// ReplyTicket POST /api/ticket/reply (Admin only)
func ReplyTicket(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		util.Fail(c, "权限不足")
		return
	}
	var req struct {
		ID           int    `json:"id"`
		ReplyContent string `json:"replyContent"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 || req.ReplyContent == "" {
		util.Fail(c, "缺少工单标识或回复内容不能为空")
		return
	}
	if err := repository.TicketUpdateReply(req.ID, req.ReplyContent); err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.OkEmpty(c)
}
