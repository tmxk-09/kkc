package model

// User 用户实体，对应 users 表
type User struct {
	ID               int    `db:"id" json:"id"`
	Username         string `db:"username" json:"username"`
	Phone            string `db:"phone" json:"phone"`
	SecurityQuestion string `db:"security_question" json:"securityQuestion"`
	SecurityAnswer   string `db:"security_answer" json:"-"`
	PasswordHash     string `db:"password_hash" json:"-"`
	TotalSpace       int64  `db:"total_space" json:"totalSpace"`
	UsedSpace        int64  `db:"used_space" json:"usedSpace"`
	CreatedAt        string `db:"created_at" json:"createdAt"`
	Status           int    `db:"status" json:"status"`
}
