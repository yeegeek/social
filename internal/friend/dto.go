package friend

import (
	"time"

	"github.com/yeegeek/uyou-go-api-starter/internal/user"
)

// SendFriendRequestRequest 发送好友请求
type SendFriendRequestRequest struct {
	FriendID uint `json:"friend_id" binding:"required"`
}

// AcceptFriendRequestRequest 接受好友请求
type AcceptFriendRequestRequest struct {
	FriendshipID uint `json:"friendship_id" binding:"required"`
}

// RejectFriendRequestRequest 拒绝好友请求
type RejectFriendRequestRequest struct {
	FriendshipID uint `json:"friendship_id" binding:"required"`
}

// UpdateFriendRemarkRequest 更新好友备注
type UpdateFriendRemarkRequest struct {
	FriendID uint   `json:"friend_id" binding:"required"`
	Remark   string `json:"remark" binding:"required,max=255"`
}

// UpdateFriendGroupRequest 更新好友分组
type UpdateFriendGroupRequest struct {
	FriendID  uint   `json:"friend_id" binding:"required"`
	GroupName string `json:"group_name" binding:"required,max=100"`
}

// BlockUserRequest 拉黑用户
type BlockUserRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Reason string `json:"reason" binding:"max=500"`
}

// FriendshipResponse 好友关系响应
type FriendshipResponse struct {
	ID          uint             `json:"id"`
	UserID      uint             `json:"user_id"`
	FriendID    uint             `json:"friend_id"`
	Status      FriendshipStatus `json:"status"`
	GroupName   string           `json:"group_name,omitempty"`
	Remark      string           `json:"remark,omitempty"`
	RequestedAt time.Time        `json:"requested_at"`
	AcceptedAt  *time.Time       `json:"accepted_at,omitempty"`
	Friend      *FriendInfo      `json:"friend,omitempty"`
}

// FriendInfo 好友信息
type FriendInfo struct {
	ID           uint       `json:"id"`
	Name         string     `json:"name"`
	Username     string     `json:"username,omitempty"`
	Email        string     `json:"email"`
	AvatarURL    string     `json:"avatar_url,omitempty"`
	Gender       string     `json:"gender,omitempty"`
	Country      string     `json:"country,omitempty"`
	City         string     `json:"city,omitempty"`
	Bio          string     `json:"bio,omitempty"`
	IsVIP        bool       `json:"is_vip"`
	IsOnline     bool       `json:"is_online"`
	LastActiveAt *time.Time `json:"last_active_at,omitempty"`
}

// ToFriendInfo 转换用户模型为好友信息
func ToFriendInfo(u *user.User) *FriendInfo {
	if u == nil {
		return nil
	}
	return &FriendInfo{
		ID:           u.ID,
		Name:         u.Name,
		Username:     u.Username,
		Email:        u.Email,
		AvatarURL:    u.AvatarURL,
		Gender:       u.Gender,
		Country:      u.Country,
		City:         u.City,
		Bio:          u.Bio,
		IsVIP:        u.IsVIP,
		IsOnline:     u.IsOnline,
		LastActiveAt: u.LastActiveAt,
	}
}

// ToFriendshipResponse 转换好友关系模型为响应
func ToFriendshipResponse(f *Friendship) *FriendshipResponse {
	if f == nil {
		return nil
	}
	return &FriendshipResponse{
		ID:          f.ID,
		UserID:      f.UserID,
		FriendID:    f.FriendID,
		Status:      f.Status,
		GroupName:   f.GroupName,
		Remark:      f.Remark,
		RequestedAt: f.RequestedAt,
		AcceptedAt:  f.AcceptedAt,
		Friend:      ToFriendInfo(&f.Friend),
	}
}

// BlacklistResponse 黑名单响应
type BlacklistResponse struct {
	ID            uint        `json:"id"`
	UserID        uint        `json:"user_id"`
	BlockedUserID uint        `json:"blocked_user_id"`
	Reason        string      `json:"reason,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	BlockedUser   *FriendInfo `json:"blocked_user,omitempty"`
}

// ToBlacklistResponse 转换黑名单模型为响应
func ToBlacklistResponse(b *Blacklist) *BlacklistResponse {
	if b == nil {
		return nil
	}
	return &BlacklistResponse{
		ID:            b.ID,
		UserID:        b.UserID,
		BlockedUserID: b.BlockedUserID,
		Reason:        b.Reason,
		CreatedAt:     b.CreatedAt,
		BlockedUser:   ToFriendInfo(&b.BlockedUser),
	}
}
