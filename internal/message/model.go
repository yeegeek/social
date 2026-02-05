package message

import (
	"time"
)

// Conversation 对话模型
type Conversation struct {
	ID               int        `gorm:"primaryKey" json:"id"`
	User1ID          int        `gorm:"not null" json:"user1_id"`
	User2ID          int        `gorm:"not null" json:"user2_id"`
	LastMessageID    *int       `json:"last_message_id,omitempty"`
	LastMessageAt    *time.Time `json:"last_message_at,omitempty"`
	User1UnreadCount int        `gorm:"default:0" json:"user1_unread_count"`
	User2UnreadCount int        `gorm:"default:0" json:"user2_unread_count"`
	User1Deleted     bool       `gorm:"default:false" json:"user1_deleted"`
	User2Deleted     bool       `gorm:"default:false" json:"user2_deleted"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (Conversation) TableName() string {
	return "conversations"
}

// Message 消息模型
type Message struct {
	ID                int        `gorm:"primaryKey" json:"id"`
	ConversationID    int        `gorm:"not null" json:"conversation_id"`
	SenderID          int        `gorm:"not null" json:"sender_id"`
	ReceiverID        int        `gorm:"not null" json:"receiver_id"`
	Content           string     `gorm:"not null" json:"content"`
	MessageType       string     `gorm:"default:text" json:"message_type"` // text, image, video, audio, gift
	MediaURL          string     `json:"media_url,omitempty"`
	TranslatedContent string     `json:"translated_content,omitempty"`
	IsRead            bool       `gorm:"default:false" json:"is_read"`
	ReadAt            *time.Time `json:"read_at,omitempty"`
	IsRecalled        bool       `gorm:"default:false" json:"is_recalled"`
	RecalledAt        *time.Time `json:"recalled_at,omitempty"`
	IsDraft           bool       `gorm:"default:false" json:"is_draft"`
	GiftID            *int       `json:"gift_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}
