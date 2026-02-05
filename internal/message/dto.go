package message

import "time"

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	ReceiverID  int    `json:"receiver_id" validate:"required"`
	Content     string `json:"content" validate:"required"`
	MessageType string `json:"message_type" validate:"omitempty,oneof=text image video audio gift"`
	MediaURL    string `json:"media_url,omitempty"`
	GiftID      *int   `json:"gift_id,omitempty"`
}

// SaveDraftRequest 保存草稿请求
type SaveDraftRequest struct {
	ReceiverID int    `json:"receiver_id" validate:"required"`
	Content    string `json:"content" validate:"required"`
}

// RecallMessageRequest 撤回消息请求
type RecallMessageRequest struct {
	MessageID int `json:"message_id" validate:"required"`
}

// TranslateMessageRequest 翻译消息请求
type TranslateMessageRequest struct {
	MessageID      int    `json:"message_id" validate:"required"`
	TargetLanguage string `json:"target_language" validate:"required"`
}

// MessageResponse 消息响应
type MessageResponse struct {
	ID                int        `json:"id"`
	ConversationID    int        `json:"conversation_id"`
	SenderID          int        `json:"sender_id"`
	ReceiverID        int        `json:"receiver_id"`
	Content           string     `json:"content"`
	MessageType       string     `json:"message_type"`
	MediaURL          string     `json:"media_url,omitempty"`
	TranslatedContent string     `json:"translated_content,omitempty"`
	IsRead            bool       `json:"is_read"`
	ReadAt            *time.Time `json:"read_at,omitempty"`
	IsRecalled        bool       `json:"is_recalled"`
	IsDraft           bool       `json:"is_draft"`
	CreatedAt         time.Time  `json:"created_at"`
}

// ConversationResponse 对话响应
type ConversationResponse struct {
	ID            int              `json:"id"`
	OtherUserID   int              `json:"other_user_id"`
	OtherUserName string           `json:"other_user_name"`
	OtherUserAvatar string         `json:"other_user_avatar,omitempty"`
	LastMessage   *MessageResponse `json:"last_message,omitempty"`
	UnreadCount   int              `json:"unread_count"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

// ConversationListResponse 对话列表响应
type ConversationListResponse struct {
	Conversations []ConversationResponse `json:"conversations"`
	Total         int64                  `json:"total"`
}

// MessageListResponse 消息列表响应
type MessageListResponse struct {
	Messages []MessageResponse `json:"messages"`
	Total    int64             `json:"total"`
}
