package repository

import "upfile-server/model"

// ========== TicketRepo ==========

func TicketInsert(t *model.Ticket) error {
	_, err := DB.Exec(
		`INSERT INTO tickets (user_id, title, content) VALUES (?,?,?)`,
		t.UserID, t.Title, t.Content,
	)
	return err
}

func TicketSelectAll() ([]model.Ticket, error) {
	tickets := make([]model.Ticket, 0)
	err := DB.Select(&tickets, `SELECT id, user_id, title, content, COALESCE(reply_content, '') AS reply_content, status, created_at, COALESCE(updated_at, '') AS updated_at FROM tickets ORDER BY id DESC`)
	return tickets, err
}

func TicketSelectByUserID(userID int) ([]model.Ticket, error) {
	tickets := make([]model.Ticket, 0)
	err := DB.Select(&tickets, `SELECT id, user_id, title, content, COALESCE(reply_content, '') AS reply_content, status, created_at, COALESCE(updated_at, '') AS updated_at FROM tickets WHERE user_id=? ORDER BY id DESC`, userID)
	return tickets, err
}

func TicketUpdateReply(id int, replyContent string) error {
	_, err := DB.Exec(
		`UPDATE tickets SET reply_content=?, status=1, updated_at=datetime('now','localtime') WHERE id=?`,
		replyContent, id,
	)
	return err
}

// ========== MessageRepo ==========

func MessageInsert(m *model.Message) error {
	_, err := DB.Exec(
		`INSERT INTO messages (user_id, title, content) VALUES (?,?,?)`,
		m.UserID, m.Title, m.Content,
	)
	return err
}

func MessageInsertBatch(msgs []model.Message) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO messages (user_id, title, content) VALUES (?,?,?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	for _, m := range msgs {
		if _, err := stmt.Exec(m.UserID, m.Title, m.Content); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func MessageSelectByUserID(userID int) ([]model.Message, error) {
	msgs := make([]model.Message, 0)
	err := DB.Select(&msgs, `
		SELECT id, user_id, title, content, is_read, created_at 
		FROM messages 
		WHERE user_id = ? 
		   OR (user_id = 0 AND title NOT IN (SELECT title FROM messages WHERE user_id = ?))
		ORDER BY id DESC`, userID, userID)
	return msgs, err
}

func MessageCountUnread(userID int) (int, error) {
	var count int
	err := DB.Get(&count, `
		SELECT COUNT(*) 
		FROM messages 
		WHERE (user_id = ? AND is_read = 0)
		   OR (user_id = 0 AND title NOT IN (SELECT title FROM messages WHERE user_id = ?))`, userID, userID)
	return count, err
}

func MessageMarkRead(id, userID int) error {
	_, err := DB.Exec("UPDATE messages SET is_read=1 WHERE id=? AND (user_id=? OR user_id=0)", id, userID)
	return err
}

func MessageMarkAllRead(userID int) error {
	_, err := DB.Exec("UPDATE messages SET is_read=1 WHERE user_id=? OR user_id=0", userID)
	return err
}

// MessageCopyBroadcastForUser 为新注册用户分发历史全局广播系统消息
func MessageCopyBroadcastForUser(userID int) error {
	_, err := DB.Exec(`INSERT INTO messages (user_id, title, content, is_read, created_at)
		SELECT ?, title, content, 0, created_at FROM messages WHERE user_id = 0`, userID)
	return err
}

// ========== SpaceGiftRepo ==========

func SpaceGiftInsert(r *model.SpaceGiftRecord) error {
	_, err := DB.Exec(
		`INSERT INTO space_gift_records (from_user_id, to_user_id, gift_size, remark) VALUES (?,?,?,?)`,
		r.FromUserID, r.ToUserID, r.GiftSize, r.Remark,
	)
	return err
}
