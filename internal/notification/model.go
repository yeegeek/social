package notification

import (
	"time"

	"gorm.io/datatypes"
)

// Notification 通知模型
type Notification struct {
	ID       int            `gorm:"primaryKey" json:"id"`
	UserID   int            `gorm:"not null" json:"user_id"`
	SenderID *int           `json:"sender_id,omitempty"`
	Type     string         `gorm:"not null" json:"type"` // friend_request, message, like, comment, system
	Title    string         `gorm:"not null" json:"title"`
	Content  string         `json:"content,omitempty"`
	Data     datatypes.JSON `json:"data,omitempty"` // 额外数据
	IsRead   bool           `gorm:"default:false" json:"is_read"`
	ReadAt   *time.Time     `json:"read_at,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
}

// TableName 指定表名
func (Notification) TableName() string {
	return "notifications"
}
