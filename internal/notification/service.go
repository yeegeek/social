package notification

import (
	"encoding/json"
	"errors"

	"gorm.io/datatypes"
)

// Service 通知服务接口
type Service interface {
	CreateNotification(req *CreateNotificationRequest) (*NotificationResponse, error)
	GetNotifications(userID int, page, pageSize int) (*NotificationListResponse, error)
	GetUnreadCount(userID int) (int64, error)
	MarkAsRead(notificationID, userID int) error
	MarkAllAsRead(userID int) error
	DeleteNotification(notificationID, userID int) error
}

type service struct {
	repo Repository
}

// NewService 创建通知服务
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateNotification 创建通知
func (s *service) CreateNotification(req *CreateNotificationRequest) (*NotificationResponse, error) {
	var dataJSON datatypes.JSON
	if req.Data != nil {
		data, err := json.Marshal(req.Data)
		if err != nil {
			return nil, err
		}
		dataJSON = datatypes.JSON(data)
	}

	notification := &Notification{
		UserID:   req.UserID,
		SenderID: req.SenderID,
		Type:     req.Type,
		Title:    req.Title,
		Content:  req.Content,
		Data:     dataJSON,
	}

	err := s.repo.CreateNotification(notification)
	if err != nil {
		return nil, err
	}

	return s.toNotificationResponse(notification), nil
}

// GetNotifications 获取通知列表
func (s *service) GetNotifications(userID int, page, pageSize int) (*NotificationListResponse, error) {
	offset := (page - 1) * pageSize
	notifications, total, err := s.repo.GetUserNotifications(userID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	unreadCount, _ := s.repo.GetUnreadCount(userID)

	response := &NotificationListResponse{
		Notifications: make([]NotificationResponse, 0, len(notifications)),
		Total:         total,
		UnreadCount:   unreadCount,
	}

	for _, notif := range notifications {
		response.Notifications = append(response.Notifications, *s.toNotificationResponse(&notif))
	}

	return response, nil
}

// GetUnreadCount 获取未读通知数量
func (s *service) GetUnreadCount(userID int) (int64, error) {
	return s.repo.GetUnreadCount(userID)
}

// MarkAsRead 标记通知为已读
func (s *service) MarkAsRead(notificationID, userID int) error {
	notification, err := s.repo.GetNotificationByID(notificationID)
	if err != nil {
		return err
	}

	// 验证权限
	if notification.UserID != userID {
		return errors.New("unauthorized to mark notification as read")
	}

	return s.repo.MarkAsRead(notificationID)
}

// MarkAllAsRead 标记所有通知为已读
func (s *service) MarkAllAsRead(userID int) error {
	return s.repo.MarkAllAsRead(userID)
}

// DeleteNotification 删除通知
func (s *service) DeleteNotification(notificationID, userID int) error {
	notification, err := s.repo.GetNotificationByID(notificationID)
	if err != nil {
		return err
	}

	// 验证权限
	if notification.UserID != userID {
		return errors.New("unauthorized to delete notification")
	}

	return s.repo.DeleteNotification(notificationID)
}

// toNotificationResponse 转换为通知响应
func (s *service) toNotificationResponse(notification *Notification) *NotificationResponse {
	var data map[string]interface{}
	if notification.Data != nil {
		_ = json.Unmarshal(notification.Data, &data)
	}

	return &NotificationResponse{
		ID:        notification.ID,
		UserID:    notification.UserID,
		SenderID:  notification.SenderID,
		Type:      notification.Type,
		Title:     notification.Title,
		Content:   notification.Content,
		Data:      data,
		IsRead:    notification.IsRead,
		ReadAt:    notification.ReadAt,
		CreatedAt: notification.CreatedAt,
	}
}
