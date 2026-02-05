package notification

import (
	"time"

	"gorm.io/gorm"
)

// Repository 通知仓库接口
type Repository interface {
	CreateNotification(notification *Notification) error
	GetNotificationByID(id int) (*Notification, error)
	GetUserNotifications(userID int, limit, offset int) ([]Notification, int64, error)
	GetUnreadCount(userID int) (int64, error)
	MarkAsRead(notificationID int) error
	MarkAllAsRead(userID int) error
	DeleteNotification(notificationID int) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository 创建通知仓库
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateNotification 创建通知
func (r *repository) CreateNotification(notification *Notification) error {
	return r.db.Create(notification).Error
}

// GetNotificationByID 根据ID获取通知
func (r *repository) GetNotificationByID(id int) (*Notification, error) {
	var notification Notification
	err := r.db.First(&notification, id).Error
	return &notification, err
}

// GetUserNotifications 获取用户的通知列表
func (r *repository) GetUserNotifications(userID int, limit, offset int) ([]Notification, int64, error) {
	var notifications []Notification
	var total int64

	query := r.db.Model(&Notification{}).Where("user_id = ?", userID).Order("created_at DESC")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Limit(limit).Offset(offset).Find(&notifications).Error
	return notifications, total, err
}

// GetUnreadCount 获取未读通知数量
func (r *repository) GetUnreadCount(userID int) (int64, error) {
	var count int64
	err := r.db.Model(&Notification{}).Where("user_id = ? AND is_read = false", userID).Count(&count).Error
	return count, err
}

// MarkAsRead 标记通知为已读
func (r *repository) MarkAsRead(notificationID int) error {
	now := time.Now()
	return r.db.Model(&Notification{}).Where("id = ?", notificationID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// MarkAllAsRead 标记所有通知为已读
func (r *repository) MarkAllAsRead(userID int) error {
	now := time.Now()
	return r.db.Model(&Notification{}).Where("user_id = ? AND is_read = false", userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// DeleteNotification 删除通知
func (r *repository) DeleteNotification(notificationID int) error {
	return r.db.Delete(&Notification{}, notificationID).Error
}
