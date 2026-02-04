package friend

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Repository 好友仓储接口
type Repository interface {
	// Friendship operations
	CreateFriendship(ctx context.Context, friendship *Friendship) error
	GetFriendshipByID(ctx context.Context, id uint) (*Friendship, error)
	GetFriendship(ctx context.Context, userID, friendID uint) (*Friendship, error)
	UpdateFriendship(ctx context.Context, friendship *Friendship) error
	DeleteFriendship(ctx context.Context, id uint) error
	ListFriends(ctx context.Context, userID uint, status FriendshipStatus) ([]*Friendship, error)
	ListFriendRequests(ctx context.Context, userID uint) ([]*Friendship, error)
	IsFriend(ctx context.Context, userID, friendID uint) (bool, error)

	// Blacklist operations
	CreateBlacklist(ctx context.Context, blacklist *Blacklist) error
	GetBlacklist(ctx context.Context, userID, blockedUserID uint) (*Blacklist, error)
	DeleteBlacklist(ctx context.Context, id uint) error
	ListBlacklist(ctx context.Context, userID uint) ([]*Blacklist, error)
	IsBlocked(ctx context.Context, userID, blockedUserID uint) (bool, error)
}

// repository 好友仓储实现
type repository struct {
	db *gorm.DB
}

// NewRepository 创建好友仓储实例
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateFriendship 创建好友关系
func (r *repository) CreateFriendship(ctx context.Context, friendship *Friendship) error {
	return r.db.WithContext(ctx).Create(friendship).Error
}

// GetFriendshipByID 根据ID获取好友关系
func (r *repository) GetFriendshipByID(ctx context.Context, id uint) (*Friendship, error) {
	var friendship Friendship
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Friend").
		First(&friendship, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &friendship, nil
}

// GetFriendship 获取两个用户之间的好友关系
func (r *repository) GetFriendship(ctx context.Context, userID, friendID uint) (*Friendship, error) {
	var friendship Friendship
	err := r.db.WithContext(ctx).
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", 
			userID, friendID, friendID, userID).
		Preload("User").
		Preload("Friend").
		First(&friendship).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &friendship, nil
}

// UpdateFriendship 更新好友关系
func (r *repository) UpdateFriendship(ctx context.Context, friendship *Friendship) error {
	return r.db.WithContext(ctx).Save(friendship).Error
}

// DeleteFriendship 删除好友关系
func (r *repository) DeleteFriendship(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Friendship{}, id).Error
}

// ListFriends 获取好友列表
func (r *repository) ListFriends(ctx context.Context, userID uint, status FriendshipStatus) ([]*Friendship, error) {
	var friendships []*Friendship
	query := r.db.WithContext(ctx).
		Where("user_id = ? OR friend_id = ?", userID, userID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.
		Preload("User").
		Preload("Friend").
		Order("created_at DESC").
		Find(&friendships).Error
	
	return friendships, err
}

// ListFriendRequests 获取好友请求列表
func (r *repository) ListFriendRequests(ctx context.Context, userID uint) ([]*Friendship, error) {
	var friendships []*Friendship
	err := r.db.WithContext(ctx).
		Where("friend_id = ? AND status = ?", userID, StatusPending).
		Preload("User").
		Preload("Friend").
		Order("created_at DESC").
		Find(&friendships).Error
	
	return friendships, err
}

// IsFriend 检查是否为好友
func (r *repository) IsFriend(ctx context.Context, userID, friendID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Friendship{}).
		Where("((user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)) AND status = ?",
			userID, friendID, friendID, userID, StatusAccepted).
		Count(&count).Error
	
	return count > 0, err
}

// CreateBlacklist 创建黑名单记录
func (r *repository) CreateBlacklist(ctx context.Context, blacklist *Blacklist) error {
	return r.db.WithContext(ctx).Create(blacklist).Error
}

// GetBlacklist 获取黑名单记录
func (r *repository) GetBlacklist(ctx context.Context, userID, blockedUserID uint) (*Blacklist, error) {
	var blacklist Blacklist
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND blocked_user_id = ?", userID, blockedUserID).
		Preload("User").
		Preload("BlockedUser").
		First(&blacklist).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &blacklist, nil
}

// DeleteBlacklist 删除黑名单记录
func (r *repository) DeleteBlacklist(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Blacklist{}, id).Error
}

// ListBlacklist 获取黑名单列表
func (r *repository) ListBlacklist(ctx context.Context, userID uint) ([]*Blacklist, error) {
	var blacklist []*Blacklist
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("User").
		Preload("BlockedUser").
		Order("created_at DESC").
		Find(&blacklist).Error
	
	return blacklist, err
}

// IsBlocked 检查是否被拉黑
func (r *repository) IsBlocked(ctx context.Context, userID, blockedUserID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Blacklist{}).
		Where("user_id = ? AND blocked_user_id = ?", userID, blockedUserID).
		Count(&count).Error
	
	return count > 0, err
}
