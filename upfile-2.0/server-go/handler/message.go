package handler

import (
	"strconv"
	"strings"
	"upfile-server/middleware"
	"upfile-server/model"
	"upfile-server/repository"
	"upfile-server/util"

	"github.com/gin-gonic/gin"
)

// GetMessageList POST /api/message/list
func GetMessageList(c *gin.Context) {
	username := middleware.GetUsername(c)
	user, err := repository.UserSelectByUsername(username)
	if err != nil {
		util.Fail(c, "用户不存在")
		return
	}

	msgs, err := repository.MessageSelectByUserID(user.ID)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	unread, _ := repository.MessageCountUnread(user.ID)

	util.Ok(c, map[string]any{
		"list":        msgs,
		"unreadCount": unread,
	})
}

// ReadMessage POST /api/message/read
func ReadMessage(c *gin.Context) {
	username := middleware.GetUsername(c)
	user, err := repository.UserSelectByUsername(username)
	if err != nil {
		util.Fail(c, "用户失效")
		return
	}
	var req struct {
		ID int `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err == nil && req.ID > 0 {
		repository.MessageMarkRead(req.ID, user.ID)
	}
	util.OkEmpty(c)
}

// ReadAllMessages POST /api/message/read_all
func ReadAllMessages(c *gin.Context) {
	username := middleware.GetUsername(c)
	user, err := repository.UserSelectByUsername(username)
	if err != nil {
		util.Fail(c, "用户失效")
		return
	}
	repository.MessageMarkAllRead(user.ID)
	util.OkEmpty(c)
}

// SendMessage POST /api/message/send (Admin only)
func SendMessage(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		util.Fail(c, "仅限管理员发送系统通知")
		return
	}

	var req struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		TargetID any    `json:"targetId"` // "all" or int or string
		Target   any    `json:"target"`   // 兼容 target
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Content) == "" {
		util.Fail(c, "标题和内容均不能为空")
		return
	}
	if req.TargetID == nil && req.Target != nil {
		req.TargetID = req.Target
	}

	isAll := false
	targetUID := 0

	if req.TargetID == nil {
		isAll = true
	} else {
		switch v := req.TargetID.(type) {
		case string:
			trimmed := strings.TrimSpace(v)
			if trimmed == "" || trimmed == "all" || trimmed == "0" {
				isAll = true
			} else if id, err := strconv.Atoi(trimmed); err == nil {
				targetUID = id
			} else {
				util.Fail(c, "非法的目标用户ID")
				return
			}
		case float64:
			if int(v) <= 0 {
				isAll = true
			} else {
				targetUID = int(v)
			}
		case int:
			if v <= 0 {
				isAll = true
			} else {
				targetUID = v
			}
		default:
			isAll = true
		}
	}

	if isAll {
		// 1. 存入 user_id=0 作为公共全局系统通知底表
		_ = repository.MessageInsert(&model.Message{UserID: 0, Title: req.Title, Content: req.Content})

		// 2. 遍历所有现有账号，给每个人都插入一条专有记录，保证每个人独立已读未读
		users, _ := repository.UserSelectAll()
		batch := make([]model.Message, 0, len(users))
		for _, u := range users {
			batch = append(batch, model.Message{UserID: u.ID, Title: req.Title, Content: req.Content})
		}
		if len(batch) > 0 {
			if err := repository.MessageInsertBatch(batch); err != nil {
				util.Fail(c, "发送通知失败: "+err.Error())
				return
			}
		}
	} else {
		msg := model.Message{UserID: targetUID, Title: req.Title, Content: req.Content}
		if err := repository.MessageInsert(&msg); err != nil {
			util.Fail(c, "发送失败: "+err.Error())
			return
		}
	}

	util.OkEmpty(c)
}
