package model

// Ticket 工单实体，对应 tickets 表
type Ticket struct {
	ID           int    `db:"id" json:"id"`
	UserID       int    `db:"user_id" json:"userId"`
	Title        string `db:"title" json:"title"`
	Content      string `db:"content" json:"content"`
	ReplyContent string `db:"reply_content" json:"replyContent"`
	Status       int    `db:"status" json:"status"`
	CreatedAt    string `db:"created_at" json:"createdAt"`
	UpdatedAt    string `db:"updated_at" json:"updatedAt"`
}

// Message 站内信实体，对应 messages 表
type Message struct {
	ID        int    `db:"id" json:"id"`
	UserID    int    `db:"user_id" json:"userId"`
	Title     string `db:"title" json:"title"`
	Content   string `db:"content" json:"content"`
	IsRead    int    `db:"is_read" json:"isRead"`
	CreatedAt string `db:"created_at" json:"createdAt"`
}

// SpaceGiftRecord 空间赠送流水，对应 space_gift_records 表
type SpaceGiftRecord struct {
	ID         int    `db:"id" json:"id"`
	FromUserID int    `db:"from_user_id" json:"fromUserId"`
	ToUserID   int    `db:"to_user_id" json:"toUserId"`
	GiftSize   int64  `db:"gift_size" json:"giftSize"`
	CreatedAt  string `db:"created_at" json:"createdAt"`
	Remark     string `db:"remark" json:"remark"`
}
