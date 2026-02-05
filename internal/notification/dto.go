package notification

import "time"

// CreateNotificationRequest 创建通知请求
type CreateNotificationRequest struct {
	UserID   int                    `json:"user_id" validate:"required"`
	SenderID *int                   `json:"sender_id,omitempty"`
	Type     string                 `json:"type" validate:"required"`
	Title    string                 `json:"title" validate:"required"`
	Content  string                 `json:"content,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

// NotificationResponse 通知响应
type NotificationResponse struct {
	ID         int                    `json:"id"`
	UserID     int                    `json:"user_id"`
	SenderID   *int                   `json:"sender_id,omitempty"`
	SenderName string                 `json:"sender_name,omitempty"`
	Type       string                 `json:"type"`
	Title      string                 `json:"title"`
	Content    string                 `json:"content,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
	IsRead     bool                   `json:"is_read"`
	ReadAt     *time.Time             `json:"read_at,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}

// NotificationListResponse 通知列表响应
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Total         int64                  `json:"total"`
	UnreadCount   int64                  `json:"unread_count"`
}
